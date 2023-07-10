package main

import (
	"currency-service/internal/config"
	"currency-service/internal/http-server/handlers"
	"currency-service/internal/storage/sqlite"

	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()

	log := loggerInit(cfg.Env)
	log.Info("Starting currency-service", slog.String("env", cfg.Env))
	log.Debug("Debug mode")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", err)
		os.Exit(1)
	}

	/*
	m := make(map[string]float64)
	m["EUR"] = 12.22
	m["JPY"] = 22.22
	m["UAH"] = 32.22

	_ = storage.SaveCurrency(&core.LatestRates{
		Motd: core.AddMessage{Msg: "abc", Url: "def"},
		Success: true,
		Base: "USD",
		Date: "2023-07-07",
		Rates: m,
	})

	r, _ := storage.GetCurrencyByDate("2023-07-07")
	*/

	//_ = storage

	router := chi.NewRouter()
	
	router.Use(middleware.Logger)//logging routing requests
	router.Use(middleware.Recoverer)//recovery routing panic

	router.Get("/latest", handlers.New(log, storage))

	log.Info("Starting server...", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

	log.Error("Server stopped")
}

func loggerInit(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
