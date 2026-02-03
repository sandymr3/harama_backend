package middleware

import (
	"context"
	"net/http"

	"harama/internal/auth"

	"github.com/google/uuid"
)

func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantIDStr := r.Header.Get("X-Tenant-ID")
		if tenantIDStr == "" {
			http.Error(w, "X-Tenant-ID header required", http.StatusUnauthorized)
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			http.Error(w, "invalid X-Tenant-ID", http.StatusUnauthorized)
			return
		}

		ctx := auth.WithTenantID(r.Context(), tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
