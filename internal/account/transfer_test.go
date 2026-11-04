package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/nongpal/Palge-Backend/internal/account"
	"github.com/nongpal/Palge-Backend/internal/testutil"
)

func TestTransferAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "Test Deposit")

	var blc int64 = 1_000_000
	var amountToTransfer int64 = 1_000
	sender := account.Account{
		UserID:  u.ID,
		Balance: blc,
	}

	err := model.Insert(ctx, &sender)
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
			"DELETE FROM accounts WHERE id = $1", sender.ID,
		)

		if err != nil {
			t.Errorf("cleanup account %d: %v", sender.ID, err)
		}
	})

	receiver := account.Account{
		UserID:  u.ID,
		Balance: blc,
	}

	err = model.Insert(ctx, &receiver)
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
			"DELETE FROM accounts WHERE id = $1", receiver.ID,
		)

		if err != nil {
			t.Errorf("cleanup account %d: %v", receiver.ID, err)
		}
	})

	_, _, err = model.Transfer(ctx, u.ID, sender.ID, receiver.ID, amountToTransfer)
	if err != nil {
		t.Fatalf("Expected Transfer() err = %v", err)
	}

	got_sender, err := model.Get(ctx, sender.ID, sender.UserID)
	if err != nil {
		t.Fatalf("Get() err = %v", err)
	}

	if got_sender.Balance != sender.Balance-amountToTransfer {
		t.Errorf("Expected balance = %d, got %d", sender.Balance-amountToTransfer, got_sender.Balance)
	}

	got_receiver, err := model.Get(ctx, receiver.ID, receiver.UserID)
	if err != nil {
		t.Fatalf("Get() err = %v", err)
	}

	if got_receiver.Balance != receiver.Balance+amountToTransfer {
		t.Errorf("Expected balance = %d, got %d", receiver.Balance+amountToTransfer, got_receiver.Balance)
	}
}
