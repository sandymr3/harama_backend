# HARaMA — Project Structure

## Overview Directory Tree

```
harrama/
├── backend/                      # Go backend services
│   ├── cmd/                      # Application entrypoints
│   │   ├── api/                  # Main API server
│   │   │   └── main.go
│   │   ├── worker/               # Background job processor
│   │   │   └── main.go
│   │   └── migrate/              # Database migrations runner
│   │       └── main.go
│   │
│   ├── internal/                 # Private application code
│   │   ├── api/                  # HTTP handlers & routing
│   │   │   ├── router.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── tenant.go
│   │   │   │   ├── ratelimit.go
│   │   │   │   └── logging.go
│   │   │   └── handlers/
│   │   │       ├── exam.go       # Exam CRUD endpoints
│   │   │       ├── submission.go # Upload & submission endpoints
│   │   │       ├── grading.go    # Grading endpoints
│   │   │       ├── override.go   # Teacher override endpoints
│   │   │       ├── analytics.go  # Dashboard & reports
│   │   │       └── health.go     # Health checks
│   │   │
│   │   ├── domain/               # Core business logic (entities)
│   │   │   ├── exam.go
│   │   │   ├── question.go
│   │   │   ├── rubric.go
│   │   │   ├── submission.go
│   │   │   ├── grade.go
│   │   │   ├── student.go
│   │   │   └── teacher.go
│   │   │
│   │   ├── service/              # Business logic layer
│   │   │   ├── exam_service.go
│   │   │   ├── grading_service.go
│   │   │   ├── ocr_service.go
│   │   │   ├── segmentation_service.go
│   │   │   ├── feedback_service.go
│   │   │   └── analytics_service.go
│   │   │
│   │   ├── grading/              # Core grading engine
│   │   │   ├── engine.go         # Main grading orchestrator
│   │   │   ├── evaluator.go      # Single evaluator logic
│   │   │   ├── multi_eval.go     # Multi-evaluator coordinator
│   │   │   ├── confidence.go     # Confidence calculation
│   │   │   ├── variance.go       # Variance & disagreement detection
│   │   │   ├── partial_credit.go # Partial credit engine
│   │   │   └── profiles/         # Subject-specific configs
│   │   │       ├── mathematics.go
│   │   │       ├── english.go
│   │   │       └── science.go
│   │   │
│   │   ├── ai/                   # AI provider abstraction
│   │   │   ├── provider.go       # Interface definition
│   │   │   ├── gemini/
│   │   │   │   ├── client.go     # Gemini API client
│   │   │   │   ├── grading.go    # Grading prompts
│   │   │   │   ├── vision.go     # Multimodal/diagram
│   │   │   │   ├── feedback.go   # Student feedback generation
│   │   │   │   └── prompts/      # Prompt templates
│   │   │   │       ├── base_grading.txt
│   │   │   │       ├── rubric_enforcer.txt
│   │   │   │       ├── reasoning_validator.txt
│   │   │   │       ├── structural_analyzer.txt
│   │   │   │       └── multimodal_grading.txt
│   │   │   └── cache.go          # Response caching
│   │   │
│   │   ├── ocr/                  # OCR processing
│   │   │   ├── processor.go      # OCR orchestrator
│   │   │   ├── google_vision.go  # Google Vision API
│   │   │   ├── tesseract.go      # Tesseract fallback
│   │   │   ├── confidence.go     # OCR confidence scoring
│   │   │   └── correction.go     # Gemini-based correction
│   │   │
│   │   ├── segmentation/         # Answer segmentation
│   │   │   ├── detector.go       # Question boundary detection
│   │   │   ├── diagram.go        # Diagram extraction
│   │   │   └── layout.go         # Spatial analysis
│   │   │
│   │   ├── storage/              # File & object storage
│   │   │   ├── minio.go          # MinIO client
│   │   │   ├── upload.go         # Upload handling
│   │   │   └── retrieval.go      # File retrieval
│   │   │
│   │   ├── repository/           # Data access layer
│   │   │   ├── postgres/
│   │   │   │   ├── exam_repo.go
│   │   │   │   ├── submission_repo.go
│   │   │   │   ├── grade_repo.go
│   │   │   │   ├── override_repo.go
│   │   │   │   ├── feedback_repo.go
│   │   │   │   └── audit_repo.go
│   │   │   └── qdrant/
│   │   │       └── vector_repo.go
│   │   │
│   │   ├── worker/               # Background job processing
│   │   │   ├── pool.go           # Worker pool
│   │   │   ├── jobs/
│   │   │   │   ├── ocr_job.go
│   │   │   │   ├── segmentation_job.go
│   │   │   │   ├── grading_job.go
│   │   │   │   └── feedback_job.go
│   │   │   └── queue.go          # Job queue (PostgreSQL-based)
│   │   │
│   │   ├── auth/                 # Authentication & authorization
│   │   │   ├── jwt.go
│   │   │   ├── permissions.go
│   │   │   └── tenant.go
│   │   │
│   │   ├── config/               # Configuration management
│   │   │   ├── config.go
│   │   │   └── validator.go
│   │   │
│   │   └── pkg/                  # Shared utilities
│   │       ├── logger/
│   │       │   └── logger.go
│   │       ├── errors/
│   │       │   └── errors.go
│   │       ├── validators/
│   │       │   └── validators.go
│   │       └── utils/
│   │           ├── image.go      # Image processing helpers
│   │           └── math.go       # Math utilities
│   │
│   ├── migrations/               # Database migrations
│   │   ├── 001_initial_schema.up.sql
│   │   ├── 001_initial_schema.down.sql
│   │   ├── 002_add_multimodal.up.sql
│   │   ├── 002_add_multimodal.down.sql
│   │   ├── 003_add_feedback.up.sql
│   │   └── 003_add_feedback.down.sql
│   │
│   ├── tests/                    # Tests
│   │   ├── unit/
│   │   │   ├── grading_test.go
│   │   │   ├── confidence_test.go
│   │   │   └── variance_test.go
│   │   ├── integration/
│   │   │   ├── api_test.go
│   │   │   └── grading_flow_test.go
│   │   └── fixtures/
│   │       ├── sample_exam.pdf
│   │       └── sample_submission.pdf
│   │
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   └── README.md
│
├── frontend/                     # React/Next.js frontend
│   ├── public/
│   │   ├── icons/
│   │   └── images/
│   │
│   ├── src/
│   │   ├── app/                  # Next.js App Router
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx          # Landing page
│   │   │   ├── dashboard/
│   │   │   │   └── page.tsx      # Teacher dashboard
│   │   │   ├── exams/
│   │   │   │   ├── page.tsx      # Exam list
│   │   │   │   ├── [id]/
│   │   │   │   │   └── page.tsx  # Exam detail
│   │   │   │   └── create/
│   │   │   │       └── page.tsx  # Create exam
│   │   │   ├── grading/
│   │   │   │   └── [id]/
│   │   │   │       └── page.tsx  # Grading interface
│   │   │   └── analytics/
│   │   │       └── page.tsx      # Analytics dashboard
│   │   │
│   │   ├── components/           # Reusable components
│   │   │   ├── ui/               # Base UI components
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   ├── Card.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   └── Badge.tsx
│   │   │   ├── exam/
│   │   │   │   ├── ExamCard.tsx
│   │   │   │   ├── QuestionEditor.tsx
│   │   │   │   └── RubricBuilder.tsx
│   │   │   ├── grading/
│   │   │   │   ├── GradingView.tsx
│   │   │   │   ├── AnswerDisplay.tsx
│   │   │   │   ├── AIReasoningPanel.tsx
│   │   │   │   ├── ConfidenceMeter.tsx
│   │   │   │   ├── VarianceIndicator.tsx
│   │   │   │   └── OverrideForm.tsx
│   │   │   ├── submission/
│   │   │   │   ├── UploadZone.tsx
│   │   │   │   └── SubmissionList.tsx
│   │   │   └── analytics/
│   │   │       ├── GradingTrends.tsx
│   │   │       └── QuestionDifficulty.tsx
│   │   │
│   │   ├── lib/                  # Utilities & helpers
│   │   │   ├── api.ts            # API client
│   │   │   ├── auth.ts           # Auth helpers
│   │   │   └── utils.ts          # Common utilities
│   │   │
│   │   ├── hooks/                # Custom React hooks
│   │   │   ├── useExams.ts
│   │   │   ├── useGrading.ts
│   │   │   └── useAuth.ts
│   │   │
│   │   ├── types/                # TypeScript types
│   │   │   ├── exam.ts
│   │   │   ├── grading.ts
│   │   │   └── api.ts
│   │   │
│   │   └── styles/
│   │       └── globals.css
│   │
│   ├── package.json
│   ├── tsconfig.json
│   ├── next.config.js
│   └── tailwind.config.js
│
├── docs/                         # Documentation
│   ├── API.md                    # API documentation
│   ├── ARCHITECTURE.md           # System architecture
│   ├── DEPLOYMENT.md             # Deployment guide
│   ├── GRADING_LOGIC.md          # Grading algorithm details
│   └── PROMPTS.md                # Prompt engineering guide
│
├── scripts/                      # Utility scripts
│   ├── seed_data.go              # Database seeding
│   ├── test_grading.sh           # Grading test script
│   └── deploy.sh                 # Deployment script
│
├── docker-compose.yml            # Local development stack
├── docker-compose.prod.yml       # Production stack
├── .env.example                  # Environment variables template
├── .gitignore
├── LICENSE
└── README.md                     # Project overview
```

