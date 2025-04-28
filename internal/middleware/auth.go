package middleware

import (
	"context"
	"fmt"
	"gateway-api/internal/lib/validation"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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
