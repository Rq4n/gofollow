package service

import (
	"context"
	"time"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/google/uuid"
)

type (
	// ClientServiceInterface defines client-related methods
	ClientServiceInterface interface {
		CreateNewClient(ctx context.Context, userID uuid.UUID, name, email, invoiceLink string, sendDate time.Time) error
		GetClientByID(ctx context.Context, id uuid.UUID) (*repository.Client, error)
	}

	// EmailJobServiceInterface defines email job-related methods
	EmailJobServiceInterface interface {
		GetEmailByJobID(ctx context.Context, id uuid.UUID) (*repository.EmailJob, error)
		GetPendingEmailJobs(ctx context.Context) ([]repository.EmailJob, error)
		MarkJobAsCompleted(ctx context.Context, id uuid.UUID) error
		MarkJobAsFailed(ctx context.Context, id uuid.UUID) error
		TryMarkJobAsProcessing(ctx context.Context, id uuid.UUID) (int64, error)
	}

	// UserServiceInterface defines user-related methods
	UserServiceInterface interface {
		CreateNewUser(ctx context.Context, email, password string) error
		GetUserByName(ctx context.Context, name, password string) (*repository.GetUserByNameRow, error)
	}
)

type AppService struct {
	UserService     *UserService
	ClientService   *ClientService
	EmailJobService *EmailJobService
}

func NewAppService(repo repository.Querier) *AppService {
	return &AppService{
		UserService:     NewUserService(repo),
		ClientService:   NewClientService(repo),
		EmailJobService: NewEmailJobService(repo),
	}
}
