package serverHandler

import (
	"context"
	"encoding/json"
	"gateway-api/internal/domain/models"
	"gateway-api/internal/lib/utils/formatTimestamp"
	"gateway-api/internal/lib/validation"
	"gateway-api/internal/middleware/auth"
	"gateway-api/internal/services/server"
	"net/http"
	"time"
)

type ServerHandler struct {
	service *server.Server
	timeout time.Duration
}

func NewServerHandler(service *server.Server, timeout time.Duration) *ServerHandler {
	return &ServerHandler{service: service, timeout: timeout}
}

func (h *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.EncryptedCreateServerModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, err, http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(auth.UserIDKey).(string)
	req.UserID = userID

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.CreateServer(ctx, &req)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *ServerHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req models.EncryptedServerModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, err, http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(auth.UserIDKey).(string)
	req.UserID = userID

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.UpdateServer(ctx, &req)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *ServerHandler) Get(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("id")
	if serverID == "" {
		http.Error(w, "missing id param", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(auth.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.Server(ctx, serverID, userID)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	dto := models.EncryptedServerModel{
		ID:                   resp.GetId(),
		EncryptedIp:          resp.GetEncryptedIp(),
		EncryptedPort:        resp.GetEncryptedPort(),
		EncryptedDisplayName: resp.GetEncryptedDisplayName(),
		IsMonitoringEnabled:  resp.GetIsMonitoringEnabled(),
		CreatedAt:            formatTimestamp.FormatTimestamp(resp.GetCreatedAt()),
	}

	_ = json.NewEncoder(w).Encode(dto)
}

func (h *ServerHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(auth.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.ServersList(ctx, userID)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	servers := make([]models.EncryptedServerModel, 0, len(resp.Servers))
	for _, server := range resp.Servers {
		servers = append(servers, models.EncryptedServerModel{
			ID:                   server.GetId(),
			EncryptedIp:          server.GetEncryptedIp(),
			EncryptedPort:        server.GetEncryptedPort(),
			EncryptedDisplayName: server.GetEncryptedDisplayName(),
			IsMonitoringEnabled:  server.GetIsMonitoringEnabled(),
			CreatedAt:            formatTimestamp.FormatTimestamp(server.GetCreatedAt()),
		})
	}
	_ = json.NewEncoder(w).Encode(servers)
}

func (h *ServerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("id")
	if serverID == "" {
		http.Error(w, "missing id param", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(auth.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.DeleteServer(ctx, serverID, userID)
	if err != nil {
		validation.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}
