package app

import (
	grpcapp "gateway-api/internal/app/grpc"
	"gateway-api/internal/client"
	"gateway-api/internal/config"
	"gateway-api/internal/lib/validation"
	"gateway-api/internal/services/auth"
	"gateway-api/internal/services/forecast"
	"gateway-api/internal/services/metrics"
	"gateway-api/internal/services/server"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, config config.Config) *App {
	authClient, err := client.NewAuthClient(config.APIs.AuthAPI)
	if err != nil {
		log.Error("failed to create auth client", slog.Any("err", err))
		panic(err)
	}

	serversClient, err := client.NewServersClient(config.APIs.ServersAPI)
	if err != nil {
		log.Error("failed to create servers client", slog.Any("err", err))
		panic(err)
	}

	authService := auth.New(log, authClient)
	serversService := server.New(log, serversClient)
	monitoringService := metrics.New(log, serversClient)
	forecastService := forecast.New(log, serversClient)

	errMapper := validation.NewErrorMapper()
	grpcApp := grpcapp.New(log, config.GRPCConfig.Port, authService, serversService, monitoringService, forecastService, errMapper)
	return &App{GRPCServer: grpcApp}
}
