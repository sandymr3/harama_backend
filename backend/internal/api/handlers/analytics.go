package handlers

import (
	"encoding/json"
	"net/http"

	"harama/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(s *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: s}
}

func (h *AnalyticsHandler) GetGradingTrends(w http.ResponseWriter, r *http.Request) {
	examIDStr := r.URL.Query().Get("exam_id")
	if examIDStr == "" {
		http.Error(w, "exam_id query parameter is required", http.StatusBadRequest)
		return
	}
	
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		http.Error(w, "invalid exam_id", http.StatusBadRequest)
		return
	}

	trends, err := h.service.GetGradingTrends(r.Context(), examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trends)
}

func (h *AnalyticsHandler) ExportGrades(w http.ResponseWriter, r *http.Request) {
	examIDStr := chi.URLParam(r, "id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		http.Error(w, "invalid exam id", http.StatusBadRequest)
		return
	}

	var req struct {
		Format string `json:"format"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Default to CSV if body is empty or invalid?
		req.Format = "csv"
	}
	if req.Format == "" {
		req.Format = "csv"
	}

	data, contentType, err := h.service.ExportGrades(r.Context(), examID, req.Format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=grades.csv")
	w.Write(data)
}