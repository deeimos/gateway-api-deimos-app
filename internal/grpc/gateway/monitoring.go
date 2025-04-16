package gateway

import (
	"context"

	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"
)

type Monitoring interface {
	StreamServerMetrics(ctx context.Context, serverID, userID string, stream gateway_apiv1.PublicMonitoring_StreamServerMetricsServer) error
}

func (s *serverApi) StreamServerMetrics(req *gateway_apiv1.PublicServerMetricsRequest, stream gateway_apiv1.PublicMonitoring_StreamServerMetricsServer) error {
	// if req.GetServerId() == "" {
	// 	return status.Error(codes.InvalidArgument, "Отсутствует ID сервера")
	// }

	// if req.GetToken() == "" {
	// 	return status.Error(codes.Unauthenticated, "Отсутствует токен")
	// }

	// userID, err := s.auth.ValidateAccessToken(req.GetToken())
	// if err != nil {
	// 	return status.Error(codes.Unauthenticated, "Неверный или недействительный токен")
	// }

	// return s.servers.StreamServerMetrics(stream.Context(), req.GetServerId(), userID, stream)
	panic("implement me")
}
