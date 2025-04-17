package gateway

import (
	"context"
	"gateway-api/internal/domain/models"
	"gateway-api/internal/middleware"

	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"
	servers_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/servers-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "userId is undefined")
	}
	server, err := s.servers.CreateServer(ctx, &models.EncryptedCreateServerModel{
		UserID:               userID,
		EncryptedIP:          req.GetEncryptedIp(),
		EncryptedPort:        req.GetEncryptedPort(),
		EncryptedDisplayName: req.GetEncryptedDisplayName(),
		IsMonitoringEnabled:  req.IsMonitoringEnabled,
	})
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.CreateServerResponse{
		Id: server.Id,
	}, nil
}
func (s *serverApi) UpdateServer(ctx context.Context, req *gateway_apiv1.UpdateServerRequest) (*gateway_apiv1.UpdateServerResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "userId is undefined")
	}
	server, err := s.servers.UpdateServer(ctx, &models.EncryptedServerModel{
		ID:                   req.GetId(),
		UserID:               userID,
		EncryptedIP:          req.GetEncryptedIp(),
		EncryptedPort:        req.GetEncryptedPort(),
		EncryptedDisplayName: req.GetEncryptedDisplayName(),
		IsMonitoringEnabled:  req.IsMonitoringEnabled,
	})
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.UpdateServerResponse{
		Id: server.Id,
	}, nil
}
func (s *serverApi) GetServersList(ctx context.Context, req *gateway_apiv1.GetServersListRequest) (*gateway_apiv1.GetServersListResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "userId is undefined")
	}
	clientType, err := clientType(ctx)
	if err != nil {
		return nil, err
	}
	list, err := s.servers.ServersList(ctx, userID, clientType)
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}

	var result []*gateway_apiv1.GetServerResponse
	for _, server := range list.Servers {
		result = append(result, &gateway_apiv1.GetServerResponse{
			Id:                   server.Id,
			EncryptedIp:          server.EncryptedIp,
			EncryptedPort:        server.EncryptedPort,
			EncryptedDisplayName: server.EncryptedDisplayName,
			IsMonitoringEnabled:  server.IsMonitoringEnabled,
		})
	}
	return &gateway_apiv1.GetServersListResponse{
		Servers: result,
	}, nil
}
func (s *serverApi) GetServer(ctx context.Context, req *gateway_apiv1.GetServerRequest) (*gateway_apiv1.GetServerResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "userId is undefined")
	}

	clientType, err := clientType(ctx)
	if err != nil {
		return nil, err
	}

	server, err := s.servers.Server(ctx, req.GetId(), userID, clientType)
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.GetServerResponse{Id: server.Id, EncryptedIp: server.EncryptedIp, EncryptedPort: server.EncryptedPort, EncryptedDisplayName: server.EncryptedDisplayName}, nil
}

func (s *serverApi) DeleteServer(ctx context.Context, req *gateway_apiv1.DeleteServerRequest) (*gateway_apiv1.DeleteServerResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "userId is undefined")
	}
	server, err := s.servers.DeleteServer(ctx, req.GetId(), userID)
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.DeleteServerResponse{
		Id: server.Id,
	}, nil
}

func clientType(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Metadata отсутствует")
	}

	clientTypes := md.Get("client-type")
	if len(clientTypes) == 0 {
		return "", status.Error(codes.Unauthenticated, "client-type не указан")
	}

	return clientTypes[0], nil
}
