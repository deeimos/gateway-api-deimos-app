package forecast

import (
	"context"
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
