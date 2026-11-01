package account

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/nongpal/Palge-Backend/internal/apperr"
	"github.com/nongpal/Palge-Backend/internal/validator"
)

type Account struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   int64     `json:"balance"`
	Status    string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Version   int       `json:"-"`
}

type Model struct {
	DB *sql.DB
}

func ValidateAccount(v *validator.Validator, account *Account) {
	//v.Check(strings.TrimSpace(account.Owner) != "", "owner", "must be provided")
	v.Check(account.Balance >= 0, "balance", "must not be negative")
}

func ValidateAmount(v *validator.Validator, amount int64) {
	v.Check(amount > 0, "amount", "must be greater than 0")
}

func (m *Model) execInTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			err = errors.Join(err, rbErr)
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Model) Insert(ctx context.Context, account *Account) error {
	query := `
	INSERT INTO accounts (user_id, balance)
	VALUES ($1, $2)
	RETURNING id, user_id, balance, status, created_at, updated_at
	`

	return m.DB.QueryRowContext(
		ctx,
		query,
		account.UserID,
		account.Balance,
	).Scan(
		&account.ID,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
}

func (m *Model) Get(ctx context.Context, id, user_id int64) (*Account, error) {
	if id < 1 {
		return nil, apperr.New(apperr.CodeRecordNotFound)
	}
	query := `
		SELECT id, user_id, owner, balance
		FROM accounts
		WHERE id = $1 AND user_id = $2
	`

	var account Account

	err := m.DB.QueryRowContext(ctx, query, id, user_id).Scan(
		&account.ID,
		&account.Balance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.New(apperr.CodeRecordNotFound)
		}
		return nil, err
	}

	return &account, nil
}

func (a *Model) GetAll(ctx context.Context) ([]*Account, error) {
	query := `
	SELECT id, user_id, balance
	FROM accounts
	ORDER BY id ASC
	`

	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rows != nil {
			err = errors.Join(err, rows.Close())
		}
	}()

	var accounts []*Account

	for rows.Next() {
		var account Account

		err = rows.Scan(
			&account.ID, &account.UserID, &account.Balance,
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

func (m *Model) Deposit(ctx context.Context, id, user_id, amount int64) (*Account, error) {
	query := `
	UPDATE accounts
	SET balance = balance + $1
	WHERE id = $2 AND user_id = $3
	RETURNING id, user_id,  balance, updated_at
	`

	account := &Account{}

	err := m.DB.QueryRowContext(ctx, query, amount, id, user_id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.New(apperr.CodeRecordNotFound)
		}
		return nil, err
	}

	return account, nil
}

func (m *Model) Withdraw(ctx context.Context, id, user_id, amount int64) (*Account, error) {
	query := `
	UPDATE accounts
	SET balance = balance - $1
	WHERE id = $2 AND balance >= $1
	RETURNING id, user_id, balance, updated_at
	`

	var account Account
	err := m.DB.QueryRowContext(ctx, query, amount, id, user_id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var exists bool
			checkErr := m.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", id).Scan(&exists)
			if checkErr == nil && exists {
				return nil, apperr.New(apperr.CodeInsufficientBalance)
			}
			return nil, apperr.New(apperr.CodeRecordNotFound)
		}
		return nil, err
	}

	return &account, nil
}

func (m *Model) Transfer(ctx context.Context, userID, from, to, amount int64) (*Account, *Account, error) {
	if from == to {
		return nil, nil, apperr.New(apperr.CodeSameAccountTransfer)
	}
	var sender, receiver Account
	err := m.execInTx(ctx, func(tx *sql.Tx) (err error) {
		firstID, secondID := from, to
		if from > to {
			firstID, secondID = to, from
		}

		queryLock := `
		SELECT id, user_id, balance
		FROM accounts
		WHERE id in ($1, $2)
		ORDER BY id
		FOR UPDATE
		`
		rows, err := tx.QueryContext(ctx, queryLock, firstID, secondID)
		if err != nil {
			return err
		}

		defer func() {
			if rows != nil {
				err = errors.Join(
					err,
					rows.Close(),
				)
			}
		}()

		accounts := make(map[int64]*Account)
		for rows.Next() {
			var acc Account
			if err := rows.Scan(
				&acc.ID,
				&acc.UserID,
				&acc.Balance,
			); err != nil {
				return err
			}
			accounts[acc.ID] = &acc
		}

		if err := rows.Err(); err != nil {
			return err
		}

		if err := rows.Close(); err != nil {
			return err
		}

		accFrom, existFrom := accounts[from]
		_, existTo := accounts[to]

		if !existFrom || !existTo {
			return apperr.New(apperr.CodeRecordNotFound)
		}

		if accFrom.UserID != userID {
			return apperr.New(apperr.CodePermissionDenied)
		}

		if accFrom.Balance < amount {
			return apperr.New(apperr.CodeInsufficientBalance)
		}

		debitQuery := `
			UPDATE accounts
			SET balance = balance - $1
			WHERE id = $2
			RETURNING id, balance, updated_at
		`
		if err := tx.QueryRowContext(ctx, debitQuery, amount, from).Scan(
			&sender.ID,
			&sender.Balance,
			&sender.UpdatedAt,
		); err != nil {
			return err
		}

		creditQuery := `
			UPDATE accounts
			SET balance = balance + $1
			WHERE id = $2
			RETURNING id, balance, updated_at
		`
		if err := tx.QueryRowContext(ctx, creditQuery, amount, to).Scan(
			&receiver.ID,
			&receiver.Balance,
			&receiver.UpdatedAt,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &sender, &receiver, nil
}
