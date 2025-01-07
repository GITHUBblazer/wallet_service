package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func (a *API) DepositHandler(w http.ResponseWriter, r *http.Request) {
	userID, amount, err := parseRequestParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.walletService.Deposit(r.Context(), userID, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Deposit successful"))
}

func (a *API) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	userID, amount, err := parseRequestParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.walletService.Withdraw(r.Context(), userID, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Withdrawal successful"))
}

func (a *API) TransferHandler(w http.ResponseWriter, r *http.Request) {
	fromUserID, toUserID, amount, err := parseTransferRequestParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.walletService.Transfer(r.Context(), fromUserID, toUserID, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Transfer successful"))
}

func (a *API) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	balance, err := a.walletService.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Balance: %.2f", balance)))
}

func (a *API) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	history, err := a.walletService.GetTransactionHistory(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建交易历史响应
	response := "Transaction History:\n"
	for _, transaction := range history {
		response += fmt.Sprintf("ID: %d, Type: %s, Amount: %.2f, Time: %s\n",
			transaction.ID, transaction.TransactionType, transaction.Amount, transaction.TransactionTime.Format("2006-01-02 15:04:05"))
	}

	w.Write([]byte(response))
}

func parseRequestParams(r *http.Request) (int, float64, error) {
	userIDStr := r.URL.Query().Get("user_id")
	amountStr := r.URL.Query().Get("amount")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid user ID")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid amount")
	}

	return userID, amount, nil
}

func parseTransferRequestParams(r *http.Request) (int, int, float64, error) {
	fromUserIDStr := r.URL.Query().Get("from_user_id")
	toUserIDStr := r.URL.Query().Get("to_user_id")
	amountStr := r.URL.Query().Get("amount")

	fromUserID, err := strconv.Atoi(fromUserIDStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Invalid from user ID")
	}

	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Invalid to user ID")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Invalid amount")
	}

	return fromUserID, toUserID, amount, nil
}
