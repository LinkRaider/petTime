package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PetType struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Icon   *string         `json:"icon,omitempty"`
	Config json.RawMessage `json:"config,omitempty"`
}

type Mood string

const (
	MoodHappy     Mood = "happy"
	MoodContent   Mood = "content"
	MoodTired     Mood = "tired"
	MoodSad       Mood = "sad"
	MoodExcited   Mood = "excited"
	MoodBored     Mood = "bored"
)

type Pet struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	PetTypeID      string     `json:"pet_type_id"`
	PetType        *PetType   `json:"pet_type,omitempty"`
	Name           string     `json:"name"`
	Breed          *string    `json:"breed,omitempty"`
	AvatarURL      *string    `json:"avatar_url,omitempty"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	TotalXP        int        `json:"total_xp"`
	Level          int        `json:"level"`
	Mood           Mood       `json:"mood"`
	StreakDays     int        `json:"streak_days"`
	LastActivityAt *time.Time `json:"last_activity_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreatePetInput struct {
	PetTypeID string     `json:"pet_type_id" validate:"required"`
	Name      string     `json:"name" validate:"required,min=1,max=100"`
	Breed     *string    `json:"breed,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
}

type UpdatePetInput struct {
	Name      *string    `json:"name,omitempty"`
	Breed     *string    `json:"breed,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
}

type PetStats struct {
	TotalActivities   int     `json:"total_activities"`
	TotalDuration     int     `json:"total_duration_seconds"`
	TotalDistance     float64 `json:"total_distance_meters"`
	CurrentStreak     int     `json:"current_streak"`
	LongestStreak     int     `json:"longest_streak"`
	XPToNextLevel     int     `json:"xp_to_next_level"`
	LevelProgress     float64 `json:"level_progress"`
}

func CalculateLevel(xp int) int {
	// Level formula: level = floor(sqrt(xp / 100)) + 1
	// Level 1: 0-99 XP
	// Level 2: 100-399 XP
	// Level 3: 400-899 XP
	// etc.
	level := 1
	threshold := 100
	for xp >= threshold {
		level++
		threshold = level * level * 100
	}
	return level
}

func XPForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	return (level - 1) * (level - 1) * 100
}

func XPToNextLevel(currentXP int) int {
	currentLevel := CalculateLevel(currentXP)
	nextLevelXP := XPForLevel(currentLevel + 1)
	return nextLevelXP - currentXP
}
