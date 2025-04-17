package middleware

import (
	"context"
	"fmt"
	"gateway-api/internal/lib/validation"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
)

type Auth interface {
	GetUser(ctx context.Context, token string) (*auth_apiv1.GetUserResponse, error)
}
type contextKey string

const UserIDKey contextKey = "UserID"

func NewAuthInterceptor(auth Auth, errMapper *validation.ErrorMapper) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		fmt.Println("Method called:", info.FullMethod)

		if isWhitelisted(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMapper.HandleGRPC(validation.ErrNoMetadata)
		}

		tokens := md.Get("token")
		if len(tokens) == 0 || strings.TrimSpace(tokens[0]) == "" {
			return nil, errMapper.HandleGRPC(validation.ErrInvalidToken)
		}

		user, err := auth.GetUser(ctx, tokens[0])
		if err != nil {
			return nil, errMapper.HandleGRPC(validation.ErrInvalidToken)
		}

		ctx = context.WithValue(ctx, UserIDKey, user.Id)
		return handler(ctx, req)
	}
}

func isWhitelisted(method string) bool {
	return strings.HasSuffix(method, "/Login") ||
		strings.HasSuffix(method, "/Register") ||
		strings.HasSuffix(method, "/Refresh")
}
