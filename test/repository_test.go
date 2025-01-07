package unit

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
	"wallet-service/internal/model"
	"wallet-service/internal/repository/postgres"
)

// 测试获取钱包功能
func TestPostgresRepository_GetWallet(t *testing.T) {
	// 创建模拟数据库连接和对象
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresRepository(db)

	// 模拟查询钱包成功的情况
	rows := sqlmock.NewRows([]string{"user_id", "balance", "last_updated"}).
		AddRow(1, 100.00, time.Now())
	mock.ExpectQuery("SELECT user_id, balance, last_updated FROM wallets WHERE user_id = \\$1").
		WithArgs(1).WillReturnRows(rows)

	wallet, err := repo.GetWallet(context.Background(), 1)
	if err != nil {
		t.Errorf("获取钱包时预期无错误，实际错误：%v", err)
	}
	if wallet.UserID != 1 || wallet.Balance != 100.00 {
		t.Errorf("预期钱包用户ID为1，余额为100.00，实际：%+v", wallet)
	}

	// 验证所有期望的操作都被执行
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的期望：%v", err)
	}
}

// 测试更新钱包余额功能
func TestPostgresRepository_UpdateWalletBalance(t *testing.T) {
	// 创建模拟数据库连接和对象
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresRepository(db)

	// 模拟更新钱包余额成功的情况
	mock.ExpectExec("UPDATE wallets SET balance = balance + \\$1, last_updated = \\$2 WHERE user_id = \\$3").
		WithArgs(50.00, time.Now(), 1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateWalletBalance(context.Background(), 1, 50.00)
	if err != nil {
		t.Errorf("更新钱包余额时预期无错误，实际错误：%v", err)
	}

	// 验证所有期望的操作都被执行
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的期望：%v", err)
	}
}

// 测试插入交易记录功能
func TestPostgresRepository_InsertTransaction(t *testing.T) {
	// 创建模拟数据库连接和对象
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresRepository(db)

	// 模拟插入交易记录成功的情况
	mock.ExpectExec("INSERT INTO transactions \\(user_id, transaction_type, amount, transaction_time\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs(1, "deposit", 100.00, time.Now()).WillReturnResult(sqlmock.NewResult(0, 1))

	transaction := model.Transaction{
		UserID:          1,
		TransactionType: "deposit",
		Amount:          100.00,
		TransactionTime: time.Now(),
	}
	err = repo.InsertTransaction(context.Background(), transaction)
	if err != nil {
		t.Errorf("插入交易记录时预期无错误，实际错误：%v", err)
	}

	// 验证所有期望的操作都被执行
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的期望：%v", err)
	}
}

// 测试获取交易历史功能
func TestPostgresRepository_GetTransactionHistory(t *testing.T) {
	// 创建模拟数据库连接和对象
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresRepository(db)

	// 模拟查询交易历史成功的情况
	rows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "transaction_time"}).
		AddRow(1, 1, "deposit", 100.00, time.Now()).
		AddRow(2, 1, "withdrawal", 50.00, time.Now())
	mock.ExpectQuery("SELECT id, user_id, transaction_type, amount, transaction_time FROM transactions WHERE user_id = \\$1 ORDER BY transaction_time DESC").
		WithArgs(1).WillReturnRows(rows)

	history, err := repo.GetTransactionHistory(context.Background(), 1)
	if err != nil {
		t.Errorf("获取交易历史时预期无错误，实际错误：%v", err)
	}
	if len(history) != 2 {
		t.Errorf("预期交易历史有2条记录，实际：%d", len(history))
	}

	// 验证所有期望的操作都被执行
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的期望：%v", err)
	}
}
