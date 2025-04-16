package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
)

type Auth interface {
	GetUser(ctx context.Context, token string) (*auth_apiv1.GetUserResponse, error)
}
type contextKey string

const userIDKey contextKey = "UserID"

func NewAuthInterceptor(auth Auth) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		whitelisted := []string{
			"/deimos.public_api.v1.PublicAuth/Login",
			"/deimos.public_api.v1.PublicAuth/Register",
			"/deimos.public_api.v1.PublicAuth/Refresh",
		}

		for _, method := range whitelisted {
			if info.FullMethod == method {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Metadata отсутствует")
		}

		tokens := md.Get("token")
		if len(tokens) == 0 || strings.TrimSpace(tokens[0]) == "" {
			return nil, status.Error(codes.Unauthenticated, "Token отсутствует")
		}

		user, err := auth.GetUser(ctx, tokens[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Недействительный токен")
		}

		ctx = context.WithValue(ctx, userIDKey, user.Id)
		return handler(ctx, req)
	}
}
