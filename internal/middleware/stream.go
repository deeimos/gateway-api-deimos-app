package middleware

import (
	"context"
	"fmt"
	"gateway-api/internal/lib/validation"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type wrappedServerStream struct {
	grpc.ServerStream
	WrappedContext context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

func WrapServerStream(stream grpc.ServerStream) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}
	return &wrappedServerStream{ServerStream: stream}
}

func NewAuthStreamInterceptor(auth Auth, errMapper *validation.ErrorMapper) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		fmt.Println("Stream method called:", info.FullMethod)

		if isWhitelisted(info.FullMethod) {
			return handler(srv, ss)
		}

		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errMapper.HandleGRPC(validation.ErrNoMetadata)
		}

		tokens := md.Get("token")
		if len(tokens) == 0 || strings.TrimSpace(tokens[0]) == "" {
			return errMapper.HandleGRPC(validation.ErrInvalidToken)
		}

		user, err := auth.GetUser(ss.Context(), tokens[0])
		if err != nil {
			return errMapper.HandleGRPC(validation.ErrInvalidToken)
		}

		wrapped := WrapServerStream(ss)
		wrapped.WrappedContext = context.WithValue(ss.Context(), UserIDKey, user.Id)

		return handler(srv, wrapped)
	}
}
