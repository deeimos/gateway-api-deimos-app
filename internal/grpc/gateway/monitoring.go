package gateway

import (
	"context"
	"gateway-api/internal/middleware"

	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"
	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Monitoring interface {
	StreamServerMetrics(ctx context.Context, serverID, userID string, stream servers_apiv1.ServersAPI_StreamServerMetricsServer) error
}

type streamAdapter struct {
	gateway_apiv1.PublicMonitoring_StreamServerMetricsServer
}

func (a *streamAdapter) Send(metric *servers_apiv1.ServerMetric) error {
	return a.PublicMonitoring_StreamServerMetricsServer.Send(&gateway_apiv1.PublicServerMetric{
		CpuUsage:      metric.CpuUsage,
		MemoryUsage:   metric.MemoryUsage,
		DiskUsage:     metric.DiskUsage,
		LoadAvg_1:     metric.LoadAvg_1,
		LoadAvg_5:     metric.LoadAvg_5,
		LoadAvg_15:    metric.LoadAvg_15,
		NetworkRx:     metric.NetworkRx,
		NetworkTx:     metric.NetworkTx,
		DiskRead:      metric.DiskRead,
		DiskWrite:     metric.DiskWrite,
		ProcessCount:  metric.ProcessCount,
		IoWait:        metric.IoWait,
		UptimeSeconds: metric.UptimeSeconds,
		Temperature:   metric.Temperature,
		Status:        metric.Status,
		Timestamp:     metric.Timestamp,
	})
}

func (s *serverApi) StreamServerMetrics(
	req *gateway_apiv1.PublicServerMetricsRequest,
	stream gateway_apiv1.PublicMonitoring_StreamServerMetricsServer,
) error {
	userID, ok := stream.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		return status.Error(codes.Internal, "userId is undefined")
	}

	err := s.monitoring.StreamServerMetrics(
		stream.Context(),
		req.GetServerId(),
		userID,
		&streamAdapter{stream},
	)
	if err != nil {
		return s.errMapper.HandleGRPC(err)
	}

	return nil
}
