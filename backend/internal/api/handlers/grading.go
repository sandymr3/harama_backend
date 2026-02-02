package handlers

import (
	"encoding/json"
	"net/http"

	"harama/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GradingHandler struct {
	service *service.GradingService
}

func NewGradingHandler(s *service.GradingService) *GradingHandler {
	return &GradingHandler{service: s}
}

// GetGrades returns all grades for a submission
func (h *GradingHandler) GetGrades(w http.ResponseWriter, r *http.Request) {
	subIDStr := chi.URLParam(r, "id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		http.Error(w, "invalid submission id", http.StatusBadRequest)
		return
	}
	
	grades, err := h.service.GetGrades(r.Context(), subID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}
