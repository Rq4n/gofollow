package service

import (
	"context"
	"log"

	"github.com/Rq4n/gofollow/internal/repository"
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
		log.Printf("failed to create user %v", err)
		return err
	}
	return nil
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *UserService) GetUserByName(ctx context.Context, name, password string) (*repository.GetUserByNameRow, error) {
	user, err := s.repo.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if err := CheckPassword(password, user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}
