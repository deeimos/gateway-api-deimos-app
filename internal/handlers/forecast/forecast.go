package forecastHandler

import (
	"context"
	"encoding/json"
	"gateway-api/internal/domain/models"
	"gateway-api/internal/lib/utils/formatTimestamp"
	"gateway-api/internal/lib/validation"
	"gateway-api/internal/middleware/auth"
	"gateway-api/internal/services/forecast"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
)

type ForecastHandler struct {
	service *forecast.Forecast
	timeout time.Duration
}

func NewForecastHandler(service *forecast.Forecast, timeout time.Duration) *ForecastHandler {
	return &ForecastHandler{
		service: service,
		timeout: timeout,
	}
}

func (h *ForecastHandler) Forecast(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		validation.WriteError(w, validation.ErrInvalidToken, http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		validation.WriteError(w, validation.ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.ServerForecast(ctx, serverID, userID)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var dto []models.ServerForecast
	for _, point := range resp.Forecasts {
		dto = append(dto, convertForecastPoint(point))
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"server_id": resp.ServerId,
		"forecasts": dto,
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *ForecastHandler) Stream(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		validation.WriteError(w, validation.ErrInvalidToken, http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		validation.WriteError(w, validation.ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	streamClient, err := h.service.StreamServerForecast(ctx, serverID, userID)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}

	for {
		point, err := streamClient.Recv()
		if err != nil {
			conn.WriteJSON(map[string]string{"error": "stream error"})
			break
		}
		dto := convertForecastPoint(point)

		if err := conn.WriteJSON(dto); err != nil {
			break
		}
	}
}

func convertForecastPoint(p *servers_apiv1.ServerForecastPoint) models.ServerForecast {
	return models.ServerForecast{
		Timestamp:               formatTimestamp.FormatTimestamp(p.GetTimestamp()),
		CpuLoad:                 p.GetCpuLoad(),
		MemoryLoad:              p.GetMemoryLoad(),
		DiskUsage:               p.GetDiskUsage(),
		LoadAvg1:                p.GetLoadAvg_1(),
		LoadAvg5:                p.GetLoadAvg_5(),
		NetworkRx:               p.GetNetworkRx(),
		NetworkTx:               p.GetNetworkTx(),
		AvailabilityProbability: p.GetAvailabilityProbability(),
		Status:                  p.GetStatus(),
	}
}
