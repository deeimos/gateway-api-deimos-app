package gateway

import (
	"context"
	"gateway-api/internal/domain/models"

	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"
	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Servers interface {
	CreateServer(ctx context.Context, encryptedServer *models.EncryptedCreateServerModel) (*servers_apiv1.CreateServerResponse, error)
	UpdateServer(ctx context.Context, encryptedServer *models.EncryptedServerModel) (*servers_apiv1.UpdateServerResponse, error)
	ServersList(ctx context.Context, userID string, clientType string) (*servers_apiv1.GetServersListResponse, error)
	Server(ctx context.Context, serverID string, userID string, clientType string) (*servers_apiv1.GetServerResponse, error)
	DeleteServer(ctx context.Context, serverID, userID string) (*servers_apiv1.DeleteServerResponse, error)
}

func (s *serverApi) CreateServer(ctx context.Context, req *gateway_apiv1.CreateServerRequest) (*gateway_apiv1.CreateServerResponse, error) {
	// if req.GetUserId() == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует ID пользователя")
	// }

	// id, err := s.servers.CreateServer(ctx, &models.EncryptedCreateServerModel{
	// 	UserID:               req.GetUserId(),
	// 	EncryptedIP:          req.GetEncryptedIp(),
	// 	EncryptedPort:        req.GetEncryptedPort(),
	// 	EncryptedDisplayName: req.GetEncryptedDisplayName(),
	// 	IsMonitoringEnabled:  req.IsMonitoringEnabled,
	// })
	// if err != nil {
	// 	return nil, s.errMapper.HandleGRPC(err)
	// }
	// return &servers_apiv1.CreateServerResponse{Id: id}, nil
	panic("implement me")
}
func (s *serverApi) UpdateServer(ctx context.Context, req *gateway_apiv1.UpdateServerRequest) (*gateway_apiv1.UpdateServerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует ID сервера")
	}
	// if req.GetUserId() == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует ID пользователя")
	// }
	// id, err := s.servers.CreateServer(ctx, &models.EncryptedCreateServerModel{
	// 	UserID:               req.GetUserId(),
	// 	EncryptedIP:          req.GetEncryptedIp(),
	// 	EncryptedPort:        req.GetEncryptedPort(),
	// 	EncryptedDisplayName: req.GetEncryptedDisplayName(),
	// 	IsMonitoringEnabled:  req.IsMonitoringEnabled,
	// })
	// if err != nil {
	// 	return nil, s.errMapper.HandleGRPC(err)
	// }
	// return &servers_apiv1.UpdateServerResponse{Id: id}, nil
	panic("implement me")
}
func (s *serverApi) GetServersList(ctx context.Context, req *gateway_apiv1.GetServersListRequest) (*gateway_apiv1.GetServersListResponse, error) {
	// clientType, err := checkMetadata(ctx)
	// if err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует clientType")
	// }
	// if req.GetUserId() == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует ID пользователя")
	// }
	// serversList, err := s.servers.ServersList(ctx, req.GetUserId(), clientType)
	// if err != nil {
	// 	return nil, s.errMapper.HandleGRPC(err)
	// }
	// var response servers_apiv1.GetServersListResponse
	// for _, srv := range serversList {
	// 	response.Servers = append(response.Servers, &servers_apiv1.GetServerResponse{
	// 		Id:                   srv.ID,
	// 		EncryptedIp:          srv.EncryptedIP,
	// 		EncryptedPort:        srv.EncryptedPort,
	// 		EncryptedDisplayName: srv.EncryptedDisplayName,
	// 		IsMonitoringEnabled:  srv.IsMonitoringEnabled,
	// 		CreatedAt:            timestamppb.New(srv.CreatedAt),
	// 	})
	// }

	// return &response, nil
	panic("implement me")
}
func (s *serverApi) GetServer(ctx context.Context, req *gateway_apiv1.GetServerRequest) (*gateway_apiv1.GetServerResponse, error) {
	// clientType, err := checkMetadata(ctx)
	// if err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует clientType")
	// }
	// if req.GetId() == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует ID сервера")
	// }
	// if req.GetUserId() == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует ID пользователя")
	// }
	// server, err := s.servers.Server(ctx, req.GetId(), req.GetUserId(), clientType)
	// if err != nil {
	// 	return nil, s.errMapper.HandleGRPC(err)
	// }
	// return &servers_apiv1.GetServerResponse{Id: server.ID, EncryptedIp: server.EncryptedIP, EncryptedPort: server.EncryptedPort, EncryptedDisplayName: server.EncryptedDisplayName}, nil

	panic("implement me")
}

func (s *serverApi) DeleteServer(ctx context.Context, req *gateway_apiv1.DeleteServerRequest) (*gateway_apiv1.DeleteServerResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует ID сервера")
	}
	// if req.GetUserId() == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "Отсутствует ID пользователя")
	// }
	// err := s.servers.DeleteServer(ctx, req.GetId(), req.GetUserId())
	// if err != nil {
	// 	return nil, s.errMapper.HandleGRPC(err)
	// }

	// return &servers_apiv1.DeleteServerResponse{
	// 	Id: req.GetId(),
	// }, nil
	panic("implement me")
}
