package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

var ErrClientNotFound = errors.New("client not found")

func (s *ClientService) CreateNewClient(ctx context.Context, userID uuid.UUID, name, email, invoiceLink string, sendDate time.Time) error {
	arg := repository.CreateNewClientParams{
		UserID:      userID,
		Name:        name,
		Email:       email,
		InvoiceLink: invoiceLink,
	}

	clientID, err := s.repo.CreateNewClient(ctx, arg)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return err
	}
	jobArg := repository.CreateEmailJobParams{
		Status:   "pending",
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

func (s *ClientService) GetClientByUUID(ctx context.Context, id uuid.UUID) (*repository.Client, error) {
	resp, err := s.repo.GetClientByUUID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClientNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &resp, nil
}

func (s *ClientService) GetAllClients(ctx context.Context, userID uuid.UUID) ([]repository.Client, error) {
	resp, err := s.repo.GetAllClients(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("db error %v", err)
	}
	return resp, nil
}