---

## Key File Contents & Responsibilities

### Backend Core Files

#### `backend/cmd/api/main.go`
```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "harrama/internal/api"
    "harrama/internal/config"
    "harrama/internal/repository/postgres"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize database
    db, err := postgres.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Initialize router
    router := api.NewRouter(cfg, db)
    
    // Start server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }
    
    // Graceful shutdown
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
        <-sigint
        
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        if err := srv.Shutdown(ctx); err != nil {
            log.Printf("Server shutdown error: %v", err)
        }
    }()
    
    log.Printf("Server started on :%s", cfg.Port)
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Server error: %v", err)
    }
}
```

#### `backend/internal/domain/grade.go`
```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type GradingResult struct {
    ID              uuid.UUID       `json:"id"`
    SubmissionID    uuid.UUID       `json:"submission_id"`
    QuestionID      uuid.UUID       `json:"question_id"`
    Score           float64         `json:"score"`
    MaxScore        int             `json:"max_score"`
    Confidence      float64         `json:"confidence"`
    Reasoning       string          `json:"reasoning"`
    CriteriaMet     []string        `json:"criteria_met"`
    MistakesFound   []string        `json:"mistakes_found"`
    AIEvaluatorID   string          `json:"ai_evaluator_id"`
    CreatedAt       time.Time       `json:"created_at"`
}

type MultiEvalResult struct {
    Evaluations     []GradingResult `json:"evaluations"`
    Variance        float64         `json:"variance"`
    MeanScore       float64         `json:"mean_score"`
    ConsensusScore  float64         `json:"consensus_score"`
    ShouldEscalate  bool            `json:"should_escalate"`
}

type FinalGrade struct {
    ID              uuid.UUID       `json:"id"`
    SubmissionID    uuid.UUID       `json:"submission_id"`
    QuestionID      uuid.UUID       `json:"question_id"`
    FinalScore      float64         `json:"final_score"`
    AIScore         *float64        `json:"ai_score,omitempty"`
    OverrideScore   *float64        `json:"override_score,omitempty"`
    Confidence      float64         `json:"confidence"`
    Status          GradeStatus     `json:"status"`
    GradedBy        *uuid.UUID      `json:"graded_by,omitempty"`
    UpdatedAt       time.Time       `json:"updated_at"`
}

type GradeStatus string

const (
    GradeStatusPending     GradeStatus = "pending"
    GradeStatusAutoGraded  GradeStatus = "auto_graded"
    GradeStatusReview      GradeStatus = "needs_review"
    GradeStatusOverridden  GradeStatus = "overridden"
    GradeStatusFinal       GradeStatus = "final"
)
```

