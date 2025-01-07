package repository

import (
	"database/sql"
	"wallet-service/internal/repository/interface"
	"wallet-service/internal/repository/postgres"
)

func NewRepository(db *sql.DB) _interface.WalletRepository {
	return postgres.NewPostgresRepository(db)
}
