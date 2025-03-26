package main

import (
	authclient "ChadProgress/internal/auth_client/http"
	"ChadProgress/internal/config"
	"ChadProgress/internal/http_server/handlers/url/authorization"
	userhandler "ChadProgress/internal/http_server/handlers/url/user"
	"ChadProgress/internal/lib/logger/handlers/slogpretty"
	http2 "ChadProgress/internal/middleware/auth"
	userauthservice "ChadProgress/internal/services/authorization"
	userservice "ChadProgress/internal/services/user"
	"ChadProgress/storage/postgres"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.Username,
		cfg.DB.DBPassword,
		cfg.DB.Host, cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	storage, err := postgres.New(dsn)
	if err != nil {
		log.Error("failed to init storage:", slog.String("errormsg", err.Error()))
		return
	}

	authServiceClient := authclient.NewAuthClient(cfg.AuthClient.BaseURL, log, time.Second*10)
	userAuthService := userauthservice.NewUserAuthService(storage, authServiceClient, log)
	userAuthHandler := authorization.NewUserAuthHandler(userAuthService, log)

	userService := userservice.NewUserService(storage, log)
	userHandler := userhandler.NewUserHandler(log, userService)

	router := chi.NewRouter()

	serverAddr := cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	authMiddleware := http2.AuthMiddleware(authServiceClient)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешаем все источники
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// Open endpoints
	router.Route("/authorization", func(r chi.Router) {
		r.Post("/register", userAuthHandler.Register)
		r.Post("/login", userAuthHandler.Login)
	})

	// Protected endpoints
	router.Route("/user", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/trainers/profile", userHandler.CreateTrainer)
		r.Get("/trainers/profile", userHandler.GetTrainerProfile)
		r.Get("/trainers/clients", userHandler.GetTrainersClients)
		r.Post("/training-plan", userHandler.CreatePlan)

		r.Post("/clients/profile", userHandler.CreateClient)
		r.Patch("/clients/select-trainers", userHandler.SelectTrainer)
		r.Get("/clients/profile", userHandler.GetClientProfile)
		r.Post("/clients/metrics", userHandler.AddMetrics)
		r.Get("/clients/metrics", userHandler.GetMetrics)
	})

	log.Info("server started", slog.String("servaddr", serverAddr))
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
