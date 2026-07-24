package data

import "github.com/nongpal/Palge-Backend/internal/validator"

type Account struct {
	ID      int64  `json:"id"`
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

func ValidateAmount(v *validator.Validator, amount int64) {

}
