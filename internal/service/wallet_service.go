package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
	"wallet-service/internal/model"
	"wallet-service/internal/repository/interface"
)

// walletServiceImpl 结构体实现了WalletService接口
type walletServiceImpl struct {
	repo _interface.WalletRepository
}

// NewWalletService 创建并返回一个WalletService实例
func NewWalletService(repo _interface.WalletRepository) WalletService {
	return &walletServiceImpl{repo: repo}
}

// handleWalletNotFoundError 辅助函数，统一处理钱包不存在的错误情况
func (s *walletServiceImpl) handleWalletNotFoundError(userID int, err error) error {
	if errors.Is(err, _interface.ErrWalletNotFound) {
		return fmt.Errorf("Wallet not found for user ID %d", userID)
	}
	return err
}

// Deposit 实现存款功能
func (s *walletServiceImpl) Deposit(ctx context.Context, userID int, amount float64) error {
	if amount <= 0 {
		logrus.Errorf("Invalid deposit amount: %f for user ID: %d", amount, userID)
		return fmt.Errorf("Invalid deposit amount")
	}

	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return s.handleWalletNotFoundError(userID, err)
	}
	if wallet == nil {
		newWallet := model.Wallet{
			UserID:      userID,
			Balance:     amount,
			LastUpdated: time.Now(),
		}
		err = s.repo.InsertWallet(ctx, newWallet)
		if err != nil {
			logrus.Errorf("Error creating new wallet with initial deposit for user ID %d: %v", userID, err)
			return err
		}
	} else {
		logrus.Debugf("Going to update wallet balance for user ID %d. Current balance: %f, Deposit amount: %f", userID, wallet.Balance, amount)
		err = s.repo.UpdateWalletBalance(ctx, userID, amount)
		if err != nil {
			logrus.Errorf("Error updating wallet balance for user ID %d: %v", userID, err)
			return err
		}
		logrus.Debugf("Wallet balance updated successfully for user ID %d. New balance: %f", userID, wallet.Balance+amount)
	}

	// 记录交易
	transaction := model.Transaction{
		UserID:          userID,
		TransactionType: "deposit",
		Amount:          amount,
		TransactionTime: time.Now(),
	}
	err = s.repo.InsertTransaction(ctx, transaction)
	if err != nil {
		logrus.Errorf("Error inserting deposit transaction for user ID %d: %v", userID, err)
		return err
	}

	logrus.Infof("Deposit successful for user ID %d. New balance: %f", userID, wallet.Balance+amount)
	return nil
}

// Withdraw 实现取款功能
func (s *walletServiceImpl) Withdraw(ctx context.Context, userID int, amount float64) error {
	if amount <= 0 {
		logrus.Errorf("Invalid withdrawal amount: %f for user ID: %d", amount, userID)
		return fmt.Errorf("Invalid withdrawal amount")
	}

	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return s.handleWalletNotFoundError(userID, err)
	}
	if wallet == nil {
		logrus.Errorf("Wallet not found for user ID %d", userID)
		return fmt.Errorf("Wallet not found")
	}

	if wallet.Balance < amount {
		logrus.Errorf("Insufficient balance for user ID %d. Current balance: %f, Withdrawal amount: %f", userID, wallet.Balance, amount)
		return fmt.Errorf("Insufficient balance")
	}

	err = s.repo.UpdateWalletBalance(ctx, userID, -amount)
	if err != nil {
		logrus.Errorf("Error updating wallet balance during withdrawal for user ID %d: %v", userID, err)
		return err
	}

	// 记录交易
	transaction := model.Transaction{
		UserID:          userID,
		TransactionType: "withdrawal",
		Amount:          amount,
		TransactionTime: time.Now(),
	}
	err = s.repo.InsertTransaction(ctx, transaction)
	if err != nil {
		logrus.Errorf("Error inserting withdrawal transaction for userID %d: %v", userID, err)
		return err
	}

	logrus.Infof("Withdrawal successful for user ID %d. New balance: %f", userID, wallet.Balance-amount)
	return nil
}

