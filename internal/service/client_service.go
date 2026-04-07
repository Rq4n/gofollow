package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

	clientID, err := s.repo.CreateNewClient(ctx, arg)
	if err != nil {
		log.Print("failed to create client")
		return err
	}
	jobArg := repository.CreateEmailJobParams{
		ClientID: clientID,
		SendAt: pgtype.Timestamptz{
			Time:  sendDate,
			Valid: true,
		},
	}

	if err := s.repo.CreateEmailJob(ctx, jobArg); err != nil {
		return fmt.Errorf("Failed to create email job: %w", err)
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
