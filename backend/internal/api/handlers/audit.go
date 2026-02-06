package handlers

import (
	"encoding/json"
	"net/http"

	"harama/internal/service"

	"github.com/go-chi/chi/v5"
)

type AuditHandler struct {
	service *service.AuditService
}

func NewAuditHandler(s *service.AuditService) *AuditHandler {
	return &AuditHandler{service: s}
}

func (h *AuditHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("type")
	entityID := chi.URLParam(r, "id")

	if entityType == "" || entityID == "" {
		http.Error(w, "entity type and id are required", http.StatusBadRequest)
		return
	}

	logs, err := h.service.GetLogsForEntity(r.Context(), entityType, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
