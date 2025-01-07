package service

import (
	"context"

	"wallet-service/internal/model"
)

type WalletService interface {
	Deposit(ctx context.Context, userID int, amount float64) error
	Withdraw(ctx context.Context, userID int, amount float64) error
	Transfer(ctx context.Context, fromUserID, toUserID int, amount float64) error
	GetBalance(ctx context.Context, userID int) (float64, error)
	GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error)
}
