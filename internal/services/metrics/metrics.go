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

func (metrics *Metrics) StreamServerMetrics(ctx context.Context, serverID, userID string, stream servers_apiv1.ServersAPI_StreamServerMetricsServer) error {
	const op = "server.StreamServerMetrics"

	log := metrics.log.With(slog.String("op", op))
	log.Info("GRPC")
	streamClient, err := metrics.serversClient.Client.StreamServerMetrics(ctx, &servers_apiv1.ServerMetricsRequest{
		ServerId: serverID,
		UserId:   userID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for {
		metric, err := streamClient.Recv()
		if err != nil {
			return fmt.Errorf("%s: failed to receive from stream: %w", op, err)
		}

		if err := stream.Send(metric); err != nil {
			return fmt.Errorf("%s: failed to forward metric: %w", op, err)
		}
	}
}
