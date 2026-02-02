package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Exam struct {
	bun.BaseModel `bun:"table:exams,alias:e"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Title     string     `bun:"title,notnull" json:"title"`
	Subject   string     `bun:"subject" json:"subject"`
	Questions []Question `bun:"rel:has-many,join:id=exam_id" json:"questions"`
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	TenantID  uuid.UUID  `bun:"tenant_id,notnull,type:uuid" json:"tenant_id"`
}
