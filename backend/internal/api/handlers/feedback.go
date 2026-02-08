package handlers

import (
	"encoding/json"
	"net/http"

	"harama/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	service *service.FeedbackService
}

func NewFeedbackHandler(s *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{service: s}
}

func (h *FeedbackHandler) CaptureOverride(w http.ResponseWriter, r *http.Request) {
	subID, _ := uuid.Parse(chi.URLParam(r, "submission_id"))
	questionID, _ := uuid.Parse(chi.URLParam(r, "question_id"))

	var body struct {
		Score    float64 `json:"score"`
		NewScore float64 `json:"new_score"`
		Reason   string  `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Support both "score" and "new_score" field names from frontend
	score := body.Score
	if body.NewScore != 0 {
		score = body.NewScore
	}

	err := h.service.CaptureOverrideFeedback(r.Context(), subID, questionID, score, body.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FeedbackHandler) GetStudentFeedback(w http.ResponseWriter, r *http.Request) {
	subID, _ := uuid.Parse(chi.URLParam(r, "submission_id"))
	questionID, _ := uuid.Parse(chi.URLParam(r, "question_id"))
	studentName := r.URL.Query().Get("name")

	feedback, err := h.service.GenerateStudentFeedback(r.Context(), subID, questionID, studentName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"feedback": feedback})
}

func (h *FeedbackHandler) AnalyzePatterns(w http.ResponseWriter, r *http.Request) {
	questionID, _ := uuid.Parse(chi.URLParam(r, "question_id"))

	result, err := h.service.AnalyzeQuestionPatterns(r.Context(), questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *FeedbackHandler) AdaptRubric(w http.ResponseWriter, r *http.Request) {
	questionID, _ := uuid.Parse(chi.URLParam(r, "question_id"))

	err := h.service.AdaptRubric(r.Context(), questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
