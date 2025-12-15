package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joaosantos/pettime/internal/middleware"
	"github.com/joaosantos/pettime/internal/models"
	"github.com/joaosantos/pettime/internal/services"
)

type PetHandler struct {
	petService *services.PetService
}

func NewPetHandler(petService *services.PetService) *PetHandler {
	return &PetHandler{petService: petService}
}

type CreatePetRequest struct {
	PetTypeID string `json:"pet_type_id"`
	Name      string `json:"name"`
	Breed     string `json:"breed,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
}

type UpdatePetRequest struct {
	Name      *string `json:"name,omitempty"`
	Breed     *string `json:"breed,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
}

func (h *PetHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreatePetRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PetTypeID == "" || req.Name == "" {
		respondError(w, http.StatusBadRequest, "Pet type and name are required")
		return
	}

	input := models.CreatePetInput{
		PetTypeID: req.PetTypeID,
		Name:      req.Name,
	}

	if req.Breed != "" {
		input.Breed = &req.Breed
	}
	if req.AvatarURL != "" {
		input.AvatarURL = &req.AvatarURL
	}

	pet, err := h.petService.Create(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, services.ErrInvalidPetType) {
			respondError(w, http.StatusBadRequest, "Invalid pet type")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create pet")
		return
	}

	respondCreated(w, pet)
}

func (h *PetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	petID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	pet, err := h.petService.GetByID(r.Context(), userID, petID)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			respondError(w, http.StatusNotFound, "Pet not found")
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get pet")
		return
	}

	respondSuccess(w, pet)
}

func (h *PetHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pets, err := h.petService.GetByUserID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list pets")
		return
	}

	if pets == nil {
		pets = []*models.Pet{}
	}

	respondSuccess(w, pets)
}

func (h *PetHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	petID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	var req UpdatePetRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	input := models.UpdatePetInput{
		Name:      req.Name,
		Breed:     req.Breed,
		AvatarURL: req.AvatarURL,
	}

	pet, err := h.petService.Update(r.Context(), userID, petID, input)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			respondError(w, http.StatusNotFound, "Pet not found")
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update pet")
		return
	}

	respondSuccess(w, pet)
}

func (h *PetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	petID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	if err := h.petService.Delete(r.Context(), userID, petID); err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			respondError(w, http.StatusNotFound, "Pet not found")
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "Access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete pet")
		return
	}

	respondNoContent(w)
}

func (h *PetHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	petID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	pet, stats, err := h.petService.GetStats(r.Context(), userID, petID)
	if err != nil {
		if errors.Is(err, services.ErrPetNotFound) {
			respondError(w, http.StatusNotFound, "Pet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	respondSuccess(w, map[string]interface{}{
		"pet":   pet,
		"stats": stats,
	})
}

func (h *PetHandler) ListPetTypes(w http.ResponseWriter, r *http.Request) {
	petTypes, err := h.petService.GetAllPetTypes(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list pet types")
		return
	}

	respondSuccess(w, petTypes)
}
