package handlers

import (
	"encoding/json"
	"net/http"

	"harama/internal/auth"
	"harama/internal/domain"
	"harama/internal/service"
	"harama/internal/worker"
	"harama/internal/worker/jobs"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubmissionHandler struct {
	ocrService     *service.OCRService
	gradingService *service.GradingService
	workerPool     *worker.WorkerPool
}

func NewSubmissionHandler(ocr *service.OCRService, grading *service.GradingService, pool *worker.WorkerPool) *SubmissionHandler {
	return &SubmissionHandler{
		ocrService:     ocr,
		gradingService: grading,
		workerPool:     pool,
	}
}

func (h *SubmissionHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	tenantID, err := auth.GetTenantID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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
	sub.TenantID = tenantID

	if err := h.ocrService.CreateSubmission(r.Context(), &sub); err != nil {
		http.Error(w, "failed to create submission: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Trigger OCR processing asynchronously using worker pool
	h.workerPool.Submit(&jobs.OCRJob{
		SubmissionID: sub.ID,
		Service:      h.ocrService,
	})

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

	// Submit grading job to worker pool
	h.workerPool.Submit(&jobs.GradingJob{
		SubmissionID: subID,
		Service:      h.gradingService,
	})

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
