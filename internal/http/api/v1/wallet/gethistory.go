package wallet

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"infotecs-test-task/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

type TransactionsHistory []TransactionResponse

type TransactionResponse struct {
	Time   string  `json:"time"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

// @Summary Get History
// @Tags transaction
// @Description get transactions history
// @ID get-transactions
// @Param walletId path string false "wallet id"
// @Success 200 {object} TransactionsHistory
// @Failure 400 "Bad Request"
// @Failure 404 "Wallet Not Found"
// @Failure 500 "Internal Server Error"
// @Router /wallet/{walletId}/history [get]
func GetHistory(log *slog.Logger, service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "GetHistory"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		walletId := chi.URLParam(r, "walletId")
		if walletId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := service.GetWallet(walletId)
		if err != nil {
			if errors.Is(err, storage.ErrWalletNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		history, err := service.GetTransactionsHistory(walletId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := make(TransactionsHistory, len(history))
		for i, v := range history {
			resp[i] = TransactionResponse{
				Time:   v.Date.Format(time.RFC3339),
				From:   v.From,
				To:     v.To,
				Amount: v.Amount,
			}
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp)
	}
}
