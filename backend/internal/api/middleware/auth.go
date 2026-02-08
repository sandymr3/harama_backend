package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"harama/internal/auth"

	"github.com/google/uuid"
)

// SupabaseAuthMiddleware validates Supabase JWT tokens from the Authorization header.
// It extracts the user's sub (user_id) claim and uses it as the tenant_id.
// Falls back to X-Tenant-ID header for development/backward compatibility.
func SupabaseAuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try Authorization: Bearer <token> first
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				claims, err := validateHS256JWT(token, jwtSecret)
				if err != nil {
					http.Error(w, "invalid or expired token", http.StatusUnauthorized)
					return
				}

				// Extract user ID from 'sub' claim
				sub, ok := claims["sub"].(string)
				if !ok || sub == "" {
					http.Error(w, "token missing sub claim", http.StatusUnauthorized)
					return
				}

				userID, err := uuid.Parse(sub)
				if err != nil {
					http.Error(w, "invalid user id in token", http.StatusUnauthorized)
					return
				}

				// Use user_id as tenant_id (single-tenant per user model)
				ctx := auth.WithTenantID(r.Context(), userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Fallback: X-Tenant-ID header (for dev/testing)
			tenantIDStr := r.Header.Get("X-Tenant-ID")
			if tenantIDStr != "" {
				tenantID, err := uuid.Parse(tenantIDStr)
				if err != nil {
					http.Error(w, "invalid X-Tenant-ID", http.StatusUnauthorized)
					return
				}
				ctx := auth.WithTenantID(r.Context(), tenantID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, "authorization required: provide Bearer token or X-Tenant-ID", http.StatusUnauthorized)
		})
	}
}

// validateHS256JWT performs a minimal HS256 JWT validation.
// For production, consider using a full JWT library with expiry checks.
func validateHS256JWT(tokenStr, secret string) (map[string]interface{}, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errInvalidToken
	}

	// Verify signature
	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, errInvalidToken
	}

	// Decode payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errInvalidToken
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errInvalidToken
	}

	return claims, nil
}

var errInvalidToken = &tokenError{msg: "invalid token"}

type tokenError struct {
	msg string
}

func (e *tokenError) Error() string {
	return e.msg
}
