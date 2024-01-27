package wallet

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"infotecs-test-task/internal/lib/logger"
	"infotecs-test-task/internal/models"
	"infotecs-test-task/internal/services"
	"log/slog"
	"net/http"
	"time"
)

type SendRequest struct {
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

// @Summary Transfer Money
// @Tags wallet | transaction
// @Description transfer money from one wallet to another
// @ID transfer-money
// @Accept json
// @Param walletId path string false "wallet id"
// @Param input body SendRequest true "receiving wallet info"
// @Success 200 "Money Transfered Successfully"
// @Failure 400 "Bad Request"
// @Failure 404 "Sending Wallet Not Found"
// @Failure 507 "Failed To Store New Transaction"
// @Router /wallet/{walletId}/send [post]
func Send(log *slog.Logger, service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "Send"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		walletId := chi.URLParam(r, "walletId")
		if walletId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req SendRequest

		err := render.DecodeJSON(r.Body, &req)
		log.Debug("request", slog.Any("req", req))
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		transaction := models.Transaction{
			To:     req.To,
			From:   walletId,
			Amount: req.Amount,
			Date:   time.Now(),
		}

		err = service.TransferMoney(transaction)
		if err != nil {
			if errors.Is(err, services.ErrSendingWalletNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if !errors.Is(err, services.ErrReceivingWalletNotFound) {
				log.Error("failed to store transaction", logger.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusInsufficientStorage)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
