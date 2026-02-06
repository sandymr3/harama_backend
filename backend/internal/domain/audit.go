package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AuditLog struct {
	bun.BaseModel `bun:"table:audit_log,alias:al"`

	ID         uuid.UUID              `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	EntityType string                 `bun:"entity_type,notnull" json:"entity_type"`
	EntityID   uuid.UUID              `bun:"entity_id,notnull,type:uuid" json:"entity_id"`
	EventType  string                 `bun:"event_type,notnull" json:"event_type"`
	ActorID    *uuid.UUID             `bun:"actor_id,type:uuid" json:"actor_id,omitempty"`
	ActorType  string                 `bun:"actor_type" json:"actor_type,omitempty"`
	Changes    map[string]interface{} `bun:"changes,type:jsonb" json:"changes"`
	Metadata   map[string]interface{} `bun:"metadata,type:jsonb" json:"metadata,omitempty"`
	Hash       string                 `bun:"hash,notnull" json:"hash"`
	CreatedAt  time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}
