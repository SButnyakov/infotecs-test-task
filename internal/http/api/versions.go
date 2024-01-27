package api

import (
	"github.com/go-chi/chi/v5"
	v1 "infotecs-test-task/internal/http/api/v1"
	"infotecs-test-task/internal/http/api/v1/wallet"
	"log/slog"
)

func Versions(log *slog.Logger, service wallet.Service) chi.Router {
	router := chi.NewRouter()
	router.Mount("/v1", v1.Router(log, service))
	return router
}
