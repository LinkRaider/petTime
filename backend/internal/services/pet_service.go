package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/joaosantos/pettime/internal/models"
	"github.com/joaosantos/pettime/internal/repositories"
)

var (
	ErrPetNotFound    = errors.New("pet not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidPetType = errors.New("invalid pet type")
)

type PetService struct {
	petRepo      *repositories.PetRepository
	activityRepo *repositories.ActivityRepository
}

func NewPetService(petRepo *repositories.PetRepository, activityRepo *repositories.ActivityRepository) *PetService {
	return &PetService{
		petRepo:      petRepo,
		activityRepo: activityRepo,
	}
}

func (s *PetService) Create(ctx context.Context, userID uuid.UUID, input models.CreatePetInput) (*models.Pet, error) {
	// Validate pet type exists
	petType, err := s.petRepo.GetPetType(ctx, input.PetTypeID)
	if err != nil {
		return nil, ErrInvalidPetType
	}

	now := time.Now()
	pet := &models.Pet{
		ID:         uuid.New(),
		UserID:     userID,
		PetTypeID:  input.PetTypeID,
		PetType:    petType,
		Name:       input.Name,
		Breed:      input.Breed,
		AvatarURL:  input.AvatarURL,
		BirthDate:  input.BirthDate,
		TotalXP:    0,
		Level:      1,
		Mood:       models.MoodHappy,
		StreakDays: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.petRepo.Create(ctx, pet); err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *PetService) GetByID(ctx context.Context, userID, petID uuid.UUID) (*models.Pet, error) {
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		if errors.Is(err, repositories.ErrPetNotFound) {
			return nil, ErrPetNotFound
		}
		return nil, err
	}

	if pet.UserID != userID {
		return nil, ErrUnauthorized
	}

	return pet, nil
}

func (s *PetService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Pet, error) {
	return s.petRepo.GetByUserID(ctx, userID)
}

func (s *PetService) Update(ctx context.Context, userID, petID uuid.UUID, input models.UpdatePetInput) (*models.Pet, error) {
	pet, err := s.GetByID(ctx, userID, petID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		pet.Name = *input.Name
	}
	if input.Breed != nil {
		pet.Breed = input.Breed
	}
	if input.AvatarURL != nil {
		pet.AvatarURL = input.AvatarURL
	}
	if input.BirthDate != nil {
		pet.BirthDate = input.BirthDate
	}

	if err := s.petRepo.Update(ctx, pet); err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *PetService) Delete(ctx context.Context, userID, petID uuid.UUID) error {
	pet, err := s.GetByID(ctx, userID, petID)
	if err != nil {
		return err
	}

	return s.petRepo.Delete(ctx, pet.ID)
}

func (s *PetService) GetStats(ctx context.Context, userID, petID uuid.UUID) (*models.Pet, *models.PetStats, error) {
	pet, err := s.GetByID(ctx, userID, petID)
	if err != nil {
		return nil, nil, err
	}

	stats, err := s.activityRepo.GetPetStats(ctx, petID)
	if err != nil {
		return nil, nil, err
	}

	// Calculate level progress
	stats.CurrentStreak = pet.StreakDays
	stats.XPToNextLevel = models.XPToNextLevel(pet.TotalXP)

	currentLevelXP := models.XPForLevel(pet.Level)
	nextLevelXP := models.XPForLevel(pet.Level + 1)
	if nextLevelXP > currentLevelXP {
		stats.LevelProgress = float64(pet.TotalXP-currentLevelXP) / float64(nextLevelXP-currentLevelXP)
	}

	return pet, stats, nil
}

func (s *PetService) GetAllPetTypes(ctx context.Context) ([]*models.PetType, error) {
	return s.petRepo.GetAllPetTypes(ctx)
}

func (s *PetService) UpdateMood(ctx context.Context, petID uuid.UUID) error {
	pet, err := s.petRepo.GetByID(ctx, petID)
	if err != nil {
		return err
	}

	newMood := calculateMood(pet)
	if newMood != pet.Mood {
		return s.petRepo.UpdateMood(ctx, petID, newMood)
	}

	return nil
}

func calculateMood(pet *models.Pet) models.Mood {
	if pet.LastActivityAt == nil {
		return models.MoodBored
	}

	hoursSinceActivity := time.Since(*pet.LastActivityAt).Hours()

	switch {
	case hoursSinceActivity < 6:
		return models.MoodHappy
	case hoursSinceActivity < 12:
		return models.MoodContent
	case hoursSinceActivity < 24:
		return models.MoodTired
	case hoursSinceActivity < 48:
		return models.MoodSad
	default:
		return models.MoodBored
	}
}
