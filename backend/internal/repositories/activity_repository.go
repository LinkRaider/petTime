package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaosantos/pettime/internal/models"
)

var ErrActivityNotFound = errors.New("activity not found")

type ActivityRepository struct {
	db *pgxpool.Pool
}

func NewActivityRepository(db *pgxpool.Pool) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(ctx context.Context, activity *models.Activity) error {
	query := `
		INSERT INTO activities (id, pet_id, game_type_id, started_at, ended_at, duration_seconds, xp_earned, game_data, client_id, synced_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		activity.ID,
		activity.PetID,
		activity.GameTypeID,
		activity.StartedAt,
		activity.EndedAt,
		activity.DurationSeconds,
		activity.XPEarned,
		activity.GameData,
		activity.ClientID,
		activity.SyncedAt,
		activity.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	return nil
}

func (r *ActivityRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Activity, error) {
	query := `
		SELECT a.id, a.pet_id, a.game_type_id, a.started_at, a.ended_at, a.duration_seconds,
		       a.xp_earned, a.game_data, a.client_id, a.synced_at, a.created_at,
		       gt.id, gt.name, gt.description, gt.icon, gt.xp_config, gt.supported_pet_types, gt.enabled
		FROM activities a
		JOIN game_types gt ON a.game_type_id = gt.id
		WHERE a.id = $1
	`

	var activity models.Activity
	var gameType models.GameType

	err := r.db.QueryRow(ctx, query, id).Scan(
		&activity.ID,
		&activity.PetID,
		&activity.GameTypeID,
		&activity.StartedAt,
		&activity.EndedAt,
		&activity.DurationSeconds,
		&activity.XPEarned,
		&activity.GameData,
		&activity.ClientID,
		&activity.SyncedAt,
		&activity.CreatedAt,
		&gameType.ID,
		&gameType.Name,
		&gameType.Description,
		&gameType.Icon,
		&gameType.XPConfig,
		&gameType.SupportedPetTypes,
		&gameType.Enabled,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrActivityNotFound
		}
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	activity.GameType = &gameType
	return &activity, nil
}

func (r *ActivityRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) (*models.Activity, error) {
	query := `
		SELECT a.id, a.pet_id, a.game_type_id, a.started_at, a.ended_at, a.duration_seconds,
		       a.xp_earned, a.game_data, a.client_id, a.synced_at, a.created_at
		FROM activities a
		WHERE a.client_id = $1
	`

	var activity models.Activity
	err := r.db.QueryRow(ctx, query, clientID).Scan(
		&activity.ID,
		&activity.PetID,
		&activity.GameTypeID,
		&activity.StartedAt,
		&activity.EndedAt,
		&activity.DurationSeconds,
		&activity.XPEarned,
		&activity.GameData,
		&activity.ClientID,
		&activity.SyncedAt,
		&activity.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrActivityNotFound
		}
		return nil, fmt.Errorf("failed to get activity by client ID: %w", err)
	}

	return &activity, nil
}

func (r *ActivityRepository) List(ctx context.Context, filter models.ActivityFilter) ([]*models.Activity, error) {
	query := `
		SELECT a.id, a.pet_id, a.game_type_id, a.started_at, a.ended_at, a.duration_seconds,
		       a.xp_earned, a.game_data, a.client_id, a.synced_at, a.created_at,
		       gt.id, gt.name, gt.description, gt.icon, gt.xp_config, gt.supported_pet_types, gt.enabled
		FROM activities a
		JOIN game_types gt ON a.game_type_id = gt.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filter.PetID != nil {
		query += fmt.Sprintf(" AND a.pet_id = $%d", argIndex)
		args = append(args, *filter.PetID)
		argIndex++
	}

	if filter.GameTypeID != nil {
		query += fmt.Sprintf(" AND a.game_type_id = $%d", argIndex)
		args = append(args, *filter.GameTypeID)
		argIndex++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND a.started_at >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND a.started_at <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	query += " ORDER BY a.started_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list activities: %w", err)
	}
	defer rows.Close()

	var activities []*models.Activity
	for rows.Next() {
		var activity models.Activity
		var gameType models.GameType

		err := rows.Scan(
			&activity.ID,
			&activity.PetID,
			&activity.GameTypeID,
			&activity.StartedAt,
			&activity.EndedAt,
			&activity.DurationSeconds,
			&activity.XPEarned,
			&activity.GameData,
			&activity.ClientID,
			&activity.SyncedAt,
			&activity.CreatedAt,
			&gameType.ID,
			&gameType.Name,
			&gameType.Description,
			&gameType.Icon,
			&gameType.XPConfig,
			&gameType.SupportedPetTypes,
			&gameType.Enabled,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		activity.GameType = &gameType
		activities = append(activities, &activity)
	}

	return activities, nil
}

func (r *ActivityRepository) Update(ctx context.Context, activity *models.Activity) error {
	query := `
		UPDATE activities
		SET ended_at = $2, duration_seconds = $3, xp_earned = $4, game_data = $5, synced_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		activity.ID,
		activity.EndedAt,
		activity.DurationSeconds,
		activity.XPEarned,
		activity.GameData,
		activity.SyncedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrActivityNotFound
	}

	return nil
}

func (r *ActivityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM activities WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrActivityNotFound
	}

	return nil
}

// Game Types

func (r *ActivityRepository) GetAllGameTypes(ctx context.Context) ([]*models.GameType, error) {
	query := `
		SELECT id, name, description, icon, xp_config, supported_pet_types, enabled
		FROM game_types
		WHERE enabled = true
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get game types: %w", err)
	}
	defer rows.Close()

	var gameTypes []*models.GameType
	for rows.Next() {
		var gt models.GameType
		if err := rows.Scan(&gt.ID, &gt.Name, &gt.Description, &gt.Icon, &gt.XPConfig, &gt.SupportedPetTypes, &gt.Enabled); err != nil {
			return nil, fmt.Errorf("failed to scan game type: %w", err)
		}
		gameTypes = append(gameTypes, &gt)
	}

	return gameTypes, nil
}

func (r *ActivityRepository) GetGameType(ctx context.Context, id string) (*models.GameType, error) {
	query := `
		SELECT id, name, description, icon, xp_config, supported_pet_types, enabled
		FROM game_types
		WHERE id = $1
	`

	var gt models.GameType
	err := r.db.QueryRow(ctx, query, id).Scan(
		&gt.ID, &gt.Name, &gt.Description, &gt.Icon, &gt.XPConfig, &gt.SupportedPetTypes, &gt.Enabled,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("game type not found")
		}
		return nil, fmt.Errorf("failed to get game type: %w", err)
	}

	return &gt, nil
}

// Stats

func (r *ActivityRepository) GetPetStats(ctx context.Context, petID uuid.UUID) (*models.PetStats, error) {
	query := `
		SELECT
			COUNT(*) as total_activities,
			COALESCE(SUM(duration_seconds), 0) as total_duration,
			COALESCE(SUM((game_data->>'distance_meters')::float), 0) as total_distance
		FROM activities
		WHERE pet_id = $1 AND ended_at IS NOT NULL
	`

	var stats models.PetStats
	err := r.db.QueryRow(ctx, query, petID).Scan(
		&stats.TotalActivities,
		&stats.TotalDuration,
		&stats.TotalDistance,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet stats: %w", err)
	}

	return &stats, nil
}
