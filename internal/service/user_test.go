package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateNewUser_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := NewUserService(repo)

	err := svc.CreateNewUser(context.Background(), "foo@gmail.com", "Ry@nLima06##")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCreateNewUser_AlreadyExists(t *testing.T) {
	repo := &mockRepo{
		createUserErr: &pgconn.PgError{Code: "23505"},
	}
	svc := NewUserService(repo)

	err := svc.CreateNewUser(context.Background(), "foo@gmail.com", "Ry@nLima06##")
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestGetUserByEmail_Succes(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("Ry@nLima06##"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password %v", err)
	}
	repo := &mockRepo{
		getUserByEmailFn: func(email string) (repository.GetUserByEmailRow, error) {
			return repository.GetUserByEmailRow{
				ID:       uuid.New(),
				Email:    email,
				Password: string(hash),
			}, nil
		},
	}
	svc := NewUserService(repo)

	user, err := svc.GetUserByEmail(context.Background(), "foo@gmail.com", "Ry@nLima06##")
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	if user.Email != "foo@gmail.com" {
		t.Errorf("expected nil got %v", err)
	}
}

func TestGetUserByEmail_WrongPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("Foobar"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password %v", err)
	}
	repo := &mockRepo{
		getUserByEmailFn: func(email string) (repository.GetUserByEmailRow, error) {
			return repository.GetUserByEmailRow{
				ID:       uuid.New(),
				Email:    email,
				Password: string(hash),
			}, nil
		},
	}
	svc := NewUserService(repo)

	user, err := svc.GetUserByEmail(context.Background(), "foo@gmail.com", "SenhaErrada")
	if err == nil {
		t.Errorf("expected error got nil")
	}
	if user != nil {
		t.Errorf("expected nil got %v", err)
	}
}

// TODO: test GetUserByEmail_NotFound
