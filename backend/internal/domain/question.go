package domain

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Question struct {
	bun.BaseModel `bun:"table:questions,alias:q"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ExamID         uuid.UUID  `bun:"exam_id,notnull,type:uuid" json:"exam_id"`
	QuestionText   string     `bun:"question_text,notnull" json:"question_text"`
	Points         int        `bun:"points,notnull" json:"points"`
	AnswerType     AnswerType `bun:"answer_type,notnull" json:"answer_type"`
	QuestionNumber string     `bun:"question_number" json:"question_number,omitempty"`
	QuestionGroup  string     `bun:"question_group" json:"question_group,omitempty"`
	Rubric         *Rubric    `bun:"rel:has-one,join:id=question_id" json:"rubric"`
	VisualAids     []string   `bun:"visual_aids,type:jsonb" json:"visual_aids"`
}

type AnswerType string

const (
	AnswerTypeShortAnswer AnswerType = "short_answer"
	AnswerTypeEssay       AnswerType = "essay"
	AnswerTypeMCQ         AnswerType = "mcq"
	AnswerTypeDiagram     AnswerType = "diagram"
)
