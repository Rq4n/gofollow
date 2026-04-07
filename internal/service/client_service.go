package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/google/uuid"
)

type ClientService struct {
	repo repository.Querier
}

func NewClientService(repo repository.Querier) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

func (s *ClientService) CreateNewClient(ctx context.Context, userID uuid.UUID, name, email, invoiceLink string, sendDate time.Time) error {
	arg := repository.CreateNewClientParams{
		UserID:      userID,
		Name:        name,
		Email:       email,
		InvoiceLink: invoiceLink,
	}

	if err := s.repo.CreateNewClient(ctx, arg); err != nil {
		log.Print("failed to create client")
		return err
	}

	return nil
}

func (s *ClientService) GetClientByID(ctx context.Context, id uuid.UUID) (*repository.Client, error) {
	resp, err := s.repo.GetClientByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service_error: %w", err)
	}
	return &resp, nil
}
