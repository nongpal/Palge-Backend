package transfer

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/nongpal/Palge-Backend/internal/account"
	"github.com/nongpal/Palge-Backend/internal/apperr"
)

type Transfer struct {
	ID             int64
	UserID         int64
	SenderID       int64
	ReceiverID     int64
	Amount         int64
	IdempotencyKey string
	CreatedAt      time.Time
}

type TransferParams struct {
	UserID         int64
	From           int64
	To             int64
	Amount         int64
	IdempotencyKey string
}

type Model struct {
	DB *sql.DB
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

func (m *Model) Transfer(
	ctx context.Context,
	p TransferParams,
) (*account.Account, *account.Account, error) {
	if p.Amount <= 0 {
		return nil, nil, apperr.Field(
			apperr.CodeValidationFailed,
			"amount", "must be greater than 0",
		)
	}

	if p.From == p.To {
		return nil, nil, apperr.New(apperr.CodeSameAccountTransfer)
	}

	var sender, receiver account.Account

	err := m.execInTx(ctx, func(tx *sql.Tx) (err error) {
		firstID, secondID := p.From, p.To
		if p.From > p.To {
			firstID, secondID = p.To, p.From
		}

		queryLock := `
		SELECT id, user_id, balance, status
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

		accounts := make(map[int64]*account.Account)
		for rows.Next() {
			var acc account.Account
			if err := rows.Scan(
				&acc.ID,
				&acc.UserID,
				&acc.Balance,
				&acc.Status,
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

		accFrom, existFrom := accounts[p.From]
		accTo, existTo := accounts[p.To]

		if !existFrom || !existTo {
			return apperr.New(apperr.CodeRecordNotFound)
		}

		if accFrom.UserID != p.UserID {
			return apperr.New(apperr.CodePermissionDenied)
		}

		if accFrom.Status != "ACTIVE" || accTo.Status != "ACTIVE" {
			return apperr.New(apperr.CodeAccountInactive)
		}

		if accFrom.Balance < p.Amount {
			return apperr.New(apperr.CodeInsufficientBalance)
		}

		debitQuery := `
			UPDATE accounts
			SET balance = balance - $1,
				updated_at = NOW(),
				version = version + 1
			WHERE id = $2
			RETURNING id, user_id, balance, updated_at, version
		`
		if err := tx.QueryRowContext(ctx, debitQuery, p.Amount, p.From).Scan(
			&sender.ID,
			&sender.UserID,
			&sender.Balance,
			&sender.UpdatedAt,
			&sender.Version,
		); err != nil {
			return err
		}

		creditQuery := `
			UPDATE accounts
			SET balance = balance + $1,
				updated_at = NOW(),
				version = version + 1
			WHERE id = $2
			RETURNING id, user_id, balance, updated_at, version
		`
		if err := tx.QueryRowContext(ctx, creditQuery, p.Amount, p.To).Scan(
			&receiver.ID,
			&receiver.UserID,
			&receiver.Balance,
			&receiver.UpdatedAt,
			&receiver.Version,
		); err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO transfers (
				user_id, sender_account_id, receiver_account_id,
				amount, idempotency_key
			)
			VALUES ($1, $2, $3, $4, NULLIF($5, ''))
		`, p.UserID, p.From, p.To, p.Amount, p.IdempotencyKey)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &sender, &receiver, nil
}
