package auth

import (
	"context"
	"fmt"
	"gateway-api/internal/client"
	"log/slog"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
)

type Auth struct {
	log        *slog.Logger
	authClient *client.AuthClient
}

func New(log *slog.Logger, authClient *client.AuthClient) *Auth {
	return &Auth{
		log:        log,
		authClient: authClient,
	}
}

func (auth *Auth) Register(ctx context.Context, name string, email string, password string) (*auth_apiv1.RegisterResponse, error) {
	const op = "auth.RegisterUser"

	log := auth.log.With(slog.String("op", op))

	log.Info("GRPC")
	resp, err := auth.authClient.Client.Register(ctx, &auth_apiv1.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (auth *Auth) Login(ctx context.Context, email string, password string) (*auth_apiv1.LoginResponse, error) {
	const op = "auth.Login"

	log := auth.log.With(slog.String("op", op))

	log.Info("GRPC")
	resp, err := auth.authClient.Client.Login(ctx, &auth_apiv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (auth *Auth) Refresh(ctx context.Context, refreshToken string) (*auth_apiv1.RefreshResponse, error) {
	const op = "auth.Refresh"
	log := auth.log.With(slog.String("op", op))

	log.Info("GRPC")
	resp, err := auth.authClient.Client.Refresh(ctx, &auth_apiv1.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (auth *Auth) GetUser(ctx context.Context, token string) (*auth_apiv1.GetUserResponse, error) {
	const op = "auth.GetUser"

	log := auth.log.With(slog.String("op", op))

	log.Info("GRPC")
	resp, err := auth.authClient.Client.GetUser(ctx, &auth_apiv1.GetUserRequest{
		Token: token,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}
