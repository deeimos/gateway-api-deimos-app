package router

import (
	authHandler "gateway-api/internal/handlers/auth"
	metricsHandler "gateway-api/internal/handlers/metrics"
	serverHandler "gateway-api/internal/handlers/server"
	authMiddleware "gateway-api/internal/middleware/auth"
	loggerMiddleware "gateway-api/internal/middleware/logger"
	authService "gateway-api/internal/services/auth"

	forecastHandler "gateway-api/internal/handlers/forecast"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(hostedFront string, log *slog.Logger, authService *authService.Auth, authHandler *authHandler.AuthHandler, serverHandler *serverHandler.ServerHandler, metricsHandler *metricsHandler.MetricsHandler, forecastHandler *forecastHandler.ForecastHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", hostedFront},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(loggerMiddleware.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/refresh", authHandler.Refresh)
	})

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware(authService))

		r.Route("/user", func(r chi.Router) {
			r.Get("/me", authHandler.User)
		})

		r.Route("/server", func(r chi.Router) {
			r.Get("/list", serverHandler.List)
			r.Get("/", serverHandler.Get)
			r.Post("/create", serverHandler.Create)
			r.Put("/update", serverHandler.Update)
			r.Delete("/delete/{id}", serverHandler.Delete)
		})

		r.Route("/metrics", func(r chi.Router) {
			r.Get("/", metricsHandler.Stream)
		})

		r.Route("/forecast", func(r chi.Router) {
			r.Get("/", forecastHandler.Forecast)
			// r.Get("/stream", forecastHandler.Stream) Пока что без него
		})
	})
	return router
}
