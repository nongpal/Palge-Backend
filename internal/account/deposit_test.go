package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/nongpal/Palge-Backend/internal/account"
	"github.com/nongpal/Palge-Backend/internal/apperr"
	"github.com/nongpal/Palge-Backend/internal/testutil"
)

func TestDepositAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "Test Deposit")

	var blc int64 = 1000000
	var amountToDeposit int64 = 1000000
	acc := account.Account{
		UserID:  u.ID,
		Balance: blc,
	}

	err := model.Insert(ctx, &acc)

	if err != nil {
		t.Fatalf("Inser() err = %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		_, err := db.ExecContext(
			cleanupCtx,
			"DELETE FROM accounts WHERE id = $1", acc.ID,
		)

		if err != nil {
			t.Errorf("cleanup account %d: %v", acc.ID, err)
		}
	})

	got, err := model.Deposit(ctx, acc.ID, acc.UserID, amountToDeposit)
	if err != nil {
		t.Fatalf("Deposit() err = %v", err)
	}

	if got.Balance != acc.Balance+amountToDeposit {
		t.Errorf("Expected balance = %d, got %d", acc.Balance+amountToDeposit, got.Balance)
	}

}

func TestDepositAccountToOtherUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	owner := testutil.CreateTestUser(t, db, "Test Deposit")
	other := testutil.CreateTestUser(t, db, "Test Other User")

	var initialBalance int64 = 1000000
	var amountToDeposit int64 = 1000000

	acc := account.Account{
		UserID:  owner.ID,
		Balance: initialBalance,
	}

	err := model.Insert(ctx, &acc)

	if err != nil {
		t.Fatalf("Inser() err = %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		_, err := db.ExecContext(
			cleanupCtx,
			"DELETE FROM accounts WHERE id = $1", acc.ID,
		)

		if err != nil {
			t.Errorf("cleanup account %d: %v", acc.ID, err)
		}
	})

	got, err := model.Deposit(ctx, acc.ID, other.ID, amountToDeposit)
	if !apperr.Is(err, apperr.CodeRecordNotFound) {
		t.Fatalf("Deposit() err = %v, want %v", err, apperr.CodeRecordNotFound)
	}

	if got != nil {
		t.Fatalf("Deposit() account = %+v, want nil", got)
	}

	after, err := model.Get(ctx, acc.ID, owner.ID)
	if err != nil {
		t.Fatalf("Get() err = %v", err)
	}

	if after.Balance != initialBalance {
		t.Errorf(
			"Balance = %d, want %d",
			after.Balance,
			initialBalance,
		)
	}
}
