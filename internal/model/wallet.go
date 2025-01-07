package model

import "time"

type Wallet struct {
	UserID      int       `json:"user_id"`
	Balance     float64   `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
}

type Transaction struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	TransactionTime time.Time `json:"transaction_time"`
}
