package main

import (
	"log/slog"
	"net/http"
	"os"

	"cordell/internal/config"
	"cordell/internal/web"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	server := web.NewServer(logger)

	logger.Info("starting Cordell HTTP server", "address", cfg.HTTPAddress)

	if err := http.ListenAndServe(cfg.HTTPAddress, server.Routes()); err != nil {
		logger.Error("Cordell HTTP server stopped with error", "error", err)
		os.Exit(1)
	}
}
