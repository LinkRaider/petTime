package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaosantos/pettime/internal/models"
)

var ErrPetNotFound = errors.New("pet not found")

type PetRepository struct {
	db *pgxpool.Pool
}

func NewPetRepository(db *pgxpool.Pool) *PetRepository {
	return &PetRepository{db: db}
}

func (r *PetRepository) Create(ctx context.Context, pet *models.Pet) error {
	query := `
		INSERT INTO pets (id, user_id, pet_type_id, name, breed, avatar_url, birth_date, total_xp, level, mood, streak_days, last_activity_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Exec(ctx, query,
		pet.ID,
		pet.UserID,
		pet.PetTypeID,
		pet.Name,
		pet.Breed,
		pet.AvatarURL,
		pet.BirthDate,
		pet.TotalXP,
		pet.Level,
		pet.Mood,
		pet.StreakDays,
		pet.LastActivityAt,
		pet.CreatedAt,
		pet.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create pet: %w", err)
	}

	return nil
}

func (r *PetRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Pet, error) {
	query := `
		SELECT p.id, p.user_id, p.pet_type_id, p.name, p.breed, p.avatar_url, p.birth_date,
		       p.total_xp, p.level, p.mood, p.streak_days, p.last_activity_at, p.created_at, p.updated_at,
		       pt.id, pt.name, pt.icon, pt.config
		FROM pets p
		JOIN pet_types pt ON p.pet_type_id = pt.id
		WHERE p.id = $1
	`

	var pet models.Pet
	var petType models.PetType

	err := r.db.QueryRow(ctx, query, id).Scan(
		&pet.ID,
		&pet.UserID,
		&pet.PetTypeID,
		&pet.Name,
		&pet.Breed,
		&pet.AvatarURL,
		&pet.BirthDate,
		&pet.TotalXP,
		&pet.Level,
		&pet.Mood,
		&pet.StreakDays,
		&pet.LastActivityAt,
		&pet.CreatedAt,
		&pet.UpdatedAt,
		&petType.ID,
		&petType.Name,
		&petType.Icon,
		&petType.Config,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPetNotFound
		}
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}

	pet.PetType = &petType
	return &pet, nil
}

func (r *PetRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Pet, error) {
	query := `
		SELECT p.id, p.user_id, p.pet_type_id, p.name, p.breed, p.avatar_url, p.birth_date,
		       p.total_xp, p.level, p.mood, p.streak_days, p.last_activity_at, p.created_at, p.updated_at,
		       pt.id, pt.name, pt.icon, pt.config
		FROM pets p
		JOIN pet_types pt ON p.pet_type_id = pt.id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pets: %w", err)
	}
	defer rows.Close()

	var pets []*models.Pet
	for rows.Next() {
		var pet models.Pet
		var petType models.PetType

		err := rows.Scan(
			&pet.ID,
			&pet.UserID,
			&pet.PetTypeID,
			&pet.Name,
			&pet.Breed,
			&pet.AvatarURL,
			&pet.BirthDate,
			&pet.TotalXP,
			&pet.Level,
			&pet.Mood,
			&pet.StreakDays,
			&pet.LastActivityAt,
			&pet.CreatedAt,
			&pet.UpdatedAt,
			&petType.ID,
			&petType.Name,
			&petType.Icon,
			&petType.Config,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pet: %w", err)
		}

		pet.PetType = &petType
		pets = append(pets, &pet)
	}

	return pets, nil
}

func (r *PetRepository) Update(ctx context.Context, pet *models.Pet) error {
	query := `
		UPDATE pets
		SET name = $2, breed = $3, avatar_url = $4, birth_date = $5,
		    total_xp = $6, level = $7, mood = $8, streak_days = $9,
		    last_activity_at = $10, updated_at = $11
		WHERE id = $1
	`

	pet.UpdatedAt = time.Now()
	result, err := r.db.Exec(ctx, query,
		pet.ID,
		pet.Name,
		pet.Breed,
		pet.AvatarURL,
		pet.BirthDate,
		pet.TotalXP,
		pet.Level,
		pet.Mood,
		pet.StreakDays,
		pet.LastActivityAt,
		pet.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update pet: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPetNotFound
	}

	return nil
}

func (r *PetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM pets WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pet: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPetNotFound
	}

	return nil
}

func (r *PetRepository) AddXP(ctx context.Context, petID uuid.UUID, xp int) error {
	query := `
		UPDATE pets
		SET total_xp = total_xp + $2,
		    level = (
		        SELECT CASE
		            WHEN total_xp + $2 >= 8100 THEN 10
		            WHEN total_xp + $2 >= 6400 THEN 9
		            WHEN total_xp + $2 >= 4900 THEN 8
		            WHEN total_xp + $2 >= 3600 THEN 7
		            WHEN total_xp + $2 >= 2500 THEN 6
		            WHEN total_xp + $2 >= 1600 THEN 5
		            WHEN total_xp + $2 >= 900 THEN 4
		            WHEN total_xp + $2 >= 400 THEN 3
		            WHEN total_xp + $2 >= 100 THEN 2
		            ELSE 1
		        END
		    ),
		    last_activity_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, petID, xp)
	if err != nil {
		return fmt.Errorf("failed to add XP: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPetNotFound
	}

	return nil
}

func (r *PetRepository) UpdateStreak(ctx context.Context, petID uuid.UUID, streakDays int) error {
	query := `
		UPDATE pets
		SET streak_days = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, petID, streakDays)
	if err != nil {
		return fmt.Errorf("failed to update streak: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPetNotFound
	}

	return nil
}

func (r *PetRepository) UpdateMood(ctx context.Context, petID uuid.UUID, mood models.Mood) error {
	query := `
		UPDATE pets
		SET mood = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, petID, mood)
	if err != nil {
		return fmt.Errorf("failed to update mood: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPetNotFound
	}

	return nil
}

// Pet Types

func (r *PetRepository) GetAllPetTypes(ctx context.Context) ([]*models.PetType, error) {
	query := `SELECT id, name, icon, config FROM pet_types ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet types: %w", err)
	}
	defer rows.Close()

	var petTypes []*models.PetType
	for rows.Next() {
		var pt models.PetType
		if err := rows.Scan(&pt.ID, &pt.Name, &pt.Icon, &pt.Config); err != nil {
			return nil, fmt.Errorf("failed to scan pet type: %w", err)
		}
		petTypes = append(petTypes, &pt)
	}

	return petTypes, nil
}

func (r *PetRepository) GetPetType(ctx context.Context, id string) (*models.PetType, error) {
	query := `SELECT id, name, icon, config FROM pet_types WHERE id = $1`

	var pt models.PetType
	err := r.db.QueryRow(ctx, query, id).Scan(&pt.ID, &pt.Name, &pt.Icon, &pt.Config)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("pet type not found")
		}
		return nil, fmt.Errorf("failed to get pet type: %w", err)
	}

	return &pt, nil
}
