package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/nongpal/Palge-Backend/internal/account"
	"github.com/nongpal/Palge-Backend/internal/testutil"
)

func TestWithdrawAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "Test Deposit")

	var blc int64 = 1_000_000
	var amountToWithdraw int64 = 1_000
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

	_, err = model.Withdraw(ctx, acc.ID, acc.UserID, amountToWithdraw)
	if err != nil {
		t.Fatalf("Withdaw() err = %v", err)
	}

	got, err := model.Get(ctx, acc.ID, acc.UserID)
	if err != nil {
		t.Fatalf("Get() err = %v", err)
	}

	if got.Balance != acc.Balance-amountToWithdraw {
		t.Errorf("Expected balance = %d, got %d", acc.Balance-amountToWithdraw, got.Balance)
	}
}
