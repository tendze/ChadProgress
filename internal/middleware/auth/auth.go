package middleware

import (
	"ChadProgress/internal/models"
	"context"
	"net/http"
	"strings"
)

// TokenValidator interface for auth service structures
type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}

// AuthMiddleware проверяет JWT токен и добавляет email пользователя в контекст запроса.
func AuthMiddleware(tokenValidator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractTokenFromHeader(r)
			if token == "" {
				http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
				return
			}

			userEmail, err := tokenValidator.ValidateToken(context.Background(), token)
			if err != nil || userEmail == "" {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), models.ContextUserKey, userEmail)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractTokenFromHeader extracts token from Authorization header.
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
