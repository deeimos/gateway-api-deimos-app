package client

import (
	"crypto/tls"
	"fmt"
	"gateway-api/internal/config"

	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ServersClient struct {
	Conn   *grpc.ClientConn
	Client servers_apiv1.ServersAPIClient
}

func NewServersClient(cfg config.API) (*ServersClient, error) {
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

	return &ServersClient{
		Conn:   clientConn,
		Client: servers_apiv1.NewServersAPIClient(clientConn),
	}, nil
}
