package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joaosantos/pettime/internal/config"
	"github.com/joaosantos/pettime/internal/database"
	"github.com/joaosantos/pettime/internal/handlers"
	"github.com/joaosantos/pettime/internal/middleware"
	"github.com/joaosantos/pettime/internal/repositories"
	"github.com/joaosantos/pettime/internal/services"
	"github.com/joaosantos/pettime/pkg/jwt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
	}

	// Connect to database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT manager
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.Pool)
	petRepo := repositories.NewPetRepository(db.Pool)
	activityRepo := repositories.NewActivityRepository(db.Pool)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtManager, cfg.JWT.RefreshTokenTTL)
	petService := services.NewPetService(petRepo, activityRepo)
	activityService := services.NewActivityService(activityRepo, petRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	petHandler := handlers.NewPetHandler(petService)
	activityHandler := handlers.NewActivityHandler(activityService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Health(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy"}`))
			return
		}
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/social", authHandler.SocialLogin)
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/logout", authHandler.Logout)
		})

		// Public reference data
		r.Get("/pet-types", petHandler.ListPetTypes)
		r.Get("/game-types", activityHandler.ListGameTypes)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Pets
			r.Route("/pets", func(r chi.Router) {
				r.Post("/", petHandler.Create)
				r.Get("/", petHandler.List)
				r.Get("/{id}", petHandler.GetByID)
				r.Put("/{id}", petHandler.Update)
				r.Delete("/{id}", petHandler.Delete)
				r.Get("/{id}/stats", petHandler.GetStats)
			})

			// Activities
			r.Route("/activities", func(r chi.Router) {
				r.Post("/", activityHandler.Create)
				r.Get("/", activityHandler.List)
				r.Get("/{id}", activityHandler.GetByID)
				r.Put("/{id}", activityHandler.Update)
				r.Post("/sync", activityHandler.Sync)
			})
		})
	})

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
