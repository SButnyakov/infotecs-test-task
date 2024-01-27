package wallet

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"infotecs-test-task/internal/storage"
	"log/slog"
	"net/http"
)

type GetWalletResponse struct {
	Id      string  `json:"id"`
	Balance float32 `json:"balance"`
}

// @Summary Get Wallet
// @Tags wallet
// @Description get wallet
// @ID get-wallet
// @Param walletId path string false "wallet id"
// @Success 200 {object} GetWalletResponse
// @Failure 400 "Bad Request"
// @Failure 404 "Wallet Not Found"
// @Failure 500 "Internal Server Error"
// @Router /wallet/{walletId} [get]
func GetWallet(log *slog.Logger, service Service) http.HandlerFunc {
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

		wallet, err := service.GetWallet(walletId)
		if err != nil {
			if errors.Is(err, storage.ErrWalletNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, GetWalletResponse{
			Id:      wallet.Id,
			Balance: wallet.Balance,
		})
	}
}
