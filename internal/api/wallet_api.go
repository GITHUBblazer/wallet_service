package api

import (
	"net/http"

	"wallet-service/internal/service"
)

type API struct {
	walletService service.WalletService
}

func NewAPI(walletService service.WalletService) *API {
	return &API{walletService: walletService}
}

func (a *API) Routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/deposit", a.DepositHandler)
	router.HandleFunc("/withdraw", a.WithdrawHandler)
	router.HandleFunc("/transfer", a.TransferHandler)
	router.HandleFunc("/balance", a.BalanceHandler)
	router.HandleFunc("/history", a.HistoryHandler)

	return router
}
