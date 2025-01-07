package unit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
	"wallet-service/internal/model"
	"wallet-service/internal/service"
)

// 这里简单定义一个类似ErrWalletNotFound的错误变量，你需根据实际情况替换为真实定义
var ErrWalletNotFound = errors.New("wallet not found")

// 辅助函数，用于创建简单的钱包对象
func createWallet(userID int, balance float64) *model.Wallet {
	return &model.Wallet{
		UserID:      userID,
		Balance:     balance,
		LastUpdated: time.Now(),
	}
}

// MockWalletRepository 结构体用于模拟WalletRepository接口的实现
type MockWalletRepository struct {
	getWalletFunc             func(ctx context.Context, userID int) (*model.Wallet, error)
	updateWalletBalanceFunc   func(ctx context.Context, userID int, amount float64) error
	insertTransactionFunc     func(ctx context.Context, transaction model.Transaction) error
	insertWallet              func(ctx context.Context, wallet model.Wallet) error
	getTransactionHistoryFunc func(ctx context.Context, userID int) ([]model.Transaction, error)
}

// GetWallet 方法实现了WalletRepository接口的GetWallet方法，通过调用内部的函数来获取钱包信息
func (m *MockWalletRepository) GetWallet(ctx context.Context, userID int) (*model.Wallet, error) {
	if m.getWalletFunc != nil {
		return m.getWalletFunc(ctx, userID)
	}
	return nil, nil
}

// UpdateWalletBalance 方法实现了WalletRepository接口的UpdateWalletBalance方法，通过调用内部的函数来更新钱包余额
func (m *MockWalletRepository) UpdateWalletBalance(ctx context.Context, userID int, amount float64) error {
	if m.updateWalletBalanceFunc != nil {
		return m.updateWalletBalanceFunc(ctx, userID, amount)
	}
	return nil
}

// InsertTransaction 方法实现了WalletRepository接口的InsertTransaction方法，通过调用内部的函数来插入交易记录
func (m *MockWalletRepository) InsertTransaction(ctx context.Context, transaction model.Transaction) error {
	if m.insertTransactionFunc != nil {
		return m.insertTransactionFunc(ctx, transaction)
	}
	return nil
}

// InsertWallet 方法实现了WalletRepository接口的InsertWallet方法，通过调用内部的函数来插入钱包信息
func (m *MockWalletRepository) InsertWallet(ctx context.Context, wallet model.Wallet) error {
	if m.insertWallet != nil {
		return m.insertWallet(ctx, wallet)
	}
	return nil
}

// GetTransactionHistory 方法实现了WalletRepository接口的GetTransactionHistory方法，通过调用内部的函数来获取交易历史记录
func (m *MockWalletRepository) GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	if m.getTransactionHistoryFunc != nil {
		return m.getTransactionHistoryFunc(ctx, userID)
	}
	return nil, nil
}

// 测试存款功能
func TestWalletService_Deposit(t *testing.T) {
	// 模拟获取钱包不存在（即需要创建新钱包）的情况
	mockRepo := &MockWalletRepository{
		getWalletFunc: func(ctx context.Context, userID int) (*model.Wallet, error) {
			return nil, ErrWalletNotFound
		},
		insertWallet: func(ctx context.Context, wallet model.Wallet) error {
			return nil
		},
	}

	walletService := service.NewWalletService(mockRepo)

	// 模拟插入新钱包和插入交易记录都成功的情况
	mockRepo.insertTransactionFunc = func(ctx context.Context, transaction model.Transaction) error {
		return nil
	}

	err := walletService.Deposit(context.Background(), 1, 100.00)
	if err != nil {
		t.Errorf("存款时预期无错误，实际错误：%v", err)
	}

	// 模拟获取钱包时出错的情况
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		return nil, errors.New("模拟获取钱包出错")
	}
	err = walletService.Deposit(context.Background(), 1, 100.00)
	if err == nil {
		t.Errorf("获取钱包出错时，预期应该返回错误，实际无错误")
	}

	// 模拟插入新钱包失败的情况
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		return nil, ErrWalletNotFound
	}
	mockRepo.insertWallet = func(ctx context.Context, wallet model.Wallet) error {
		return errors.New("模拟插入新钱包失败")
	}
	err = walletService.Deposit(context.Background(), 1, 100.00)
	if err == nil {
		t.Errorf("插入新钱包失败时，预期应该返回错误，实际无错误")
	}

	// 模拟插入交易记录失败的情况
	mockRepo.insertTransactionFunc = func(ctx context.Context, transaction model.Transaction) error {
		return errors.New("模拟插入交易记录失败")
	}
	err = walletService.Deposit(context.Background(), 1, 100.00)
	if err == nil {
		t.Errorf("插入交易记录失败时，预期应该返回错误，实际无错误")
	}
}

