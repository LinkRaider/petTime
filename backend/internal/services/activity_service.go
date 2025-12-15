package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/joaosantos/pettime/internal/models"
	"github.com/joaosantos/pettime/internal/repositories"
)

var (
	ErrActivityNotFound = errors.New("activity not found")
	ErrInvalidGameType  = errors.New("invalid game type")
)

type ActivityService struct {
	activityRepo *repositories.ActivityRepository
	petRepo      *repositories.PetRepository
}

func NewActivityService(activityRepo *repositories.ActivityRepository, petRepo *repositories.PetRepository) *ActivityService {
	return &ActivityService{
		activityRepo: activityRepo,
		petRepo:      petRepo,
	}
}

func (s *ActivityService) Create(ctx context.Context, userID uuid.UUID, input models.CreateActivityInput) (*models.Activity, error) {
	// Verify pet belongs to user
	pet, err := s.petRepo.GetByID(ctx, input.PetID)
	if err != nil {
		return nil, ErrPetNotFound
	}
	if pet.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Verify game type exists and is supported for this pet type
	gameType, err := s.activityRepo.GetGameType(ctx, input.GameTypeID)
	if err != nil {
		return nil, ErrInvalidGameType
	}

	if !isGameTypeSupported(gameType, pet.PetTypeID) {
		return nil, ErrInvalidGameType
	}

	now := time.Now()
	activity := &models.Activity{
		ID:         uuid.New(),
		PetID:      input.PetID,
		GameTypeID: input.GameTypeID,
		GameType:   gameType,
		StartedAt:  input.StartedAt,
		EndedAt:    input.EndedAt,
		GameData:   input.GameData,
		ClientID:   input.ClientID,
		CreatedAt:  now,
	}

	// If activity is already completed, calculate XP
	if input.EndedAt != nil {
		duration := int(input.EndedAt.Sub(input.StartedAt).Seconds())
		activity.DurationSeconds = &duration
		activity.XPEarned = s.calculateXP(gameType, activity)

		// Update pet XP
		if err := s.petRepo.AddXP(ctx, pet.ID, activity.XPEarned); err != nil {
			return nil, err
		}

		// Update streak
		if err := s.updateStreak(ctx, pet); err != nil {
			return nil, err
		}
	}

	if input.ClientID != nil {
		syncedAt := now
		activity.SyncedAt = &syncedAt
	}

	if err := s.activityRepo.Create(ctx, activity); err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *ActivityService) GetByID(ctx context.Context, userID, activityID uuid.UUID) (*models.Activity, error) {
	activity, err := s.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		if errors.Is(err, repositories.ErrActivityNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	// Verify ownership through pet
	pet, err := s.petRepo.GetByID(ctx, activity.PetID)
	if err != nil {
		return nil, err
	}
	if pet.UserID != userID {
		return nil, ErrUnauthorized
	}

	return activity, nil
}

func (s *ActivityService) List(ctx context.Context, userID uuid.UUID, filter models.ActivityFilter) ([]*models.Activity, error) {
	// If filtering by pet, verify ownership
	if filter.PetID != nil {
		pet, err := s.petRepo.GetByID(ctx, *filter.PetID)
		if err != nil {
			return nil, err
		}
		if pet.UserID != userID {
			return nil, ErrUnauthorized
		}
	}

	return s.activityRepo.List(ctx, filter)
}

func (s *ActivityService) Update(ctx context.Context, userID, activityID uuid.UUID, input models.UpdateActivityInput) (*models.Activity, error) {
	activity, err := s.GetByID(ctx, userID, activityID)
	if err != nil {
		return nil, err
	}

	if input.EndedAt != nil && activity.EndedAt == nil {
		activity.EndedAt = input.EndedAt
		duration := int(input.EndedAt.Sub(activity.StartedAt).Seconds())
		activity.DurationSeconds = &duration

		// Get game type for XP calculation
		gameType, err := s.activityRepo.GetGameType(ctx, activity.GameTypeID)
		if err != nil {
			return nil, err
		}

		activity.XPEarned = s.calculateXP(gameType, activity)

		// Update pet XP
		if err := s.petRepo.AddXP(ctx, activity.PetID, activity.XPEarned); err != nil {
			return nil, err
		}

		// Update streak
		pet, err := s.petRepo.GetByID(ctx, activity.PetID)
		if err != nil {
			return nil, err
		}
		if err := s.updateStreak(ctx, pet); err != nil {
			return nil, err
		}
	}

	if input.GameData != nil {
		activity.GameData = input.GameData
	}

	if err := s.activityRepo.Update(ctx, activity); err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *ActivityService) Sync(ctx context.Context, userID uuid.UUID, activities []models.SyncActivityInput) ([]*models.Activity, error) {
	var synced []*models.Activity

	for _, input := range activities {
		// Check if already synced
		existing, err := s.activityRepo.GetByClientID(ctx, input.ClientID)
		if err == nil && existing != nil {
			synced = append(synced, existing)
			continue
		}

		// Create new activity
		createInput := models.CreateActivityInput{
			PetID:      input.PetID,
			GameTypeID: input.GameTypeID,
			StartedAt:  input.StartedAt,
			EndedAt:    input.EndedAt,
			GameData:   input.GameData,
			ClientID:   &input.ClientID,
		}

		activity, err := s.Create(ctx, userID, createInput)
		if err != nil {
			continue // Skip failed syncs
		}

		synced = append(synced, activity)
	}

	return synced, nil
}

func (s *ActivityService) GetAllGameTypes(ctx context.Context) ([]*models.GameType, error) {
	return s.activityRepo.GetAllGameTypes(ctx)
}

func (s *ActivityService) calculateXP(gameType *models.GameType, activity *models.Activity) int {
	var xpConfig struct {
		BaseXPPerMinute    float64 `json:"base_xp_per_minute"`
		DistanceBonusPerKM float64 `json:"distance_bonus_per_km"`
		StreakMultiplier   float64 `json:"streak_multiplier"`
		XPPerThrow         int     `json:"xp_per_throw"`
		ComboBonus         int     `json:"combo_bonus"`
		FrenzyMultiplier   float64 `json:"frenzy_multiplier"`
	}

	if err := json.Unmarshal(gameType.XPConfig, &xpConfig); err != nil {
		return 0
	}

	xp := 0

	switch gameType.ID {
	case "walk":
		// Base XP from duration
		if activity.DurationSeconds != nil {
			minutes := float64(*activity.DurationSeconds) / 60
			xp += int(minutes * xpConfig.BaseXPPerMinute)
		}

		// Distance bonus
		var walkData models.WalkGameData
		if err := json.Unmarshal(activity.GameData, &walkData); err == nil {
			distanceKM := walkData.DistanceMeters / 1000
			xp += int(distanceKM * xpConfig.DistanceBonusPerKM)
		}

	case "fetch":
		var fetchData models.FetchGameData
		if err := json.Unmarshal(activity.GameData, &fetchData); err == nil {
			xp += fetchData.Throws * xpConfig.XPPerThrow
			xp += (fetchData.MaxCombo / 5) * xpConfig.ComboBonus
			if fetchData.FrenzyModeActivated {
				xp = int(float64(xp) * xpConfig.FrenzyMultiplier)
			}
		}
	}

	return xp
}

func (s *ActivityService) updateStreak(ctx context.Context, pet *models.Pet) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if pet.LastActivityAt == nil {
		return s.petRepo.UpdateStreak(ctx, pet.ID, 1)
	}

	lastActivityDay := time.Date(
		pet.LastActivityAt.Year(),
		pet.LastActivityAt.Month(),
		pet.LastActivityAt.Day(),
		0, 0, 0, 0, pet.LastActivityAt.Location(),
	)

	daysSinceLastActivity := int(today.Sub(lastActivityDay).Hours() / 24)

	switch daysSinceLastActivity {
	case 0:
		// Same day, no streak change
		return nil
	case 1:
		// Consecutive day, increment streak
		return s.petRepo.UpdateStreak(ctx, pet.ID, pet.StreakDays+1)
	default:
		// Streak broken, reset to 1
		return s.petRepo.UpdateStreak(ctx, pet.ID, 1)
	}
}

func isGameTypeSupported(gameType *models.GameType, petTypeID string) bool {
	for _, supported := range gameType.SupportedPetTypes {
		if supported == petTypeID {
			return true
		}
	}
	return false
}
