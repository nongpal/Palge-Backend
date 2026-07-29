package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/nongpal/Palge-Backend/internal/validator"
)

type Account struct {
	ID      int64  `json:"id"`
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

type AccountModel struct {
	DB *sql.DB
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

func (a *AccountModel) Insert(account *Account) error {
	query := `
		INSERT INTO accounts (owner, balance)
		VALUES ($1, $2)
		RETURNING id
	`

	return a.DB.QueryRow(query, account.Owner, account.Balance).Scan(&account.ID)
}

func (a *AccountModel) Get(id int64) (*Account, error) {
	if id < 1 {
		return nil, ErrAccountNotFound
	}
	query := `
		SELECT id, owner, balance
		FROM accounts
		WHERE id = $1
	`

	var account Account

	err := a.DB.QueryRow(query, id).Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrAccountNotFound
		default:
			return nil, err
		}
	}

	return &account, nil
}

func (a *AccountModel) GetAll() ([]*Account, error) {
	query := `
	SELECT id, owner, balance
	FROM accounts
	ORDER BY id ASC
	`

	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var accounts []*Account

	for rows.Next() {
		var account Account

		err := rows.Scan(
			&account.ID, &account.Owner, &account.Balance,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *AccountModel) Deposit(id int64, amount int64) (*Account, error) {
	query := `
	UPDATE accounts
	SET balance = balance + $1
	WHERE id = $2
	RETURNING id, owner, balance
	`

	account := &Account{}

	err := a.DB.QueryRow(query, amount, id).Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (a *AccountModel) Withdraw(id int64, amount int64) (*Account, error) {
	query := `
	UPDATE accounts
	SET balance = balance - $1
	WHERE id = $2
	RETURNING id, owner, balance
	`

	account, err := a.Get(id)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	account = &Account{}

	err = a.DB.QueryRow(query, amount, id).Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	return account, nil
}

func (m *AccountModel) Transfer(from, to, amount int64) (*Account, *Account, error) {
	tx, err := m.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	defer tx.Rollback()

	getAccount := func(id int64) (*Account, error) {
		query := `
		SELECT id, owner, balance 
		FROM accounts 
		WHERE id = $1 
		FOR UPDATE
		`

		account := &Account{}
		err := tx.QueryRow(query,
			id).Scan(&account.ID, &account.Owner, &account.Balance)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrAccountNotFound
			}
			return nil, err
		}

		return account, nil
	}

	sender, err := getAccount(from)
	if err != nil {
		return nil, nil, err
	}
	receiver, err := getAccount(to)
	if err != nil {
		return nil, nil, err
	}

	if sender.Balance < amount {
		return nil, nil, ErrInsufficientBalance
	}

	debitQuery := `
	UPDATE accounts
	SET balance = balance - $1
	WHERE id = $2
	RETURNING id, owner, balance
`

	err = tx.QueryRow(debitQuery, amount, sender.ID).Scan(
		&sender.ID,
		&sender.Owner,
		&sender.Balance,
	)
	if err != nil {
		return nil, nil, err
	}

	creditQuery := `
	UPDATE accounts
	SET balance = balance + $1
	WHERE id = $2
	RETURNING id, owner, balance
`

	err = tx.QueryRow(creditQuery, amount, receiver.ID).Scan(
		&receiver.ID,
		&receiver.Owner,
		&receiver.Balance,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return sender, receiver, nil
}
