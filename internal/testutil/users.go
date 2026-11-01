package testutil

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func CreateTestUser(t testing.TB, db *sql.DB, name string) int64 {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email := name + "@example.com"
	var userID int64
	err := db.QueryRowContext(ctx,
		`
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`, name, email, []byte("password-hash"), true,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		_, err := db.ExecContext(
			cleanupCtx,
			"DELETE FROM users WHERE id = $1", userID,
		)

		if err != nil {
			t.Fatalf("cleanup users %d: %v", userID, err)
		}
	})

	return userID
}
