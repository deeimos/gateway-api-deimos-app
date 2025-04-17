package gateway

import (
	"context"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"
	"google.golang.org/grpc/metadata"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (*auth_apiv1.LoginResponse, error)
	Register(ctx context.Context, name string, email string, password string) (*auth_apiv1.RegisterResponse, error)
	Refresh(ctx context.Context, refersh string) (*auth_apiv1.RefreshResponse, error)
	GetUser(ctx context.Context, token string) (*auth_apiv1.GetUserResponse, error)
}

func (s *serverApi) Login(ctx context.Context, req *gateway_apiv1.LoginRequest) (*gateway_apiv1.LoginResponse, error) {
	user, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.LoginResponse{
		Id:           user.Id,
		Email:        user.Email,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt,
		Token:        user.Token,
		RefreshToken: user.RefreshToken,
	}, nil
}

func (s *serverApi) Register(ctx context.Context, req *gateway_apiv1.RegisterRequest) (*gateway_apiv1.RegisterResponse, error) {
	user, err := s.auth.Register(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.RegisterResponse{
		Id:           user.Id,
		Email:        user.Email,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt,
		Token:        user.Token,
		RefreshToken: user.RefreshToken,
	}, nil
}

func (s *serverApi) Refresh(ctx context.Context, req *gateway_apiv1.RefreshRequest) (*gateway_apiv1.RefreshResponse, error) {
	refreshToken, err := s.auth.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.RefreshResponse{
		Token:        refreshToken.Token,
		RefreshToken: refreshToken.RefreshToken,
	}, nil
}

func (s *serverApi) GetUser(ctx context.Context, req *gateway_apiv1.GetUserRequest) (*gateway_apiv1.GetUserResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("token")

	user, err := s.auth.GetUser(ctx, tokens[0])
	if err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	return &gateway_apiv1.GetUserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}
