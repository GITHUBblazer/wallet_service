package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"
	"wallet-service/internal/model"
	_interface "wallet-service/internal/repository/interface"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) _interface.WalletRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetWallet(ctx context.Context, userID int) (*model.Wallet, error) {
	query := "SELECT user_id, balance, last_updated FROM wallets WHERE user_id = $1"
	row := r.db.QueryRowContext(ctx, query, userID)

	var wallet model.Wallet
	err := row.Scan(&wallet.UserID, &wallet.Balance, &wallet.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &wallet, nil
}

func (p *PostgresRepository) UpdateWalletBalance(ctx context.Context, userID int, amount float64) error {
	sql := "UPDATE wallets SET balance = balance + $1, last_updated = $2 WHERE user_id = $3"
	log.Printf("Actual parameters: amount=%v, time=%v, userID=%d", amount, time.Now(), userID) // 添加日志打印
	_, err := p.db.ExecContext(ctx, sql, amount, time.Now(), userID)
	return err
}

func (r *PostgresRepository) InsertTransaction(ctx context.Context, transaction model.Transaction) error {
	query := "INSERT INTO transactions (user_id, transaction_type, amount, transaction_time) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query, transaction.UserID, transaction.TransactionType, transaction.Amount, transaction.TransactionTime)
	return err
}

func (r *PostgresRepository) GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	query := "SELECT id, user_id, transaction_type, amount, transaction_time FROM transactions WHERE user_id = $1 ORDER BY transaction_time DESC"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.TransactionType, &transaction.Amount, &transaction.TransactionTime)
		if err != nil {
			return nil, err
		}
		history = append(history, transaction)
	}

	return history, nil
}

func (p *PostgresRepository) InsertWallet(ctx context.Context, wallet model.Wallet) error {
	sql := "INSERT INTO wallets (user_id, balance, last_updated) VALUES ($1, $2, $3)"
	_, err := p.db.ExecContext(ctx, sql, wallet.UserID, wallet.Balance, wallet.LastUpdated)
	return err
}
