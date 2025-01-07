package _interface

import (
	"context"
	"wallet-service/internal/model"
)

// ErrWalletNotFound 定义表示钱包不存在的错误常量
var ErrWalletNotFound = NewWalletNotFoundError()

// WalletNotFoundError 是自定义的错误类型，用于表示钱包不存在的情况
type WalletNotFoundError struct{}

// Error 实现了error接口的Error方法，返回错误描述信息
func (e WalletNotFoundError) Error() string {
	return "wallet not found"
}

// NewWalletNotFoundError 创建并返回一个WalletNotFoundError实例
func NewWalletNotFoundError() error {
	return WalletNotFoundError{}
}

// WalletRepository 定义了钱包相关操作的仓库接口
type WalletRepository interface {
	GetWallet(ctx context.Context, userID int) (*model.Wallet, error)
	UpdateWalletBalance(ctx context.Context, userID int, amount float64) error
	InsertWallet(ctx context.Context, wallet model.Wallet) error
	InsertTransaction(ctx context.Context, transaction model.Transaction) error
	GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error)
}
