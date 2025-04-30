package gateway

import (
	"context"
	"gateway-api/internal/middleware"

	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"
	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Forecast interface {
	StreamServerForecast(ctx context.Context, serverID, userID string, stream servers_apiv1.ServersAPI_StreamServerForecastServer) error
	ServerForecast(ctx context.Context, serverID, userID string) (*servers_apiv1.ServerForecastResponse, error)
}

type forecastAdapter struct {
	gateway_apiv1.PublicForecast_StreamServerForecastServer
}

func (a *forecastAdapter) Send(point *servers_apiv1.ServerForecastPoint) error {
	return a.PublicForecast_StreamServerForecastServer.Send(&gateway_apiv1.PublicForecastPoint{
		CpuLoad:    point.CpuLoad,
		MemoryLoad: point.MemoryLoad,
		DiskUsage:  point.DiskUsage,
		LoadAvg_1:  point.LoadAvg_1,
		LoadAvg_5:  point.LoadAvg_5,
		NetworkRx:  point.NetworkRx,
		NetworkTx:  point.NetworkTx,
		Status:     point.Status,
		Timestamp:  point.Timestamp,
	})
}

func (s *serverApi) StreamServerForecast(
	req *gateway_apiv1.PublicForecastStreamRequest,
	stream gateway_apiv1.PublicForecast_StreamServerForecastServer,
) error {
	userID, ok := stream.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		return status.Error(codes.Internal, "userId is undefined")
	}

	err := s.forecast.StreamServerForecast(
		stream.Context(),
		req.GetServerId(),
		userID,
		&forecastAdapter{stream},
	)
	if err != nil {
		return s.errMapper.HandleGRPC(err)
	}

	return nil
}

func (s *serverApi) ServerForecast(ctx context.Context, req *gateway_apiv1.PublicForecastRequest) (*gateway_apiv1.PublicForecastResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "userId is undefined")
	}
	resp, err := s.forecast.ServerForecast(ctx, userID, req.GetServerId())
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}

	var points []*gateway_apiv1.PublicForecastPoint
	for _, point := range resp.Forecasts {
		points = append(points, &gateway_apiv1.PublicForecastPoint{
			CpuLoad:    point.CpuLoad,
			MemoryLoad: point.MemoryLoad,
			DiskUsage:  point.DiskUsage,
			LoadAvg_1:  point.LoadAvg_1,
			LoadAvg_5:  point.LoadAvg_5,
			NetworkRx:  point.NetworkRx,
			NetworkTx:  point.NetworkTx,
			Status:     point.Status,
			Timestamp:  point.Timestamp,
		})
	}
	return &gateway_apiv1.PublicForecastResponse{
		ServerId:  resp.ServerId,
		Forecasts: points,
	}, nil
}
