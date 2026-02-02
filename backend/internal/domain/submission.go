package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Submission struct {
	bun.BaseModel `bun:"table:submissions,alias:s"`

	ID               uuid.UUID        `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ExamID           uuid.UUID        `bun:"exam_id,notnull,type:uuid" json:"exam_id"`
	StudentID        string           `bun:"student_id,notnull" json:"student_id"`
	UploadedAt       time.Time        `bun:"uploaded_at,nullzero,notnull,default:current_timestamp" json:"uploaded_at"`
	ProcessingStatus ProcessingStatus `bun:"processing_status,notnull" json:"processing_status"`
	OCRResults       []OCRResult      `bun:"ocr_results,type:jsonb" json:"ocr_results"`
	Answers          []AnswerSegment  `bun:"answers,type:jsonb" json:"answers"`
	TenantID         uuid.UUID        `bun:"tenant_id,notnull,type:uuid" json:"tenant_id"`
}

type ProcessingStatus string

const (
	StatusPending    ProcessingStatus = "pending"
	StatusProcessing ProcessingStatus = "processing"
	StatusCompleted  ProcessingStatus = "completed"
	StatusFailed     ProcessingStatus = "failed"
)

type OCRResult struct {
	PageNumber    int           `json:"page_number"`
	RawText       string        `json:"raw_text"`
	Confidence    float64       `json:"confidence"`
	ImageURL      string        `json:"image_url"`
	BoundingBoxes []BoundingBox `json:"bounding_boxes"`
	CorrectedText *string       `json:"corrected_text"`
}

type BoundingBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AnswerSegment struct {
	ID           uuid.UUID     `json:"id"`
	SubmissionID uuid.UUID     `json:"submission_id"`
	QuestionID   uuid.UUID     `json:"question_id"`
	Text         string        `json:"text"`
	PageIndices  []int         `json:"page_indices"`
	BoundingBox  []BoundingBox `json:"bounding_box"`
	Diagrams     []string      `json:"diagrams"` // URLs to cropped diagram images
}
