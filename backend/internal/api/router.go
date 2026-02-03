package api

import (
	"net/http"

	"harama/internal/ai/gemini"
	"harama/internal/api/handlers"
	"harama/internal/api/middleware"
	"harama/internal/config"
	"harama/internal/grading"
	"harama/internal/repository/postgres"
	"harama/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func NewRouter(cfg *config.Config, db *bun.DB) *chi.Mux {
	r := chi.NewRouter()

	// 1. Initialize Repositories
	examRepo := postgres.NewExamRepo(db)
	subRepo := postgres.NewSubmissionRepo(db)
	gradeRepo := postgres.NewGradeRepo(db)
	feedbackRepo := postgres.NewFeedbackRepo(db)

	// 2. Initialize AI Provider
	aiClient, _ := gemini.NewClient(cfg.GeminiAPIKey)

	// 3. Initialize Engine
	gradingEngine := grading.NewEngine(aiClient)

	// 4. Initialize Services
	examService := service.NewExamService(examRepo)
	ocrService := service.NewOCRService(subRepo)
	gradingService := service.NewGradingService(gradeRepo, examRepo, subRepo, gradingEngine)
	feedbackService := service.NewFeedbackService(feedbackRepo, gradeRepo, examRepo, aiClient)

	// 5. Initialize Handlers
	examHandler := handlers.NewExamHandler(examService)
	submissionHandler := handlers.NewSubmissionHandler(ocrService, gradingService)
	gradingHandler := handlers.NewGradingHandler(gradingService)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)

	// 6. Global Middleware
	r.Use(middleware.TenantMiddleware)

	// 7. Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Exam Routes
		r.Post("/exams", examHandler.CreateExam)
		r.Get("/exams/{id}", examHandler.GetExam)
		r.Post("/exams/{id}/questions", examHandler.AddQuestion)
		r.Put("/questions/{id}/rubric", examHandler.SetRubric)

		// Submission Routes
		r.Post("/exams/{id}/submissions", submissionHandler.CreateSubmission)
		r.Post("/submissions/{id}/trigger-grading", submissionHandler.TriggerGrading)

		// Grading & Feedback Routes
		r.Get("/submissions/{id}/grades", gradingHandler.GetGrades)
		r.Post("/submissions/{submission_id}/questions/{question_id}/override", feedbackHandler.CaptureOverride)
		r.Get("/submissions/{submission_id}/questions/{question_id}/feedback", feedbackHandler.GetStudentFeedback)
		r.Get("/questions/{question_id}/analysis", feedbackHandler.AnalyzePatterns)
	})

	return r
}