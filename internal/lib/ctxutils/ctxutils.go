package ctxutils

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ParseClientMetadata(ctx context.Context) (clientType string, token string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.Unauthenticated, "Metadata отсутствует")
		return
	}

	clientTypes := md.Get("client-type")
	if len(clientTypes) == 0 {
		err = status.Error(codes.InvalidArgument, "client-type не указан")
		return
	}
	clientType = clientTypes[0]

	tokens := md.Get("token")
	if len(tokens) > 0 {
		token = tokens[0]
	}

	return
}
