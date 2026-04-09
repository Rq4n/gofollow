// internal/service/user_test.go
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateNewUser_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := NewUserService(repo)

	err := svc.CreateNewUser(context.Background(), "test@email.com", "password123")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCreateNewUser_AlreadyExists(t *testing.T) {
	repo := &mockRepo{
		createUserErr: &pgconn.PgError{Code: "23505"},
	}
	svc := NewUserService(repo)

	err := svc.CreateNewUser(context.Background(), "test@email.com", "password123")
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

// TODO
// test GetUserByEmail_Success
// test GetUserByEmail_WrongPassword
// test GetUserByEmail_NotFound
