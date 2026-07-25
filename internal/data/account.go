package data

import (
	"strings"

	"github.com/nongpal/Palge-Backend/internal/validator"
)

type Account struct {
	ID      int64  `json:"id"`
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

func ValidateAccount(v *validator.Validator, account *Account) {
	v.Check(strings.TrimSpace(account.Owner) != "", "owner", "must be provided")
	v.Check(account.Balance >= 0, "balance", "must not be negative")
}

func ValidateAmount(v *validator.Validator, amount int64) {
	v.Check(amount > 0, "amount", "must be greater than 0")
}

func ValidateWithdraw(v *validator.Validator, account *Account) {

}
