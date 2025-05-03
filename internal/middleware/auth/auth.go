package auth

import (
	"context"
	"gateway-api/internal/lib/validation"
	"net/http"
	"strings"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
)

type Auth interface {
	GetUser(ctx context.Context, token string) (*auth_apiv1.GetUserResponse, error)
}

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(auth Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if strings.TrimSpace(token) == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := auth.GetUser(r.Context(), token)
			if err != nil {
				validation.WriteError(w, err, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