#### `backend/internal/grading/engine.go`
```go
package grading

import (
    "context"
    "fmt"
    
    "harrama/internal/domain"
    "harrama/internal/ai"
)

type Engine struct {
    aiProvider      ai.Provider
    confidenceCalc  *ConfidenceCalculator
    varianceCalc    *VarianceCalculator
}

func NewEngine(provider ai.Provider) *Engine {
    return &Engine{
        aiProvider:     provider,
        confidenceCalc: NewConfidenceCalculator(),
        varianceCalc:   NewVarianceCalculator(),
    }
}

func (e *Engine) GradeAnswer(ctx context.Context, answer domain.AnswerSegment, rubric domain.Rubric) (*domain.FinalGrade, error) {
    // Multi-evaluator grading
    multiEval, err := e.multiEvaluatorGrade(ctx, answer, rubric)
    if err != nil {
        return nil, fmt.Errorf("multi-evaluator grading failed: %w", err)
    }
    
    // Check if escalation needed
    if multiEval.ShouldEscalate {
        return &domain.FinalGrade{
            Status:     domain.GradeStatusReview,
            Confidence: multiEval.Variance,
        }, nil
    }
    
    // Build consensus grade
    finalGrade := e.buildConsensus(multiEval)
    finalGrade.Status = domain.GradeStatusAutoGraded
    
    return finalGrade, nil
}

func (e *Engine) multiEvaluatorGrade(ctx context.Context, answer domain.AnswerSegment, rubric domain.Rubric) (*domain.MultiEvalResult, error) {
    evaluators := []string{
        "rubric_enforcer",
        "reasoning_validator",
        "structural_analyzer",
    }
    
    results := make([]domain.GradingResult, len(evaluators))
    
    // Parallel evaluation
    for i, evalID := range evaluators {
        result, err := e.aiProvider.Grade(ctx, ai.GradingRequest{
            Answer:      answer,
            Rubric:      rubric,
            EvaluatorID: evalID,
        })
        if err != nil {
            return nil, err
        }
        results[i] = result
    }
    
    // Calculate variance
    variance := e.varianceCalc.Calculate(results)
    shouldEscalate := variance > 0.15 // 15% threshold
    
    return &domain.MultiEvalResult{
        Evaluations:    results,
        Variance:       variance,
        ShouldEscalate: shouldEscalate,
    }, nil
}
```

