package handler

import (
	"context"
	"net/http"
	"strings"

	"posting-app/domain"
	"posting-app/infrastructure"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func AuthMiddleware(jwtService *infrastructure.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Bearer token required", http.StatusUnauthorized)
				return
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			user := &domain.User{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  domain.UserRole(claims.Role),
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(*domain.User)
		if !ok {
			http.Error(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		if user.Role != domain.UserRoleAdmin {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) *domain.User {
	user, _ := ctx.Value(UserContextKey).(*domain.User)
	return user
}
