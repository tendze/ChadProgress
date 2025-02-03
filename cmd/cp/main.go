package main

import (
	authclient "ChadProgress/internal/auth_client/http"
	"ChadProgress/internal/config"
	"ChadProgress/internal/http_server/handlers/url/authorization"
	"ChadProgress/internal/lib/logger/handlers/slogpretty"
	http2 "ChadProgress/internal/middleware/auth"
	userauthservice "ChadProgress/internal/services/auth"
	"ChadProgress/storage/postgres"
	"fmt"
	"github.com/go-chi/chi/v5"
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
		log.Error("failed to init storage:", err)
	}

	authServiceClient := authclient.NewAuthClient(cfg.AuthClient.BaseURL, log, time.Second*10)
	userAuthService := userauthservice.NewUserService(storage, authServiceClient, log)
	userHandler := authorization.NewUserAuthHandler(userAuthService, log)

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

	// Open endpoints
	router.Route("/authorization", func(r chi.Router) {
		router.Post("/register", userHandler.Register)
		router.Post("/login", userHandler.Login)
	})

	// Protected endpoints
	router.Route("/user", func(r chi.Router) {
		r.Use(authMiddleware)
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
