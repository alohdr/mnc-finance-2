package errorMessage

import "errors"

var (
	ErrFailedRegister      = errors.New("Failed to register user")
	ErrUserExist           = errors.New("Phone Number already registered")
	ErrFailedPayment       = errors.New("Failed to make payment")
	ErrBalancePayment      = errors.New("Balance is not enough")
	ErrInternalServerError = errors.New("internal server error")
)
