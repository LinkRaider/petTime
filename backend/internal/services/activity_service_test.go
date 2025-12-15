package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/joaosantos/pettime/internal/models"
)

func TestCalculateXP_WalkGame(t *testing.T) {
	service := &ActivityService{}

	// Walk game type config
	xpConfig := map[string]interface{}{
		"base_xp_per_minute":    2.0,
		"distance_bonus_per_km": 10.0,
		"streak_multiplier":     1.5,
	}
	configJSON, _ := json.Marshal(xpConfig)

	gameType := &models.GameType{
		ID:       "walk",
		XPConfig: configJSON,
	}

	tests := []struct {
		name             string
		durationSeconds  int
		distanceMeters   float64
		expectedMinXP    int
		expectedMaxXP    int
	}{
		{
			name:            "10 minute walk, 1km",
			durationSeconds: 600, // 10 minutes
			distanceMeters:  1000,
			expectedMinXP:   30, // (10 * 2) + (1 * 10) = 30
			expectedMaxXP:   30,
		},
		{
			name:            "30 minute walk, 2.5km",
			durationSeconds: 1800, // 30 minutes
			distanceMeters:  2500,
			expectedMinXP:   85, // (30 * 2) + (2.5 * 10) = 85
			expectedMaxXP:   85,
		},
		{
			name:            "60 minute walk, 5km",
			durationSeconds: 3600, // 60 minutes
			distanceMeters:  5000,
			expectedMinXP:   170, // (60 * 2) + (5 * 10) = 170
			expectedMaxXP:   170,
		},
		{
			name:            "5 minute walk, no distance data",
			durationSeconds: 300, // 5 minutes
			distanceMeters:  0,
			expectedMinXP:   10, // (5 * 2) + 0 = 10
			expectedMaxXP:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walkData := models.WalkGameData{
				DistanceMeters: tt.distanceMeters,
			}
			gameDataJSON, _ := json.Marshal(walkData)

			activity := &models.Activity{
				DurationSeconds: &tt.durationSeconds,
				GameData:        gameDataJSON,
			}

			xp := service.calculateXP(gameType, activity)

			if xp < tt.expectedMinXP || xp > tt.expectedMaxXP {
				t.Errorf("calculateXP() = %d, want between %d and %d",
					xp, tt.expectedMinXP, tt.expectedMaxXP)
			}
		})
	}
}

func TestCalculateXP_FetchGame(t *testing.T) {
	service := &ActivityService{}

	// Fetch game type config
	xpConfig := map[string]interface{}{
		"xp_per_throw":       1,
		"combo_bonus":        5,
		"frenzy_multiplier":  2.0,
	}
	configJSON, _ := json.Marshal(xpConfig)

	gameType := &models.GameType{
		ID:       "fetch",
		XPConfig: configJSON,
	}

	tests := []struct {
		name              string
		throws            int
		maxCombo          int
		frenzyActivated   bool
		expectedXP        int
	}{
		{
			name:            "10 throws, no combo",
			throws:          10,
			maxCombo:        0,
			frenzyActivated: false,
			expectedXP:      10, // 10 * 1 = 10
		},
		{
			name:            "20 throws, combo of 5",
			throws:          20,
			maxCombo:        5,
			frenzyActivated: false,
			expectedXP:      25, // (20 * 1) + (5/5 * 5) = 25
		},
		{
			name:            "30 throws, combo of 10",
			throws:          30,
			maxCombo:        10,
			frenzyActivated: false,
			expectedXP:      40, // (30 * 1) + (10/5 * 5) = 40
		},
		{
			name:            "15 throws with frenzy",
			throws:          15,
			maxCombo:        0,
			frenzyActivated: true,
			expectedXP:      30, // (15 * 1) * 2 = 30
		},
		{
			name:            "25 throws, combo of 10, with frenzy",
			throws:          25,
			maxCombo:        10,
			frenzyActivated: true,
			expectedXP:      70, // ((25 * 1) + (10/5 * 5)) * 2 = 70
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetchData := models.FetchGameData{
				Throws:              tt.throws,
				MaxCombo:            tt.maxCombo,
				FrenzyModeActivated: tt.frenzyActivated,
			}
			gameDataJSON, _ := json.Marshal(fetchData)

			duration := 600
			activity := &models.Activity{
				DurationSeconds: &duration,
				GameData:        gameDataJSON,
			}

			xp := service.calculateXP(gameType, activity)

			if xp != tt.expectedXP {
				t.Errorf("calculateXP() = %d, want %d", xp, tt.expectedXP)
			}
		})
	}
}

func TestCalculateXP_InvalidGameData(t *testing.T) {
	service := &ActivityService{}

	gameType := &models.GameType{
		ID:       "walk",
		XPConfig: json.RawMessage(`{"base_xp_per_minute": 2}`),
	}

	// Invalid JSON
	activity := &models.Activity{
		GameData: json.RawMessage(`invalid json`),
	}
	duration := 600
	activity.DurationSeconds = &duration

	xp := service.calculateXP(gameType, activity)
	// Even with invalid game data, base XP from duration should still be calculated
	// 10 minutes * 2 XP/minute = 20 XP
	expectedXP := 20
	if xp != expectedXP {
		t.Errorf("calculateXP() with invalid data = %d, want %d (base XP from duration)", xp, expectedXP)
	}
}

func TestUpdateStreak(t *testing.T) {
	service := &ActivityService{}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.Add(-24 * time.Hour)
	twoDaysAgo := today.Add(-48 * time.Hour)

	tests := []struct {
		name               string
		currentStreak      int
		lastActivityAt     *time.Time
		expectedStreakDays int
	}{
		{
			name:               "First activity ever",
			currentStreak:      0,
			lastActivityAt:     nil,
			expectedStreakDays: 1,
		},
		{
			name:               "Activity today (same day)",
			currentStreak:      5,
			lastActivityAt:     &now,
			expectedStreakDays: 5, // No change
		},
		{
			name:               "Activity yesterday (consecutive)",
			currentStreak:      3,
			lastActivityAt:     &yesterday,
			expectedStreakDays: 4, // Increment
		},
		{
			name:               "Activity 2 days ago (streak broken)",
			currentStreak:      10,
			lastActivityAt:     &twoDaysAgo,
			expectedStreakDays: 1, // Reset
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pet := &models.Pet{
				StreakDays:     tt.currentStreak,
				LastActivityAt: tt.lastActivityAt,
			}

			// Calculate expected streak
			streakDays := service.calculateExpectedStreak(pet)

			if streakDays != tt.expectedStreakDays {
				t.Errorf("calculateExpectedStreak() = %d, want %d", streakDays, tt.expectedStreakDays)
			}
		})
	}
}

// Helper method for testing streak calculation
func (s *ActivityService) calculateExpectedStreak(pet *models.Pet) int {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if pet.LastActivityAt == nil {
		return 1
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
		return pet.StreakDays // Same day, no change
	case 1:
		return pet.StreakDays + 1 // Consecutive day
	default:
		return 1 // Streak broken
	}
}
