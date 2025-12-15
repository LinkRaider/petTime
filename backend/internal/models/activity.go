package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GameType struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Description       *string         `json:"description,omitempty"`
	Icon              *string         `json:"icon,omitempty"`
	XPConfig          json.RawMessage `json:"xp_config,omitempty"`
	SupportedPetTypes []string        `json:"supported_pet_types"`
	Enabled           bool            `json:"enabled"`
}

type Activity struct {
	ID              uuid.UUID       `json:"id"`
	PetID           uuid.UUID       `json:"pet_id"`
	GameTypeID      string          `json:"game_type_id"`
	GameType        *GameType       `json:"game_type,omitempty"`
	StartedAt       time.Time       `json:"started_at"`
	EndedAt         *time.Time      `json:"ended_at,omitempty"`
	DurationSeconds *int            `json:"duration_seconds,omitempty"`
	XPEarned        int             `json:"xp_earned"`
	GameData        json.RawMessage `json:"game_data,omitempty"`
	ClientID        *uuid.UUID      `json:"client_id,omitempty"`
	SyncedAt        *time.Time      `json:"synced_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type CreateActivityInput struct {
	PetID      uuid.UUID       `json:"pet_id" validate:"required"`
	GameTypeID string          `json:"game_type_id" validate:"required"`
	StartedAt  time.Time       `json:"started_at" validate:"required"`
	EndedAt    *time.Time      `json:"ended_at,omitempty"`
	GameData   json.RawMessage `json:"game_data,omitempty"`
	ClientID   *uuid.UUID      `json:"client_id,omitempty"`
}

type UpdateActivityInput struct {
	EndedAt  *time.Time      `json:"ended_at,omitempty"`
	GameData json.RawMessage `json:"game_data,omitempty"`
}

type SyncActivityInput struct {
	ClientID   uuid.UUID       `json:"client_id" validate:"required"`
	PetID      uuid.UUID       `json:"pet_id" validate:"required"`
	GameTypeID string          `json:"game_type_id" validate:"required"`
	StartedAt  time.Time       `json:"started_at" validate:"required"`
	EndedAt    *time.Time      `json:"ended_at,omitempty"`
	GameData   json.RawMessage `json:"game_data,omitempty"`
}

type WalkGameData struct {
	DistanceMeters       float64     `json:"distance_meters"`
	Route                [][]float64 `json:"route,omitempty"`
	AvgSpeedKmh          float64     `json:"avg_speed_kmh,omitempty"`
	NewZonesDiscovered   []string    `json:"new_zones_discovered,omitempty"`
	Weather              string      `json:"weather,omitempty"`
}

type FetchGameData struct {
	Throws               int     `json:"throws"`
	Returns              int     `json:"returns"`
	SuccessRate          float64 `json:"success_rate"`
	MaxCombo             int     `json:"max_combo"`
	FrenzyModeActivated  bool    `json:"frenzy_mode_activated"`
}

type ActivityFilter struct {
	PetID      *uuid.UUID
	GameTypeID *string
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}
