package models

type Transaction struct {
	UserID        string
	Type          string
	Amount        float64
	Remarks       string
	BalanceBefore float64
	BalanceAfter  float64
}
