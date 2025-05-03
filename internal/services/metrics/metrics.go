package metrics

import (
	"context"
	"fmt"
	"gateway-api/internal/client"
	"log/slog"

	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
)

type Metrics struct {
	log           *slog.Logger
	serversClient *client.ServersClient
}

func New(log *slog.Logger, serversClient *client.ServersClient) *Metrics {
	return &Metrics{
		log:           log,
		serversClient: serversClient,
	}
}

func (metrics *Metrics) StreamServerMetrics(ctx context.Context, serverID, userID string) (servers_apiv1.ServersAPI_StreamServerMetricsClient, error) {
	const op = "metrics.StreamClient"

	log := metrics.log.With(slog.String("op", op))
	log.Info("Starting metrics stream via gRPC")

	streamClient, err := metrics.serversClient.Client.StreamServerMetrics(ctx, &servers_apiv1.ServerMetricsRequest{
		ServerId: serverID,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return streamClient, nil
}
