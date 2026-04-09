package service

import (
	"context"
	"fmt"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/google/uuid"
)

type EmailJobService struct {
	repo repository.Querier
}

func NewEmailJobService(repo repository.Querier) *EmailJobService {
	return &EmailJobService{
		repo: repo,
	}
}

func (s *EmailJobService) CreateEmailJob(ctx context.Context, arg repository.CreateEmailJobParams) error {
	if err := s.repo.CreateEmailJob(ctx, arg); err != nil {
		return fmt.Errorf("email_job error: %w", err)
	}
	return nil
}

func (s *EmailJobService) GetEmailByJobID(ctx context.Context, id uuid.UUID) (*repository.EmailJob, error) {
	job, err := s.repo.GetEmailJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("email_job error: %w", err)
	}

	return &job, nil
}

func (s *EmailJobService) GetPendingEmailJobs(ctx context.Context) ([]repository.EmailJob, error) {
	job, err := s.repo.GetPendingEmailJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("email_job error: %w", err)
	}

	return job, nil
}

func (s *EmailJobService) MarkJobAsCompleted(ctx context.Context, id uuid.UUID) error {
	err := s.repo.MarkJobAsCompleted(ctx, id)
	if err != nil {
		return fmt.Errorf("email_job error: %w", err)
	}
	return nil
}

func (s *EmailJobService) MarkJobAsFailed(ctx context.Context, id uuid.UUID) error {
	err := s.repo.MarkJobAsFailed(ctx, id)
	if err != nil {
		return fmt.Errorf("email_job error: %w", err)
	}
	return nil
}

func (s *EmailJobService) TryMarkJobAsProcessing(ctx context.Context, id uuid.UUID) (int64, error) {
	rows, err := s.repo.TryMarkJobAsProcessing(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("email_job error: %w", err)
	}
	return rows, nil
}
