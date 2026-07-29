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

func (m *AccountModel) execInTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *AccountModel) Insert(ctx context.Context, account *Account) error {
	query := `
		INSERT INTO accounts (owner, balance)
		VALUES ($1, $2)
		RETURNING id
	`

	return m.DB.QueryRowContext(ctx, query, account.Owner, account.Balance).Scan(&account.ID)
}

func (m *AccountModel) Get(ctx context.Context, id int64) (*Account, error) {
	if id < 1 {
		return nil, ErrAccountNotFound
	}
	query := `
		SELECT id, owner, balance
		FROM accounts
		WHERE id = $1
	`

	var account Account

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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

func (m *AccountModel) Deposit(ctx context.Context, id int64, amount int64) (*Account, error) {
	query := `
	UPDATE accounts
	SET balance = balance + $1
	WHERE id = $2
	RETURNING id, owner, balance
	`

	account := &Account{}

	err := m.DB.QueryRowContext(ctx, query, amount, id).Scan(
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

func (m *AccountModel) Withdraw(ctx context.Context, id int64, amount int64) (*Account, error) {
	query := `
	UPDATE accounts
	SET balance = balance - $1
	WHERE id = $2 AND balance >= $1
	RETURNING id, owner, balance
	`

	var account Account
	err := m.DB.QueryRowContext(ctx, query, amount, id).Scan(
		&account.ID,
		&account.Owner,
		&account.Balance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var exists bool
			checkErr := m.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", id).Scan(&exists)
			if checkErr == nil && exists {
				return nil, ErrInsufficientBalance
			}
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

func (m *AccountModel) Transfer(ctx context.Context, from, to, amount int64) (*Account, *Account, error) {
	if from == to {
		return nil, nil, ErrSameAccountTransfer
	}
	var sender, receiver Account
	err := m.execInTx(ctx, func(tx *sql.Tx) error {
		firstID, secondID := from, to
		if from > to {
			firstID, secondID = to, from
		}

		queryLock := `
		SELECT id, owner, balance 
		FROM accounts 
		WHERE id = ($1, $2)
		ORDER BY id
		FOR UPDATE
		`
		rows, err := tx.QueryContext(ctx, queryLock, firstID, secondID)
		if err != nil {
			return err
		}

		defer rows.Close()

		accounts := make(map[int64]*Account)
		for rows.Next() {
			var acc Account
			if err := rows.Scan(
				&acc.ID,
				&acc.Owner,
				&acc.Balance,
			); err != nil {
				return err
			}
			accounts[acc.ID] = &acc
		}

		if err := rows.Err(); err != nil {
			return err
		}

		accFrom, existFrom := accounts[from]
		accTo, existTo := accounts[to]

		if !existFrom || !existTo {
			return ErrAccountNotFound
		}

		if accFrom.Balance < amount {
			return ErrInsufficientBalance
		}

		debitQuery := `
			UPDATE accounts
			SET balance = balance - $1
			WHERE id = $2
			RETURNING id, owner, balance
		`
		if err := tx.QueryRowContext(ctx, debitQuery, amount, from).Scan(
			&sender.ID,
			&sender.Owner,
			&sender.Balance,
		); err != nil {
			return err
		}

		creditQuery := `UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING id, owner, balance`
		if err := tx.QueryRowContext(ctx, creditQuery, amount, to).Scan(&receiver.ID, &receiver.Owner, &receiver.Balance); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &sender, &receiver, nil
}
