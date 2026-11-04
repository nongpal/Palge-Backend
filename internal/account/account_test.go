package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/nongpal/Palge-Backend/internal/account"
	"github.com/nongpal/Palge-Backend/internal/apperr"
	"github.com/nongpal/Palge-Backend/internal/testutil"
)

func TestAccountInsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "test insert")

	var blc int64 = 1000000
	acc := account.Account{
		UserID:  u.ID,
		Balance: blc,
	}

	err := model.Insert(ctx, &acc)

	if err != nil {
		t.Fatalf("Inser() err = %v", err)
	}

	if acc.UserID != u.ID {
		t.Errorf("Expected account belongs to user %d, got = %d", u.ID, acc.UserID)
	}

	if acc.ID <= 0 {
		t.Errorf("Expected id > 0, got = %d", acc.ID)
	}

	if acc.Balance != blc {
		t.Errorf("Expected balance = %d, got = %d", blc, acc.Balance)
	}

	if acc.Status != "ACTIVE" {
		t.Errorf("Expected account 'ACTIVE' after been created, got = %s", acc.Status)
	}

	if acc.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}

	if acc.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not set")
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

}

func TestAccountGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "test get")

	var blc int64 = 1000000
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

	got, err := model.Get(ctx, acc.ID, acc.UserID)
	if err != nil {
		if apperr.Is(err, apperr.CodeRecordNotFound) {
			t.Errorf("Expected existing account, got = %q", apperr.CodeRecordNotFound)
		}
		t.Fatalf("Get() err = %v", err)
	}

	if got.ID != acc.ID {
		t.Errorf("Expected id %d, got = %d", acc.ID, got.ID)
	}

	if got.Balance != blc {
		t.Errorf("Expected balance = %d, got = %d", blc, got.Balance)
	}
}

func TestGetNegativeID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "test insert")

	var blc int64 = 1000000
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

	got, err := model.Get(ctx, -1, acc.UserID)

	if got != nil {
		t.Error("Expected nil because not found")
	}

	if !apperr.Is(err, apperr.CodeRecordNotFound) {
		t.Errorf("Expected error = %v, got %v", apperr.CodeRecordNotFound, err)
	}
}

func TestGetNegativeUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "test insert")

	var blc int64 = 1000000
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

	got, err := model.Get(ctx, acc.ID, -1)

	if got != nil {
		t.Error("Expected nil because not found")
	}

	if !apperr.Is(err, apperr.CodeRecordNotFound) {
		t.Errorf("Expected error = %v, got %v", apperr.CodeRecordNotFound, err)
	}
}

func TestGetInvalidID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := testutil.NewTestDB(t)
	model := account.Model{DB: db}

	u := testutil.CreateTestUser(t, db, "test insert")

	var blc int64 = 1000000
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

	got, err := model.Get(ctx, 2, acc.UserID)

	if got != nil {
		t.Error("Expected nil because not found")
	}

	if !apperr.Is(err, apperr.CodeRecordNotFound) {
		t.Errorf("Expected error = %v, got %v", apperr.CodeRecordNotFound, err)
	}
}