// 测试取款功能
func TestWalletService_Withdraw(t *testing.T) {
	// 模拟获取钱包成功且余额足够的情况
	wallet := createWallet(1, 200.00)
	mockRepo := &MockWalletRepository{
		getWalletFunc: func(ctx context.Context, userID int) (*model.Wallet, error) {
			if userID == 1 {
				return wallet, nil
			}
			return nil, nil
		},
	}

	walletService := service.NewWalletService(mockRepo)

	// 模拟更新钱包余额和插入交易记录都成功的情况
	mockRepo.updateWalletBalanceFunc = func(ctx context.Context, userID int, amount float64) error {
		return nil
	}
	mockRepo.insertTransactionFunc = func(ctx context.Context, transaction model.Transaction) error {
		return nil
	}

	err := walletService.Withdraw(context.Background(), 1, 50.00)
	if err != nil {
		t.Errorf("取款时预期无错误，实际错误：%v", err)
	}

	// 模拟获取钱包时出错的情况
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		return nil, errors.New("模拟获取钱包出错")
	}
	err = walletService.Withdraw(context.Background(), 1, 50.00)
	if err == nil {
		t.Errorf("获取钱包出错时，预期 should 返回错误，实际无错误")
	}

	// 模拟钱包余额不足的情况
	wallet.Balance = 30.00
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		return wallet, nil
	}
	err = walletService.Withdraw(context.Background(), 1, 50.00)
	if err == nil {
		t.Errorf("钱包余额不足时，预期 should 返回错误，实际无错误")
	}

	// 模拟更新钱包余额失败的情况
	wallet.Balance = 200.00
	mockRepo.updateWalletBalanceFunc = func(ctx context.Context, userID int, amount float64) error {
		return errors.New("模拟更新钱包余额失败")
	}
	err = walletService.Withdraw(context.Background(), 1, 50.00)
	if err == nil {
		t.Errorf("更新钱包余额失败时，预期 should 返回错误，实际无错误")
	}

	// 模拟插入交易记录失败的情况
	mockRepo.updateWalletBalanceFunc = func(ctx context.Context, userID int, amount float64) error {
		return nil
	}
	mockRepo.insertTransactionFunc = func(ctx context.Context, transaction model.Transaction) error {
		return errors.New("模拟插入交易记录失败")
	}
	err = walletService.Withdraw(context.Background(), 1, 50.00)
	if err == nil {
		t.Errorf("插入交易记录失败时，预期 should 返回错误，实际无错误")
	}
}

