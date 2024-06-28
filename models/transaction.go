package models

import "time"

type (
	Transaction struct {
		TransferID    string    `json:"transfer_id"`
		Amount        float64   `json:"amount,omitempty"`
		Remarks       string    `json:"remarks,omitempty"`
		BalanceBefore float64   `json:"balance_before,omitempty"`
		BalanceAfter  float64   `json:"balance_after,omitempty"`
		Status        string    `json:"status,omitempty"`
		CreatedDate   time.Time `json:"created_date,omitempty"`
	}

	TopUp struct {
		Amount float64 `json:"amount"`
	}

	Transfer struct {
		RecipientID string  `json:"target_user"`
		Amount      float64 `json:"amount"`
		Remarks     string  `json:"remarks"`
	}

	Payment struct {
		PaymentID     string    `json:"payment_id"`
		Amount        float64   `json:"amount"`
		Remarks       string    `json:"remarks"`
		BalanceBefore float64   `json:"balance_before"`
		BalanceAfter  float64   `json:"balance_after"`
		CreatedDate   time.Time `json:"created_date"`
	}
)
