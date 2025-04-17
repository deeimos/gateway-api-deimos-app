package server

import (
	"context"
	"gateway-api/internal/client"
	"gateway-api/internal/domain/models"
	"gateway-api/internal/lib/validation"
	"log/slog"

	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	log           *slog.Logger
	serversClient *client.ServersClient
}

func New(log *slog.Logger, serversClient *client.ServersClient) *Server {
	return &Server{
		log:           log,
		serversClient: serversClient,
	}
}

func (s *Server) CreateServer(ctx context.Context, serverModel *models.EncryptedCreateServerModel) (*servers_apiv1.CreateServerResponse, error) {
	const op = "server.CreateServer"

	log := s.log.With(slog.String("op", op))
	log.Info("GRPC")

	resp, err := s.serversClient.Client.CreateServer(ctx, &servers_apiv1.CreateServerRequest{
		UserId:               serverModel.UserID,
		EncryptedIp:          serverModel.EncryptedIP,
		EncryptedPort:        serverModel.EncryptedPort,
		EncryptedDisplayName: serverModel.EncryptedDisplayName,
		IsMonitoringEnabled:  serverModel.IsMonitoringEnabled,
	})
	if err != nil {
		return nil, validation.HandleGRPCServiceError(log, op, err)
	}

	return resp, nil
}

func (s *Server) UpdateServer(ctx context.Context, serverModel *models.EncryptedServerModel) (*servers_apiv1.UpdateServerResponse, error) {
	const op = "server.UpdateServer"

	log := s.log.With(slog.String("op", op))
	log.Info("GRPC")

	resp, err := s.serversClient.Client.UpdateServer(ctx, &servers_apiv1.UpdateServerRequest{
		Id:                   serverModel.ID,
		UserId:               serverModel.UserID,
		EncryptedIp:          serverModel.EncryptedIP,
		EncryptedPort:        serverModel.EncryptedPort,
		EncryptedDisplayName: serverModel.EncryptedDisplayName,
		IsMonitoringEnabled:  serverModel.IsMonitoringEnabled,
	})
	if err != nil {
		return nil, validation.HandleGRPCServiceError(log, op, err)
	}

	return resp, nil
}

func (s *Server) Server(ctx context.Context, serverID string, userID string, clientType string) (*servers_apiv1.GetServerResponse, error) {
	const op = "server.UpdateServer"

	log := s.log.With(slog.String("op", op))
	log.Info("GRPC")

	ctx = metadata.AppendToOutgoingContext(ctx, "client-type", clientType)
	resp, err := s.serversClient.Client.GetServer(ctx, &servers_apiv1.GetServerRequest{
		Id:     serverID,
		UserId: userID,
	})
	if err != nil {
		return nil, validation.HandleGRPCServiceError(log, op, err)
	}

	return resp, nil
}

func (s *Server) ServersList(ctx context.Context, userID string, clientType string) (*servers_apiv1.GetServersListResponse, error) {
	const op = "server.UpdateServer"

	log := s.log.With(slog.String("op", op))
	log.Info("GRPC")

	ctx = metadata.AppendToOutgoingContext(ctx, "client-type", clientType)
	resp, err := s.serversClient.Client.GetServersList(ctx, &servers_apiv1.GetServersListRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, validation.HandleGRPCServiceError(log, op, err)
	}

	return resp, nil
}

func (s *Server) DeleteServer(ctx context.Context, serverID, userID string) (*servers_apiv1.DeleteServerResponse, error) {
	const op = "server.DeleteServer"

	log := s.log.With(slog.String("op", op))
	log.Info("GRPC")

	resp, err := s.serversClient.Client.DeleteServer(ctx, &servers_apiv1.DeleteServerRequest{
		UserId: userID,
		Id:     serverID,
	})
	if err != nil {
		return nil, validation.HandleGRPCServiceError(log, op, err)
	}

	return resp, nil
}
