package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GradingResult struct {
	ID            uuid.UUID `json:"id"`
	SubmissionID  uuid.UUID `json:"submission_id"`
	QuestionID    uuid.UUID `json:"question_id"`
	Score         float64   `json:"score"`
	MaxScore      int       `json:"max_score"`
	Confidence    float64   `json:"confidence"`
	Reasoning     string    `json:"reasoning"`
	CriteriaMet   []string  `json:"criteria_met"`
	MistakesFound []string  `json:"mistakes_found"`
	AIEvaluatorID string    `json:"ai_evaluator_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type MultiEvalResult struct {
	Evaluations    []GradingResult `json:"evaluations"`
	Variance       float64         `json:"variance"`
	MeanScore      float64         `json:"mean_score"`
	ConsensusScore float64         `json:"consensus_score"`
	Confidence     float64         `json:"confidence"`
	Reasoning      string          `json:"reasoning"`
	ShouldEscalate bool            `json:"should_escalate"`
}

type FinalGrade struct {
	bun.BaseModel `bun:"table:grades,alias:g"`

	ID            uuid.UUID   `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	SubmissionID  uuid.UUID   `bun:"submission_id,notnull,type:uuid" json:"submission_id"`
	QuestionID    uuid.UUID   `bun:"question_id,notnull,type:uuid" json:"question_id"`
	FinalScore    float64     `bun:"final_score,notnull" json:"final_score"`
	MaxScore      int         `bun:"max_score,notnull" json:"max_score"`
	AIScore       *float64    `bun:"ai_score" json:"ai_score,omitempty"`
	OverrideScore *float64    `bun:"override_score" json:"override_score,omitempty"`
	Confidence    float64     `bun:"confidence,notnull" json:"confidence"`
	Reasoning     string      `bun:"reasoning" json:"reasoning"`
	Status        GradeStatus `bun:"status,notnull" json:"status"`
	GradedBy      *uuid.UUID  `bun:"graded_by,type:uuid" json:"graded_by,omitempty"`
	UpdatedAt     time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

type GradeStatus string

const (
    GradeStatusPending     GradeStatus = "pending"
    GradeStatusAutoGraded  GradeStatus = "auto_graded"
    GradeStatusReview      GradeStatus = "needs_review"
    GradeStatusOverridden  GradeStatus = "overridden"
    GradeStatusFinal       GradeStatus = "final"
)