#### `backend/internal/ai/gemini/client.go`
```go
package gemini

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
    
    "harrama/internal/ai"
    "harrama/internal/domain"
)

type Client struct {
    client *genai.Client
    model  *genai.GenerativeModel
}

func NewClient(apiKey string) (*Client, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }
    
    model := client.GenerativeModel("gemini-3-pro")
    model.SetTemperature(0.2)
    model.SetTopK(40)
    model.SetTopP(0.95)
    
    return &Client{
        client: client,
        model:  model,
    }, nil
}

func (c *Client) Grade(ctx context.Context, req ai.GradingRequest) (domain.GradingResult, error) {
    // Load appropriate prompt template
    promptTemplate := loadPromptTemplate(req.EvaluatorID)
    
    // Build prompt
    prompt := buildGradingPrompt(promptTemplate, req.Answer, req.Rubric)
    
    // Call Gemini
    resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return domain.GradingResult{}, fmt.Errorf("gemini API error: %w", err)
    }
    
    // Parse structured response
    var result domain.GradingResult
    if err := parseResponse(resp, &result); err != nil {
        return domain.GradingResult{}, err
    }
    
    result.AIEvaluatorID = req.EvaluatorID
    return result, nil
}
```
---

### Frontend Core Files

#### `frontend/src/components/grading/GradingView.tsx`
```tsx
'use client'

import { useState, useEffect } from 'react'
import { AnswerDisplay } from './AnswerDisplay'
import { AIReasoningPanel } from './AIReasoningPanel'