package service

import (
	"context"

	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/google/uuid"
)

type mockRepo struct {
	repository.Querier
	createUserErr     error
	createClientErr   error
	createEmailJobErr error

	getUserByEmailFn  func(email string) (repository.GetUserByEmailRow, error)
	createNewClientFn func(arg repository.CreateNewClientParams) (uuid.UUID, error)
	getClientByUUIDFn func(id uuid.UUID) (repository.Client, error)
}

func (m *mockRepo) CreateNewUser(ctx context.Context, arg repository.CreateNewUserParams) error {
	return m.createUserErr
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (repository.GetUserByEmailRow, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(email)
	}
	return repository.GetUserByEmailRow{}, nil
}

func (m *mockRepo) CreateNewClient(ctx context.Context, arg repository.CreateNewClientParams) (uuid.UUID, error) {
	if m.createNewClientFn != nil {
		return m.createNewClientFn(arg)
	}
	return uuid.New(), m.createClientErr
}

func (m *mockRepo) GetClientByUUID(ctx context.Context, id uuid.UUID) (repository.Client, error) {
	if m.getClientByUUIDFn != nil {
		return m.getClientByUUIDFn(id)
	}
	return repository.Client{}, nil
}

func (m *mockRepo) CreateEmailJob(ctx context.Context, arg repository.CreateEmailJobParams) error {
	return m.createEmailJobErr
}

func (m *mockRepo) GetEmailJobByID(ctx context.Context, id uuid.UUID) (repository.EmailJob, error) {
	return repository.EmailJob{}, nil
}

func (m *mockRepo) GetPendingEmailJobs(ctx context.Context) ([]repository.EmailJob, error) {
	return nil, nil
}
func (m *mockRepo) MarkJobAsCompleted(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepo) MarkJobAsFailed(ctx context.Context, id uuid.UUID) error    { return nil }
func (m *mockRepo) TryMarkJobAsProcessing(ctx context.Context, id uuid.UUID) (int64, error) {
	return 0, nil
}
