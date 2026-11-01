package user

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nongpal/Palge-Backend/internal/apperr"
	"github.com/nongpal/Palge-Backend/internal/testutil"
)

func TestUserInsert(t *testing.T) {
	db := testutil.NewTestDB(t)
	model := Model{DB: db}

	email := fmt.Sprintf("nongpaldummy-%d@example.com", time.Now().UnixNano())

	input := User{
		Name:      "Nongpal Dummy",
		Email:     email,
		Activated: true,
	}

	if err := input.Password.Set("P455w0rd"); err != nil {
		t.Fatalf("set password: %v", err)
	}

	if err := model.Insert(&input); err != nil && apperr.Is(err, apperr.CodeDuplicateEmail) {
		t.Fatalf("insert user: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := db.ExecContext(
			ctx,
			"DELETE FROM users WHERE id = $1",
			input.ID,
		); err != nil {
			t.Errorf("cleanup user %d: %v", input.ID, err)
		}
	})

	if input.ID <= 0 {
		t.Errorf("ID = %d, want > 0", input.ID)
	}

	if input.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}

	if input.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not set")
	}

	if input.Version != 1 {
		t.Errorf("Version = %d, want 1", input.Version)
	}

	var got User

	err := db.QueryRowContext(
		context.Background(),
		`
		SELECT id, created_at, updated_at, name, email, activated, version
		FROM users
		WHERE id = $1
		`,
		input.ID,
	).Scan(
		&got.ID,
		&got.CreatedAt,
		&got.UpdatedAt,
		&got.Name,
		&got.Email,
		&got.Activated,
		&got.Version,
	)
	if err != nil {
		t.Fatalf("query inserted user: %v", err)
	}

	if got.ID != input.ID {
		t.Errorf("ID = %d, want %d", got.ID, input.ID)
	}

	if got.Name != input.Name {
		t.Errorf("Name = %q, want %q", got.Name, input.Name)
	}

	if got.Email != email {
		t.Errorf("Email = %q, want %q", got.Email, email)
	}

	if !got.Activated {
		t.Error("Activated = false, want true")
	}

	if got.Version != 1 {
		t.Errorf("Version = %d, want 1", got.Version)
	}
}

func TestUserInsertDuplicateEmail(t *testing.T) {
	db := testutil.NewTestDB(t)
	model := Model{DB: db}

	email := fmt.Sprintf("nongpaldup-%d@example.com", time.Now().UnixNano())

	first := User{
		Name:      "Nongpal Dummy",
		Email:     email,
		Activated: true,
	}

	if err := first.Password.Set("P455w0rd"); err != nil {
		t.Fatalf("set password: %v", err)
	}

	if err := model.Insert(&first); err != nil {
		t.Fatalf("insert first user: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := db.ExecContext(
			ctx,
			"DELETE FROM users WHERE id = $1",
			first.ID,
		); err != nil {
			t.Errorf("cleanup user %d: %v", first.ID, err)
		}
	})

	second := User{
		Name:      "Nongpal Dummy 2",
		Email:     email,
		Activated: true,
	}

	if err := second.Password.Set("P455w0rd"); err != nil {
		t.Fatalf("set password: %v", err)
	}

	err := model.Insert(&second)
	if !apperr.Is(err, apperr.CodeDuplicateEmail) {
		t.Fatalf(
			"insert duplicate email: got error %v, want %v",
			err,
			apperr.CodeDuplicateEmail,
		)
	}
}

func TestGetByEmail(t *testing.T) {
	db := testutil.NewTestDB(t)
	model := Model{DB: db}

	email := fmt.Sprintf("nongpaldup-%d@example.com", time.Now().UnixNano())

	input := User{
		Name:      "Nongpal Dummy",
		Email:     email,
		Activated: true,
	}

	if err := input.Password.Set("P455w0rd"); err != nil {
		t.Fatalf("set password: %v", err)
	}

	if err := model.Insert(&input); err != nil {
		t.Fatalf("insert first user: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := db.ExecContext(
			ctx,
			"DELETE FROM users WHERE id = $1",
			input.ID,
		); err != nil {
			t.Errorf("cleanup user %d: %v", input.ID, err)
		}
	})

	user, err := model.GetByEmail(email)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatalf("Expected user, got nil")
	}

	if user.ID != input.ID {
		t.Errorf("ID = %d, want %d", user.ID, input.ID)
	}

	if user.Email != input.Email {
		t.Errorf("Email = %q, want %q", user.Email, input.Email)
	}

	if user.Password.hash == nil {
		t.Errorf("Expected password hash to be persisted")
	}

	user, err = model.GetByEmail("nonexistnongpal@email.com")

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if user != nil {
		t.Fatalf("Expected nil user on error, got %+v", user)
	}

	if !apperr.Is(err, apperr.CodeRecordNotFound) {
		t.Fatalf("Expected CodeRecordNotFound, got %v", err)
	}
}

func TestInsertNotNullViolations(t *testing.T) {
	db := testutil.NewTestDB(t)
	model := Model{DB: db}

	email := fmt.Sprintf("nongpaldup-%d@example.com", time.Now().UnixNano())

	input := User{
		Name:      "Nongpal Dummy",
		Email:     email,
		Activated: true,
	}

	if err := model.Insert(&input); !apperr.Is(err, apperr.CodeNotNullViolation) {
		t.Fatalf(
			"insert null password into user: got error %v, want %v",
			err,
			apperr.CodeNotNullViolation,
		)
	}
}

func TestMatchPassword(t *testing.T) {
	tpw := "P455w0rd"
	wpw := "wrongpw"

	db := testutil.NewTestDB(t)
	model := Model{DB: db}

	email := fmt.Sprintf("matchpw-%d@example.com", time.Now().UnixNano())

	input := User{
		Name:      "Nongpal Dummy",
		Email:     email,
		Activated: true,
	}

	if err := input.Password.Set(tpw); err != nil {
		t.Fatalf("set password: %v", err)
	}

	if err := model.Insert(&input); err != nil {
		t.Fatalf("insert user: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", input.ID); err != nil {
			t.Errorf("cleanup user %d: %v", input.ID, err)
		}
	})

	user, err := model.GetByEmail(email)
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}

	if user.Password.hash == nil {
		t.Fatal("expected password hash to be loaded from database")
	}

	match, err := user.Password.Match(tpw)
	if err != nil {
		t.Fatalf("match with correct password: %v", err)
	}
	if !match {
		t.Fatal("expected correct password to match")
	}

	match, err = user.Password.Match(wpw)
	if err != nil {
		t.Fatalf("match with wrong password: %v", err)
	}
	if match {
		t.Fatal("expected wrong password to NOT match")
	}
}
