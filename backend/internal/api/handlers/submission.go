package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"harama/internal/domain"
	"harama/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubmissionHandler struct {
	ocrService     *service.OCRService
	gradingService *service.GradingService
}

func NewSubmissionHandler(ocr *service.OCRService, grading *service.GradingService) *SubmissionHandler {
	return &SubmissionHandler{
		ocrService:     ocr,
		gradingService: grading,
	}
}

func (h *SubmissionHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	examIDStr := chi.URLParam(r, "id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		http.Error(w, "invalid exam id", http.StatusBadRequest)
		return
	}

	// This is where file upload handling would go.
	// For now, we assume metadata is posted.
	var sub domain.Submission
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sub.ExamID = examID
	// Hardcode tenant for now
	sub.TenantID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

	// Trigger OCR processing (async in real app)
	if err := h.ocrService.ProcessSubmission(r.Context(), sub.ID); err != nil {
		// Log error but don't fail request? For now, fail.
		http.Error(w, "failed to start processing: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(sub)
}

func (h *SubmissionHandler) TriggerGrading(w http.ResponseWriter, r *http.Request) {
	subIDStr := chi.URLParam(r, "id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		http.Error(w, "invalid submission id", http.StatusBadRequest)
		return
	}

	// Async in real world
	go func() {
		// Create a detached context or use background
		h.gradingService.GradeSubmission(context.Background(), subID)
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "grading_started"}`))
}

func (h *SubmissionHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	subIDStr := chi.URLParam(r, "id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		http.Error(w, "invalid submission id", http.StatusBadRequest)
		return
	}

	sub, err := h.ocrService.GetByID(r.Context(), subID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}
