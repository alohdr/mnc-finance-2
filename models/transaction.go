package models

type (
	Transaction struct {
		UserID        string
		Type          string
		Amount        float64
		Remarks       string
		BalanceBefore float64
		BalanceAfter  float64
	}

	TopUp struct {
		UserID  string  `json:"user_id"`
		Amount  float64 `json:"amount"`
		Remarks string  `json:"remarks"`
	}
)
