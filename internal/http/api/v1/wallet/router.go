package wallet

import (
	"github.com/go-chi/chi/v5"
	"infotecs-test-task/internal/models"
	"log/slog"
)

const packagePath = "http.api.v1.wallet."

type Service interface {
	CreateWallet() (models.Wallet, error)
	GetWallet(id string) (models.Wallet, error)
	TransferMoney(transaction models.Transaction) error
	GetTransactionsHistory(id string) ([]models.Transaction, error)
}

func Router(log *slog.Logger, service Service) chi.Router {
	router := chi.NewRouter()
	router.Post("/", Create(log, service))
	router.Post("/{walletId}/send", Send(log, service))
	router.Get("/{walletId}/history", GetHistory(log, service))
	router.Get("/{walletId}", GetWallet(log, service))
	return router
}
