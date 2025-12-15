package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/joaosantos/pettime/internal/models"
	"github.com/joaosantos/pettime/internal/repositories"
	"github.com/joaosantos/pettime/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	userRepo        *repositories.UserRepository
	jwtManager      *jwt.Manager
	refreshTokenTTL time.Duration
}

func NewAuthService(userRepo *repositories.UserRepository, jwtManager *jwt.Manager, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		jwtManager:      jwtManager,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, input models.CreateUserInput) (*models.User, *models.AuthTokens, error) {
	// Check if user exists
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	passwordHash := string(hashedPassword)

	// Create user
	now := time.Now()
	user := &models.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: &passwordHash,
		Name:         input.Name,
		AuthProvider: models.AuthProviderEmail,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			return nil, nil, ErrUserExists
		}
		return nil, nil, err
	}

	// Generate tokens
	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, input models.LoginInput) (*models.User, *models.AuthTokens, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if user.PasswordHash == nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) SocialLogin(ctx context.Context, input models.SocialLoginInput) (*models.User, *models.AuthTokens, error) {
	// Try to find existing user by provider
	user, err := s.userRepo.GetByProvider(ctx, input.Provider, input.ProviderID)
	if err != nil && !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, nil, err
	}

	if user == nil {
		// Try to find by email
		user, err = s.userRepo.GetByEmail(ctx, input.Email)
		if err != nil && !errors.Is(err, repositories.ErrUserNotFound) {
			return nil, nil, err
		}

		if user == nil {
			// Create new user
			providerID := input.ProviderID
			now := time.Now()
			user = &models.User{
				ID:             uuid.New(),
				Email:          input.Email,
				Name:           input.Name,
				AuthProvider:   input.Provider,
				AuthProviderID: &providerID,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, nil, err
			}
		}
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthTokens, error) {
	tokenHash := hashToken(refreshToken)

	storedToken, err := s.userRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}

	// Delete old refresh token
	_ = s.userRepo.DeleteRefreshToken(ctx, tokenHash)

	// Generate new tokens
	return s.generateTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return s.userRepo.DeleteRefreshToken(ctx, tokenHash)
}

func (s *AuthService) generateTokens(ctx context.Context, user *models.User) (*models.AuthTokens, error) {
	// Generate access token
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, err
	}
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	// Store refresh token
	now := time.Now()
	tokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: now.Add(s.refreshTokenTTL),
		CreatedAt: now,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, tokenRecord); err != nil {
		return nil, err
	}

	return &models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtManager.GetAccessTokenTTL().Seconds()),
	}, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
