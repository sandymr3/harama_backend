package handlers

import (
	"encoding/json"
	"net/http"

	"harama/internal/domain"
	"harama/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ExamHandler struct {
	service *service.ExamService
}

func NewExamHandler(s *service.ExamService) *ExamHandler {
	return &ExamHandler{service: s}
}

func (h *ExamHandler) CreateExam(w http.ResponseWriter, r *http.Request) {
	var exam domain.Exam
	if err := json.NewDecoder(r.Body).Decode(&exam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For now, hardcode a tenant ID
	exam.TenantID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

	if err := h.service.CreateExam(r.Context(), &exam); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exam)
}

func (h *ExamHandler) GetExam(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid exam id", http.StatusBadRequest)
		return
	}

	exam, err := h.service.GetExam(r.Context(), id)
	if err != nil {
		http.Error(w, "exam not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exam)
}

func (h *ExamHandler) AddQuestion(w http.ResponseWriter, r *http.Request) {
	examIDStr := chi.URLParam(r, "id")
	examID, err := uuid.Parse(examIDStr)
	if err != nil {
		http.Error(w, "invalid exam id", http.StatusBadRequest)
		return
	}

	var question domain.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.AddQuestion(r.Context(), examID, &question); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

func (h *ExamHandler) SetRubric(w http.ResponseWriter, r *http.Request) {
	questionIDStr := chi.URLParam(r, "id")
	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	var rubric domain.Rubric
	if err := json.NewDecoder(r.Body).Decode(&rubric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SetRubric(r.Context(), questionID, &rubric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rubric)
}
