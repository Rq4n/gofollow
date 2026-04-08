package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.Querier
}

func NewUserService(repo repository.Querier) *UserService {
	return &UserService{
		repo: repo,
	}
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

func (s *UserService) CreateNewUser(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	arg := repository.CreateNewUserParams{
		Email:    email,
		Password: string(hash),
	}

	if err := s.repo.CreateNewUser(ctx, arg); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email, password string) (*repository.GetUserByEmailRow, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	return &user, nil
}
