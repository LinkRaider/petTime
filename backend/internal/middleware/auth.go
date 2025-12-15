package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/joaosantos/pettime/pkg/jwt"
)

type contextKey string

const (
	userIDKey    contextKey = "userID"
	userEmailKey contextKey = "userEmail"
)

type AuthMiddleware struct {
	jwtManager *jwt.Manager
}

func NewAuthMiddleware(jwtManager *jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Unauthorized","message":"Missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error":"Unauthorized","message":"Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			http.Error(w, `{"error":"Unauthorized","message":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, userEmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}

func GetUserEmail(ctx context.Context) string {
	email, ok := ctx.Value(userEmailKey).(string)
	if !ok {
		return ""
	}
	return email
}
