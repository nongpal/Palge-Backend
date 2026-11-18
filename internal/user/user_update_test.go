package user_test

import (
	"testing"

	"github.com/nongpal/Palge-Backend/internal/apperr"
)

func TestUserUpdateName(t *testing.T) {
	input, model, _ := newTestUserFixture(t)

	beforeUpdate := input.UpdatedAt
	beforeVersion := input.Version

	input.Name = "Page"
	if err := model.Update(&input); err != nil {
		t.Fatalf("update user: %v", err)
	}

	if input.Name != "Page" {
		t.Errorf("Name = %q, want %q", input.Name, "Page")
	}

	if input.Version != beforeVersion+1 {
		t.Errorf("Version = %d, want %d", input.Version, beforeVersion+1)
	}

	if !input.UpdatedAt.After(beforeUpdate) {
		t.Errorf("UpdatedAt = %v, want after %v", input.UpdatedAt, beforeUpdate)
	}

	got, err := model.GetByEmail(input.Email)
	if err != nil {
		t.Fatalf("get updated user: %v", err)
	}

	if got.Name != "Page" {
		t.Errorf("stored Name = %q, want %q", got.Name, "Page")
	}

	if got.Version != beforeVersion+1 {
		t.Errorf("stored Version = %d, want %d", got.Version, beforeVersion+1)
	}

	if !got.UpdatedAt.Equal(input.UpdatedAt) {
		t.Errorf("stored UpdatedAt = %v, want %v", got.UpdatedAt, input.UpdatedAt)
	}

	if got.Email != input.Email {
		t.Errorf("stored Email = %q, want %q", got.Email, input.Email)
	}

	if got.Activated != input.Activated {
		t.Errorf("stored Activated = %v, want %v", got.Activated, input.Activated)
	}

	match, err := got.Password.Match("P455w0rd")
	if err != nil {
		t.Fatalf("match password after update: %v", err)
	}
	if !match {
		t.Fatal("expected password to remain valid after update")
	}
}

func TestUserNotUpdate(t *testing.T) {
	input, model, _ := newTestUserFixture(t)

	beforeUpdate := input.UpdatedAt
	beforeVersion := input.Version
	originalName := input.Name

	input.Name = originalName
	if err := model.Update(&input); err != nil {
		t.Fatalf("update user: %v", err)
	}

	if !input.UpdatedAt.Equal(beforeUpdate) {
		t.Errorf("UpdatedAt = %v, want %v (unchanged when no fields updated)", input.UpdatedAt, beforeUpdate)
	}

	if input.Version != beforeVersion {
		t.Errorf("Version = %d, want %d (unchanged when no fields updated)", input.Version, beforeVersion)
	}

	got, err := model.GetByEmail(input.Email)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	if got.Name != originalName {
		t.Errorf("stored Name = %q, want %q", got.Name, originalName)
	}

	if got.Version != beforeVersion {
		t.Errorf("stored Version = %d, want %d", got.Version, beforeVersion)
	}

	if !got.UpdatedAt.Equal(beforeUpdate) {
		t.Errorf("stored UpdatedAt = %v, want %v", got.UpdatedAt, beforeUpdate)
	}
}

func TestUserUpdateButEmailDuplicated(t *testing.T) {
	user1, _, _ := newTestUserFixture(t)
	user2, model2, _ := newTestUserFixture(t)

	originalEmail := user2.Email
	beforeUpdate := user2.UpdatedAt
	beforeVersion := user2.Version

	user2.Email = user1.Email
	err := model2.Update(&user2)
	if !apperr.Is(err, apperr.CodeDuplicateEmail) {
		t.Fatalf("update user with duplicate email: got %v, want %v", err, apperr.CodeDuplicateEmail)
	}

	if !user2.UpdatedAt.Equal(beforeUpdate) {
		t.Errorf("UpdatedAt = %v, want %v (unchanged on failed update)", user2.UpdatedAt, beforeUpdate)
	}

	if user2.Version != beforeVersion {
		t.Errorf("Version = %d, want %d (unchanged on failed update)", user2.Version, beforeVersion)
	}

	got, err := model2.GetByEmail(originalEmail)
	if err != nil {
		t.Fatalf("get user after failed update: %v", err)
	}

	if got.Email != originalEmail {
		t.Errorf("stored Email = %q, want %q", got.Email, originalEmail)
	}

	if got.Version != beforeVersion {
		t.Errorf("stored Version = %d, want %d", got.Version, beforeVersion)
	}

	if !got.UpdatedAt.Equal(beforeUpdate) {
		t.Errorf("stored UpdatedAt = %v, want %v", got.UpdatedAt, beforeUpdate)
	}
}
