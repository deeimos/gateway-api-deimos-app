package app

import (
	httpapp "gateway-api/internal/app/http"
	"gateway-api/internal/client"
	"gateway-api/internal/config"
	authHandler "gateway-api/internal/handlers/auth"
	forecastHandler "gateway-api/internal/handlers/forecast"
	metricsHandler "gateway-api/internal/handlers/metrics"
	serverHandler "gateway-api/internal/handlers/server"
	"gateway-api/internal/services/auth"
	"gateway-api/internal/services/forecast"
	"gateway-api/internal/services/metrics"
	"gateway-api/internal/services/server"
	"log/slog"
)

type App struct {
	HttpServer *httpapp.App
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
	serverService := server.New(log, serversClient)
	metrcisService := metrics.New(log, serversClient)
	forecastService := forecast.New(log, serversClient)

	authHandler := authHandler.NewAuthHandler(authService, config.Timeout)
	serverHandler := serverHandler.NewServerHandler(serverService, config.Timeout)
	metricsHandler := metricsHandler.NewMetricsHandler(metrcisService, config.Timeout)
	forecastHandler := forecastHandler.NewForecastHandler(forecastService, config.Timeout)

	httpApp := httpapp.New(log, config, authService, authHandler, serverHandler, metricsHandler, forecastHandler)
	return &App{HttpServer: httpApp}
}
