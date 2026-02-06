package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"harama/internal/domain"

	"github.com/uptrace/bun"
)

type AuditRepo struct {
	db *bun.DB
}

func NewAuditRepo(db *bun.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) GetLastHash(ctx context.Context) (string, error) {
	var lastLog domain.AuditLog
	err := r.db.NewSelect().
		Model(&lastLog).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		// If no logs exist, return a starting seed
		return "initial_seed", nil
	}
	return lastLog.Hash, nil
}

func (r *AuditRepo) Save(ctx context.Context, log *domain.AuditLog) error {
	lastHash, err := r.GetLastHash(ctx)
	if err != nil {
		return err
	}

	// Calculate hash: SHA-256(previousHash + entityType + entityID + eventType + actorID + changesJSON)
	changesJSON, _ := json.Marshal(log.Changes)
	actorID := ""
	if log.ActorID != nil {
		actorID = log.ActorID.String()
	}

	data := fmt.Sprintf("%s|%s|%s|%s|%s", log.EntityType, log.EntityID, log.EventType, actorID, string(changesJSON))
	hasher := sha256.New()
	hasher.Write([]byte(lastHash + data))
	log.Hash = hex.EncodeToString(hasher.Sum(nil))

	_, err = r.db.NewInsert().Model(log).Exec(ctx)
	return err
}

func (r *AuditRepo) GetByEntity(ctx context.Context, entityType string, entityID string) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.NewSelect().
		Model(&logs).
		Where("entity_type = ?", entityType).
		Where("entity_id = ?", entityID).
		Order("created_at DESC").
		Scan(ctx)
	return logs, err
}
