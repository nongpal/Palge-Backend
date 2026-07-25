package data

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")

	ErrInvalidAmount   = errors.New("amount must be greater than 0")
	ErrNegativeBalance = errors.New("balance must not be negative")
	ErrEmptyOwner      = errors.New("owner must be provided")

	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrSameAccountTransfer = errors.New("sender and receiver must be different")
	ErrSenderNotFound      = errors.New("sender account not found")
	ErrReceiverNotFound    = errors.New("receiver account not found")
)
