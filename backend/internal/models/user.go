package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderApple  AuthProvider = "apple"
)

type User struct {
	ID             uuid.UUID    `json:"id"`
	Email          string       `json:"email"`
	PasswordHash   *string      `json:"-"`
	Name           string       `json:"name"`
	AvatarURL      *string      `json:"avatar_url,omitempty"`
	AuthProvider   AuthProvider `json:"auth_provider"`
	AuthProviderID *string      `json:"-"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type CreateUserInput struct {
	Email          string       `json:"email" validate:"required,email"`
	Password       string       `json:"password" validate:"required,min=8"`
	Name           string       `json:"name" validate:"required,min=2"`
	AuthProvider   AuthProvider `json:"auth_provider"`
	AuthProviderID string       `json:"auth_provider_id"`
}

type UpdateUserInput struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SocialLoginInput struct {
	Provider   AuthProvider `json:"provider" validate:"required"`
	Token      string       `json:"token" validate:"required"`
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	ProviderID string       `json:"provider_id"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
