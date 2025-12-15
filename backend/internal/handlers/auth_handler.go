package handlers

import (
	"errors"
	"net/http"

	"github.com/joaosantos/pettime/internal/models"
	"github.com/joaosantos/pettime/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SocialLoginRequest struct {
	Provider   string `json:"provider"`
	Token      string `json:"token"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProviderID string `json:"provider_id"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User   *models.User        `json:"user"`
	Tokens *models.AuthTokens  `json:"tokens"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondError(w, http.StatusBadRequest, "Email, password, and name are required")
		return
	}

	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	input := models.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	user, tokens, err := h.authService.Register(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			respondError(w, http.StatusConflict, "User already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	respondCreated(w, AuthResponse{User: user, Tokens: tokens})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	input := models.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	user, tokens, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	respondSuccess(w, AuthResponse{User: user, Tokens: tokens})
}

func (h *AuthHandler) SocialLogin(w http.ResponseWriter, r *http.Request) {
	var req SocialLoginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Provider == "" || req.ProviderID == "" || req.Email == "" {
		respondError(w, http.StatusBadRequest, "Provider, provider_id, and email are required")
		return
	}

	var provider models.AuthProvider
	switch req.Provider {
	case "google":
		provider = models.AuthProviderGoogle
	case "apple":
		provider = models.AuthProviderApple
	default:
		respondError(w, http.StatusBadRequest, "Invalid provider")
		return
	}

	input := models.SocialLoginInput{
		Provider:   provider,
		Token:      req.Token,
		Name:       req.Name,
		Email:      req.Email,
		ProviderID: req.ProviderID,
	}

	user, tokens, err := h.authService.SocialLogin(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	respondSuccess(w, AuthResponse{User: user, Tokens: tokens})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	respondSuccess(w, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken != "" {
		_ = h.authService.Logout(r.Context(), req.RefreshToken)
	}

	respondNoContent(w)
}
