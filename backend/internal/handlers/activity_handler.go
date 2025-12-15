package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joaosantos/pettime/internal/middleware"
	"github.com/joaosantos/pettime/internal/models"
	"github.com/joaosantos/pettime/internal/services"
)

type ActivityHandler struct {
	activityService *services.ActivityService
}

func NewActivityHandler(activityService *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

type CreateActivityRequest struct {
	PetID      string          `json:"pet_id"`
	GameTypeID string          `json:"game_type_id"`
	StartedAt  string          `json:"started_at"`
	EndedAt    *string         `json:"ended_at,omitempty"`
	GameData   json.RawMessage `json:"game_data,omitempty"`
	ClientID   *string         `json:"client_id,omitempty"`
}

type UpdateActivityRequest struct {
	EndedAt  *string         `json:"ended_at,omitempty"`
	GameData json.RawMessage `json:"game_data,omitempty"`
}

type SyncActivityRequest struct {
	ClientID   string          `json:"client_id"`
	PetID      string          `json:"pet_id"`
	GameTypeID string          `json:"game_type_id"`
	StartedAt  string          `json:"started_at"`
	EndedAt    *string         `json:"ended_at,omitempty"`
	GameData   json.RawMessage `json:"game_data,omitempty"`
}

type SyncRequest struct {
	Activities []SyncActivityRequest `json:"activities"`
}

func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateActivityRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	petID, err := uuid.Parse(req.PetID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	startedAt, err := time.Parse(time.RFC3339, req.StartedAt)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid started_at format")
		return
	}

	input := models.CreateActivityInput{
		PetID:      petID,
		GameTypeID: req.GameTypeID,
		StartedAt:  startedAt,
		GameData:   req.GameData,
	}

	if req.EndedAt != nil {
		endedAt, err := time.Parse(time.RFC3339, *req.EndedAt)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid ended_at format")
			return
		}
		input.EndedAt = &endedAt
	}

	if req.ClientID != nil {
		clientID, err := uuid.Parse(*req.ClientID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid client ID")
			return
		}
		input.ClientID = &clientID
	}

	activity, err := h.activityService.Create(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			respondError(w, http.StatusNotFound, "Pet not found")
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		if errors.Is(err, services.ErrInvalidGameType) {
			respondError(w, http.StatusBadRequest, "Invalid game type")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create activity")
		return
	}

	respondCreated(w, activity)
}

func (h *ActivityHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	activityID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	activity, err := h.activityService.GetByID(r.Context(), userID, activityID)
	if err != nil {
		if errors.Is(err, services.ErrActivityNotFound) {
			respondError(w, http.StatusNotFound, "Activity not found")
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get activity")
		return
	}

	respondSuccess(w, activity)
}

func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	filter := models.ActivityFilter{
		Limit:  50,
		Offset: 0,
	}

	if petIDStr := r.URL.Query().Get("pet_id"); petIDStr != "" {
		petID, err := uuid.Parse(petIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid pet ID")
			return
		}
		filter.PetID = &petID
	}

	if gameTypeID := r.URL.Query().Get("game_type_id"); gameTypeID != "" {
		filter.GameTypeID = &gameTypeID
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	activities, err := h.activityService.List(r.Context(), userID, filter)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to list activities")
		return
	}

	if activities == nil {
		activities = []*models.Activity{}
	}

	respondSuccess(w, activities)
}

func (h *ActivityHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	activityID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	var req UpdateActivityRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	input := models.UpdateActivityInput{
		GameData: req.GameData,
	}

	if req.EndedAt != nil {
		endedAt, err := time.Parse(time.RFC3339, *req.EndedAt)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid ended_at format")
			return
		}
		input.EndedAt = &endedAt
	}

	activity, err := h.activityService.Update(r.Context(), userID, activityID, input)
	if err != nil {
		if errors.Is(err, services.ErrActivityNotFound) {
			respondError(w, http.StatusNotFound, "Activity not found")
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update activity")
		return
	}

	respondSuccess(w, activity)
}

func (h *ActivityHandler) Sync(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req SyncRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var inputs []models.SyncActivityInput
	for _, act := range req.Activities {
		clientID, err := uuid.Parse(act.ClientID)
		if err != nil {
			continue
		}

		petID, err := uuid.Parse(act.PetID)
		if err != nil {
			continue
		}

		startedAt, err := time.Parse(time.RFC3339, act.StartedAt)
		if err != nil {
			continue
		}

		input := models.SyncActivityInput{
			ClientID:   clientID,
			PetID:      petID,
			GameTypeID: act.GameTypeID,
			StartedAt:  startedAt,
			GameData:   act.GameData,
		}

		if act.EndedAt != nil {
			endedAt, err := time.Parse(time.RFC3339, *act.EndedAt)
			if err == nil {
				input.EndedAt = &endedAt
			}
		}

		inputs = append(inputs, input)
	}

	activities, err := h.activityService.Sync(r.Context(), userID, inputs)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to sync activities")
		return
	}

	respondSuccess(w, activities)
}

func (h *ActivityHandler) ListGameTypes(w http.ResponseWriter, r *http.Request) {
	gameTypes, err := h.activityService.GetAllGameTypes(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list game types")
		return
	}

	respondSuccess(w, gameTypes)
}
