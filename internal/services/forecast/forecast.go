package forecast

import (
	"context"
	"fmt"
	"gateway-api/internal/client"
	"gateway-api/internal/lib/validation"
	"log/slog"

	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
)

type Forecast struct {
	log           *slog.Logger
	serversClient *client.ServersClient
}

func New(log *slog.Logger, serversClient *client.ServersClient) *Forecast {
	return &Forecast{
		log:           log,
		serversClient: serversClient,
	}
}
func (forecast *Forecast) ServerForecast(ctx context.Context, serverID, userID string) (*servers_apiv1.ServerForecastResponse, error) {
	const op = "forecast.ServerForecast"

	log := forecast.log.With(slog.String("op", op))
	log.Info("GRPC")

	resp, err := forecast.serversClient.Client.ServerForecast(ctx, &servers_apiv1.ServerForecastRequest{
		UserId:   userID,
		ServerId: serverID,
	})
	if err != nil {
		return nil, validation.HandleGRPCServiceError(log, op, err)
	}

	return resp, nil
}

func (forecast *Forecast) StreamServerForecast(ctx context.Context, serverID, userID string) (servers_apiv1.ServersAPI_StreamServerForecastClient, error) {
	const op = "forecast.StreamClient"

	log := forecast.log.With(slog.String("op", op))
	log.Info("Starting forecast stream via gRPC")

	streamClient, err := forecast.serversClient.Client.StreamServerForecast(ctx, &servers_apiv1.ServerForecastStreamRequest{
		ServerId: serverID,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return streamClient, nil
}