// 测试转账功能
func TestWalletService_Transfer(t *testing.T) {
	// 模拟获取转出钱包和转入钱包成功，且余额足够的情况
	fromWallet := createWallet(1, 200.00)
	toWallet := createWallet(2, 100.00)
	mockRepo := &MockWalletRepository{
		getWalletFunc: func(ctx context.Context, userID int) (*model.Wallet, error) {
			log.Printf("In getWalletFunc, userID: %d", userID)
			if userID == 1 {
				return fromWallet, nil
			}
			return nil, nil
		},
	}

	walletService := service.NewWalletService(mockRepo)

	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		if userID == 2 {
			return toWallet, nil
		}
		return nil, nil
	}

	// 模拟更新双方钱包余额和插入双方交易记录都成功的情况
	mockRepo.updateWalletBalanceFunc = func(ctx context.Context, userID int, amount float64) error {
		return nil
	}
	mockRepo.insertTransactionFunc = func(ctx context.Context, transaction model.Transaction) error {
		return nil
	}

	err := walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err != nil {
		t.Errorf("转账时预期无错误，实际错误：%v", err)
	}

	// 模拟获取转出钱包时出错的情况
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		return nil, errors.New("模拟获取转出钱包出错")
	}
	err = walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err == nil || !errors.Is(err, errors.New("模拟获取转出钱包出错")) {
		t.Errorf("获取转出钱包出错时，预期 should 返回错误，实际无错误")
	}

	// 模拟获取转入钱包时出错的情况
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		if userID == 2 {
			return nil, errors.New("模拟获取转入钱包出错")
		}
		return nil, nil
	}
	err = walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err == nil || !errors.Is(err, errors.New("模拟获取转入钱包出错")) {
		t.Errorf("获取转入钱包出错时，预期 should 返回错误，实际无错误")
	}

	// 模拟转出钱包余额不足的情况
	fromWallet.Balance = 30.00
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		if userID == 2 {
			return toWallet, nil
		}
		return fromWallet, nil
	}
	err = walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err == nil || !errors.Is(err, fmt.Errorf("Insufficient balance")) {
		t.Errorf("转出钱包余额不足时，预期 should 返回错误，实际无错误")
	}

	// 模拟更新转出钱包余额失败的情况
	fromWallet.Balance = 200.00
	mockRepo.updateWalletBalanceFunc = func(ctx context.Context, userID int, amount float64) error {
		return errors.New("模拟更新转出钱包余额失败")
	}
	err = walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err == nil || !errors.Is(err, errors.New("模拟更新转出钱包余额失败")) {
		t.Errorf("更新转出钱包余额失败时，预期 should 返回错误，实际无错误")
	}

	// 模拟更新转入钱包余额失败的情况
	mockRepo.updateWalletBalanceFunc = func(ctx context.Context, userID int, amount float64) error {
		return errors.New("模拟更新转入钱包余额失败")
	}
	err = walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err == nil || !errors.Is(err, errors.New("模拟更新转入钱包余额失败")) {
		t.Errorf("更新转入钱包余额失败时，预期 should 返回错误，实际无错误")
	}

	// 模拟插入转出交易记录失败的情况
	mockRepo.updateWalletBalanceFunc = func(ctx context.Context, userID int, amount float64) error {
		return nil
	}
	mockRepo.insertTransactionFunc = func(ctx context.Context, transaction model.Transaction) error {
		return errors.New("模拟插入转出交易记录失败")
	}
	err = walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err == nil || !errors.Is(err, errors.New("模拟插入转出交易记录失败")) {
		t.Errorf("插入转出交易记录失败时，预期 should 返回错误，实际无错误")
	}

	// 模拟插入转入交易记录失败的情况
	mockRepo.insertTransactionFunc = func(ctx context.Context, transaction model.Transaction) error {
		return errors.New("模拟插入转入交易记录失败")
	}
	err = walletService.Transfer(context.Background(), 1, 2, 50.00)
	if err == nil || !errors.Is(err, errors.New("模拟插入转入交易记录失败")) {
		t.Errorf("插入转入交易记录失败时，预期 should 返回错误，实际无错误")
	}
}

// 测试获取余额功能
func TestWalletService_GetBalance(t *testing.T) {
	// 模拟获取钱包成功的情况
	wallet := createWallet(1, 200.00)
	mockRepo := &MockWalletRepository{
		getWalletFunc: func(ctx context.Context, userID int) (*model.Wallet, error) {
			if userID == 1 {
				return wallet, nil
			}
			return nil, nil
		},
	}

	walletService := service.NewWalletService(mockRepo)

	balance, err := walletService.GetBalance(context.Background(), 1)
	if err != nil {
		t.Errorf("获取余额时预期无错误，实际错误：%v", err)
	}
	if balance != 200.00 {
		t.Errorf("获取的余额值不正确，预期为200.00，实际为：%v", balance)
	}

	// 模拟获取钱包失败的情况
	mockRepo.getWalletFunc = func(ctx context.Context, userID int) (*model.Wallet, error) {
		return nil, errors.New("模拟获取钱包出错")
	}
	balance, err = walletService.GetBalance(context.Background(), 1)
	if err == nil {
		t.Errorf("获取钱包出错时，预期 should 返回错误，实际无错误")
	}
}

// 测试获取交易历史功能
func TestWalletService_GetTransactionHistory(t *testing.T) {
	// 模拟获取交易历史成功的情况
	transactions := []model.Transaction{
		// 这里可以添加一些模拟的交易记录示例
	}
	mockRepo := &MockWalletRepository{
		getTransactionHistoryFunc: func(ctx context.Context, userID int) ([]model.Transaction, error) {
			return transactions, nil
		},
	}

	walletService := service.NewWalletService(mockRepo)

	history, err := walletService.GetTransactionHistory(context.Background(), 1)
	if err != nil {
		t.Errorf("获取交易历史时预期无错误，实际错误：%v", err)
	}
	// 这里添加对history变量的使用逻辑，例如检查返回的交易历史记录数量是否符合预期
	if len(history) != len(transactions) {
		t.Errorf("获取的交易历史记录数量不正确，预期为 %d，实际为 %d", len(transactions), len(history))
	}

	// 模拟获取交易历史失败的情况
	mockRepo.getTransactionHistoryFunc = func(ctx context.Context, userID int) ([]model.Transaction, error) {
		return nil, errors.New("模拟获取交易历史出错")
	}
	history, err = walletService.GetTransactionHistory(context.Background(), 1)
	if err == nil {
		t.Errorf("获取交易历史出错时，预期 should 返回错误，实际无错误")
	}
}