// Transfer 实现转账功能
func (s *walletServiceImpl) Transfer(ctx context.Context, fromUserID, toUserID int, amount float64) error {
	if amount <= 0 {
		logrus.Errorf("Invalid transfer amount: %f from user ID %d to user ID %d", amount, fromUserID, toUserID)
		return fmt.Errorf("Invalid transfer amount")
	}

	// 获取转出钱包
	fromWallet, err := s.repo.GetWallet(ctx, fromUserID)
	if err != nil {
		logrus.Errorf("Error getting from wallet: %v", err)
		return s.handleWalletNotFoundError(fromUserID, err)
	}
	if fromWallet == nil {
		logrus.Errorf("From wallet not found for user ID %d", fromUserID)
		logrus.Debugf("GetWallet returned nil for fromUserID %d. Check if this is correct in the test scenario.", fromUserID)
		return fmt.Errorf("From wallet not found")
	}
	logrus.Debugf("FromWallet details: UserID: %d, Balance: %f, LastUpdated: %v", fromWallet.UserID, fromWallet.Balance, fromWallet.LastUpdated)

	// 获取转入钱包
	toWallet, err := s.repo.GetWallet(ctx, toUserID)
	if err != nil {
		logrus.Errorf("Error getting to wallet: %v", err)
		return s.handleWalletNotFoundError(toUserID, err)
	}
	if toWallet == nil {
		logrus.Errorf("To wallet not found for user ID %d", toUserID)
		logrus.Debugf("GetWallet returned nil for toUserID %d. Check if this is correct in the test scenario.", toUserID)
		return fmt.Errorf("To wallet not found")
	}
	logrus.Debugf("ToWallet details: UserID: %d, Balance: %f, LastUpdated: %v", toWallet.UserID, toWallet.Balance, toWallet.LastUpdated)

	// 检查转出钱包余额是否足够
	if fromWallet.Balance < amount {
		logrus.Errorf("Insufficient balance for from user ID %d. Current balance: %f, Transfer amount: %f", fromUserID, fromWallet.Balance, amount)
		logrus.Debugf("FromWallet balance details for insufficient balance check: %+v", fromWallet)
		return fmt.Errorf("Insufficient balance")
	}

	// 扣除转出钱包金额
	err = s.repo.UpdateWalletBalance(ctx, fromUserID, -amount)
	if err != nil {
		logrus.Errorf("Error updating from wallet balance during transfer for user ID %d: %v. Returning this error to the test case.", fromUserID, err)
		return err
	}

	// 增加转入钱包金额
	err = s.repo.UpdateWalletBalance(ctx, toUserID, amount)
	if err != nil {
		logrus.Errorf("Error updating to wallet balance during transfer for user ID %d: %v. Returning this error to the test case.", toUserID, err)
		return err
	}

	// 记录转出交易
	fromTransaction := model.Transaction{
		UserID:          fromUserID,
		TransactionType: "transfer_out",
		Amount:          amount,
		TransactionTime: time.Now(),
	}
	err = s.repo.InsertTransaction(ctx, fromTransaction)
	if err != nil {
		logrus.Errorf("Error inserting transfer out transaction for user ID %d: %v. Returning this error to the test case.", fromUserID, err)
		return err
	}

	// 记录转入交易
	toTransaction := model.Transaction{
		UserID:          toUserID,
		TransactionType: "transfer_in",
		Amount:          amount,
		TransactionTime: time.Now(),
	}
	err = s.repo.InsertTransaction(ctx, toTransaction)
	if err != nil {
		logrus.Errorf("Error inserting transfer in transaction for user ID %d: %v. Returning this error to the test case.", toUserID, err)
		return err
	}

	logrus.Infof("Transfer successful from user ID %d to user ID %d. Transfer amount: %f", fromUserID, toUserID, amount)
	return nil
}

// GetBalance 获取指定用户的钱包余额
func (s *walletServiceImpl) GetBalance(ctx context.Context, userID int) (float64, error) {
	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return 0, s.handleWalletNotFoundError(userID, err)
	}
	if wallet == nil {
		logrus.Infof("Wallet not found for user ID %d. Returning balance 0", userID)
		return 0, nil
	}

	logrus.Infof("Balance retrieved for user ID %d. Balance: %f", userID, wallet.Balance)
	return wallet.Balance, nil
}

// GetTransactionHistory 获取指定用户的交易历史记录
func (s *walletServiceImpl) GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	history, err := s.repo.GetTransactionHistory(ctx, userID)
	if err != nil {
		logrus.Errorf("Error getting transaction history for user ID %d: %v", userID, err)
		return nil, err
	}

	logrus.Infof("Transaction history retrieved for user ID %d. Number of transactions: %d", userID, len(history))
	return history, nil
}
