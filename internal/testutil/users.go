package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nongpal/Palge-Backend/internal/user"
)

func CreateTestUser(t testing.TB, db *sql.DB, name string) *user.User {
	t.Helper()

	model := user.Model{DB: db}

	u := &user.User{
		Name:      name,
		Email:     fmt.Sprintf("test-%s@example.com", uuid.NewString()),
		Activated: true,
	}

	if err := u.Password.Set("P455w0rd"); err != nil {
		t.Fatalf("set password: %v", err)
	}

	if err := model.Insert(u); err != nil {
		t.Fatalf("insert user: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		_, err := db.ExecContext(
			ctx,
			"DELETE FROM users WHERE id = $1", u.ID,
		)

		if err != nil {
			t.Errorf("cleanup users %d: %v", u.ID, err)
		}
	})

	return u
}
