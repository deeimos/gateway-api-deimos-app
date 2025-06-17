package metricsHandler

import (
	"context"
	"encoding/json"
	"gateway-api/internal/domain/models"
	"gateway-api/internal/lib/utils/formatTimestamp"
	"gateway-api/internal/lib/validation"
	"gateway-api/internal/middleware/auth"
	"gateway-api/internal/services/metrics"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type MetricsHandler struct {
	service     *metrics.Metrics
	authService auth.Auth
	timeout     time.Duration
}

func NewMetricsHandler(service *metrics.Metrics, authService auth.Auth, timeout time.Duration) *MetricsHandler {
	return &MetricsHandler{service: service, authService: authService, timeout: timeout}
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

	token := r.URL.Query().Get("token")
	if strings.TrimSpace(token) == "" {
		http.Error(w, "unauthorized (no token)", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUser(r.Context(), token)
	if err != nil {
		log.Printf("[WebSocket] auth error: %v, token: %s", err, token)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID := user.Id

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer func() {
		log.Printf("[WebSocket] disconnected: user=%s server=%s", userID, serverID)
		conn.Close()
	}()

	log.Printf("[WebSocket] connected: user=%s server=%s", userID, serverID)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	stream, err := h.service.StreamServerMetrics(ctx, serverID, userID)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		log.Printf("[WebSocket] stream error: %v", err)
		return
	}

	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("[WebSocket] ping error: %v", err)
				return
			}
		}
	}()

	for {
		metric, err := stream.Recv()
		if err != nil {
			log.Printf("[WebSocket] stream.Recv error: %v", err)
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
			log.Printf("[WebSocket] marshal error: %v", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("[WebSocket] write error: %v", err)
			break
		}
	}
}
