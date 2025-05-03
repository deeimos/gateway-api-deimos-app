package metricsHandler

import (
	"context"
	"encoding/json"
	"gateway-api/internal/domain/models"
	"gateway-api/internal/lib/utils/formatTimestamp"
	"gateway-api/internal/lib/validation"
	"gateway-api/internal/middleware/auth"
	"gateway-api/internal/services/metrics"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type MetricsHandler struct {
	service *metrics.Metrics
	timeout time.Duration
}

func NewMetricsHandler(service *metrics.Metrics, timeout time.Duration) *MetricsHandler {
	return &MetricsHandler{service: service, timeout: timeout}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *MetricsHandler) Stream(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		http.Error(w, "missing server_id", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	stream, err := h.service.StreamServerMetrics(ctx, serverID, userID)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	for {
		metric, err := stream.Recv()
		if err != nil {
			break
		}
		dto := models.ServerMetric{
			CPUUsage:      metric.GetCpuUsage(),
			MemoryUsage:   metric.GetMemoryUsage(),
			DiskUsage:     metric.GetDiskUsage(),
			LoadAvg1:      metric.GetLoadAvg_1(),
			LoadAvg5:      metric.GetLoadAvg_5(),
			LoadAvg15:     metric.GetLoadAvg_15(),
			NetworkRx:     metric.GetNetworkRx(),
			NetworkTx:     metric.GetNetworkTx(),
			DiskRead:      metric.GetDiskRead(),
			DiskWrite:     metric.GetDiskWrite(),
			ProcessCount:  metric.GetProcessCount(),
			IOWait:        metric.GetIoWait(),
			UptimeSeconds: metric.GetUptimeSeconds(),
			Temperature:   metric.GetTemperature(),
			Status:        metric.GetStatus(),
			Timestamp:     formatTimestamp.FormatTimestamp(metric.GetTimestamp()),
		}

		data, err := json.Marshal(dto)
		if err != nil {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}
