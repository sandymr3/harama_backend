package api

import (
	"fmt"
	"net/http"

	"harama/internal/ai/gemini"
	"harama/internal/api/handlers"
	"harama/internal/api/middleware"
	"harama/internal/config"
	"harama/internal/grading"
	"harama/internal/ocr"
	"harama/internal/repository/postgres"
	"harama/internal/service"
	"harama/internal/storage"
	"harama/internal/worker"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

func NewRouter(cfg *config.Config, db *bun.DB) (*chi.Mux, error) {
	r := chi.NewRouter()

	// 1. Initialize Repositories
	examRepo := postgres.NewExamRepo(db)
	subRepo := postgres.NewSubmissionRepo(db)
	gradeRepo := postgres.NewGradeRepo(db)
	feedbackRepo := postgres.NewFeedbackRepo(db)
	auditRepo := postgres.NewAuditRepo(db)

	// 2. Initialize AI Provider & Infrastructure
	aiClient, err := gemini.NewClient(cfg.GeminiAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gemini client: %w", err)
	}
	
	minioStorage, err := storage.NewMinioStorage(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio storage: %w", err)
	}

	visionProcessor, err := ocr.NewGeminiOCRProcessor(cfg.GeminiAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gemini vision processor: %w", err)
	}

	// 3. Initialize Engine & Worker Pool
	gradingEngine := grading.NewEngine(aiClient)
	workerPool := worker.NewWorkerPool(5, 100)
	workerPool.Start()

	// 4. Initialize Services
	examService := service.NewExamService(examRepo, auditRepo)
	ocrService := service.NewOCRService(subRepo, auditRepo, minioStorage, visionProcessor)
	gradingService := service.NewGradingService(gradeRepo, examRepo, subRepo, auditRepo, gradingEngine)
	feedbackService := service.NewFeedbackService(feedbackRepo, gradeRepo, examRepo, auditRepo, aiClient)

	// 5. Initialize Handlers
	examHandler := handlers.NewExamHandler(examService)
	submissionHandler := handlers.NewSubmissionHandler(ocrService, gradingService, workerPool)
	gradingHandler := handlers.NewGradingHandler(gradingService)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)

	// 6. Global Middleware
	r.Use(middleware.RateLimitMiddleware(middleware.NewIPRateLimiter(50, 100)))

	// 7. Unprotected Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 8. Protected API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.TenantMiddleware)

		// Exam Routes
		r.Post("/exams", examHandler.CreateExam)
		r.Get("/exams/{id}", examHandler.GetExam)
		r.Post("/exams/{id}/questions", examHandler.AddQuestion)
		r.Put("/questions/{id}/rubric", examHandler.SetRubric)

		// Submission Routes
		r.Post("/exams/{id}/submissions", submissionHandler.CreateSubmission)
		r.Get("/submissions/{id}", submissionHandler.GetSubmission)
		r.Post("/submissions/{id}/trigger-grading", submissionHandler.TriggerGrading)

		// Grading & Feedback Routes
		r.Get("/submissions/{id}/grades", gradingHandler.GetGrades)
		r.Post("/submissions/{submission_id}/questions/{question_id}/override", feedbackHandler.CaptureOverride)
		r.Get("/submissions/{submission_id}/questions/{question_id}/feedback", feedbackHandler.GetStudentFeedback)
		r.Get("/questions/{question_id}/analysis", feedbackHandler.AnalyzePatterns)
	})

	

		return r, nil

	}

	