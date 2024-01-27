package v1

import (
	"github.com/go-chi/chi/v5"
	"infotecs-test-task/internal/http/api/v1/wallet"
	"log/slog"
)

// @title Infotecs EWallet API
// @version 1.0
// @description API HTTPServer for Infotecs EWallet Test Task

// @host localhost:8080
// @BasePath /api/v1

func Router(log *slog.Logger, service wallet.Service) chi.Router {
	router := chi.NewRouter()
	router.Mount("/wallet", wallet.Router(log, service))
	return router
}
