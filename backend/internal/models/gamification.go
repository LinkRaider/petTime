package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Achievements

type Achievement struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Icon        *string         `json:"icon,omitempty"`
	Category    *string         `json:"category,omitempty"`
	Criteria    json.RawMessage `json:"criteria"`
	XPReward    int             `json:"xp_reward"`
}

type UserAchievement struct {
	UserID        uuid.UUID    `json:"user_id"`
	AchievementID string       `json:"achievement_id"`
	Achievement   *Achievement `json:"achievement,omitempty"`
	PetID         uuid.UUID    `json:"pet_id"`
	UnlockedAt    time.Time    `json:"unlocked_at"`
}

// Cards

type CardRarity string

const (
	CardRarityCommon    CardRarity = "common"
	CardRarityRare      CardRarity = "rare"
	CardRarityEpic      CardRarity = "epic"
	CardRarityLegendary CardRarity = "legendary"
)

type Card struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	ImageURL    *string         `json:"image_url,omitempty"`
	Rarity      CardRarity      `json:"rarity"`
	Category    *string         `json:"category,omitempty"`
	DropConfig  json.RawMessage `json:"drop_config,omitempty"`
}

type UserCard struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	CardID     string     `json:"card_id"`
	Card       *Card      `json:"card,omitempty"`
	ObtainedAt time.Time  `json:"obtained_at"`
	ActivityID *uuid.UUID `json:"activity_id,omitempty"`
}

// Missions

type MissionType string

const (
	MissionTypeWalkDuration  MissionType = "walk_duration"
	MissionTypeWalkDistance  MissionType = "walk_distance"
	MissionTypeActivityCount MissionType = "activity_count"
	MissionTypeFetchThrows   MissionType = "fetch_throws"
	MissionTypeExploreZones  MissionType = "explore_zones"
)

type Mission struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"user_id"`
	MissionType  MissionType `json:"mission_type"`
	Description  string      `json:"description"`
	TargetValue  int         `json:"target_value"`
	CurrentValue int         `json:"current_value"`
	XPReward     int         `json:"xp_reward"`
	ExpiresAt    time.Time   `json:"expires_at"`
	CompletedAt  *time.Time  `json:"completed_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}

func (m *Mission) IsCompleted() bool {
	return m.CurrentValue >= m.TargetValue
}

func (m *Mission) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

func (m *Mission) Progress() float64 {
	if m.TargetValue == 0 {
		return 0
	}
	progress := float64(m.CurrentValue) / float64(m.TargetValue)
	if progress > 1 {
		return 1
	}
	return progress
}
