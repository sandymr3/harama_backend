package domain

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Rubric struct {
	bun.BaseModel `bun:"table:rubrics,alias:r"`

	ID                 uuid.UUID           `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	QuestionID         uuid.UUID           `bun:"question_id,notnull,type:uuid" json:"question_id"`
	FullCreditCriteria []Criterion         `bun:"full_credit_criteria,notnull,type:jsonb" json:"full_credit_criteria"`
	PartialCreditRules []PartialCreditRule `bun:"partial_credit_rules,notnull,type:jsonb" json:"partial_credit_rules"`
	CommonMistakes     []CommonMistake     `bun:"common_mistakes,notnull,type:jsonb" json:"common_mistakes"`
	KeyConcepts        []string            `bun:"key_concepts,type:jsonb" json:"key_concepts"`
	GradingNotes       string              `bun:"grading_notes" json:"grading_notes"`
	StrictMode         bool                `bun:"strict_mode,default:false" json:"strict_mode"`
}

type Criterion struct {
    ID          string
    Description string
    Points      float64
    Required    bool
    Category    string
}

type PartialCreditRule struct {
    ID          string
    Condition   string
    Points      float64
    Description string
    Dependencies []string
}

type CommonMistake struct {
    ID          string
    Description string
    Penalty     float64
    Category    string
    Frequency   int
}
