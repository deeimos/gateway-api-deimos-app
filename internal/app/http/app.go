package httpapp

import (
	"context"
	"fmt"
	"gateway-api/internal/config"
	authHandler "gateway-api/internal/handlers/auth"
	forecastHandler "gateway-api/internal/handlers/forecast"
	metricsHandler "gateway-api/internal/handlers/metrics"
	serverHandler "gateway-api/internal/handlers/server"
	"gateway-api/internal/router"
	authService "gateway-api/internal/services/auth"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	log    *slog.Logger
	port   int
	server *http.Server
}

func New(
	log *slog.Logger,
	config config.Config,
	authService *authService.Auth,
	authHandler *authHandler.AuthHandler,
	serverHandler *serverHandler.ServerHandler,
	metricsHandler *metricsHandler.MetricsHandler,
	forecastHandler *forecastHandler.ForecastHandler,
) *App {
	appRouter := router.NewRouter(config.HostedFront, log, authService, authHandler, serverHandler, metricsHandler, forecastHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      appRouter,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return &App{
		log:    log,
		port:   config.Port,
		server: srv,
	}
}

func (a *App) Run() {
	a.log.Info("REST API is running", slog.Int("port", a.port))

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("failed to start HTTP server", slog.String("err", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.log.Info("Shutting down HTTP server...")
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("HTTP shutdown error", slog.String("err", err.Error()))
	}
}
