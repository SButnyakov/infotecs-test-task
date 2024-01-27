package wallet

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type CreateWalletResponse struct {
	Id      string  `json:"id"`
	Balance float32 `json:"balance"`
}

// @Summary Create Wallet
// @Tags wallet
// @Description create a new wallet
// @ID create-wallet
// @Success 200 {object} CreateWalletResponse
// @Failure 507 "Failed to store new wallet"
// @Router /wallet [post]
func Create(log *slog.Logger, service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "Create"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		wallet, err := service.CreateWallet()
		if err != nil {
			w.WriteHeader(http.StatusInsufficientStorage)
			return
		}

		response := CreateWalletResponse{
			Id:      wallet.Id,
			Balance: wallet.Balance,
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response)
	}
}
