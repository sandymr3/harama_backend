package service

import (
	"context"
	"harama/internal/domain"
	"harama/internal/repository/postgres"
)

type AuditService struct {
	repo *postgres.AuditRepo
}

func NewAuditService(repo *postgres.AuditRepo) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) GetLogsForEntity(ctx context.Context, entityType string, entityID string) ([]domain.AuditLog, error) {
	return s.repo.GetByEntity(ctx, entityType, entityID)
}
