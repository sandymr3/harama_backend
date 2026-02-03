package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextKey string

const (
	TenantKey contextKey = "tenant_id"
	UserKey   contextKey = "user_id"
)

func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, TenantKey, tenantID)
}

func GetTenantID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(TenantKey)
	if val == nil {
		return uuid.Nil, errors.New("tenant_id not found in context")
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid tenant_id type in context")
	}
	return id, nil
}
