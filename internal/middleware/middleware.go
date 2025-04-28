package middleware

import (
	"context"
	"strings"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
)

type Auth interface {
	GetUser(ctx context.Context, token string) (*auth_apiv1.GetUserResponse, error)
}
type contextKey string

const UserIDKey contextKey = "UserID"

func isWhitelisted(method string) bool {
	return strings.HasSuffix(method, "/Login") ||
		strings.HasSuffix(method, "/Register") ||
		strings.HasSuffix(method, "/Refresh")
}
