package client

import (
	"crypto/tls"
	"fmt"
	"gateway-api/internal/config"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Conn   *grpc.ClientConn
	Client auth_apiv1.AuthAPIClient
}

func NewAuthClient(cfg config.API) (*AuthClient, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var creds credentials.TransportCredentials
	if cfg.UseTLS {
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	} else {
		creds = insecure.NewCredentials()
	}

	clientConn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return &AuthClient{
		Conn:   clientConn,
		Client: auth_apiv1.NewAuthAPIClient(clientConn),
	}, nil
}
