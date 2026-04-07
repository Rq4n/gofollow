package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/Rq4n/gofollow/internal/mailer"
	"github.com/Rq4n/gofollow/internal/service"
	"github.com/google/uuid"
)

type Job struct {
	ID uuid.UUID
}

type Worker struct {
	Email  *service.EmailJobService
	Client *service.ClientService
	Mail   mailer.Mailer
}

func (w *Worker) ProcessSingleJob(ctx context.Context, JobID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			log.Printf("job %s failed: %v", JobID, err)
			_ = w.Email.MarkJobAsFailed(ctx, JobID)
		}
	}()

	// this is the equivalent of a lock, so only one worker
	// can get a specific job
	rows, err := w.Email.TryMarkJobAsProcessing(ctx, JobID)
	if err != nil {
		return fmt.Errorf("Lock error: %w", err)
	}
	if rows == 0 {
		return nil
	}

	// get email_job
	job, err := w.Email.GetEmailByJobID(ctx, JobID)
	if err != nil {
		return fmt.Errorf("Failed to get job error: %w", err)
	}

	// get client
	client, err := w.Client.GetClientByID(ctx, job.ClientID)
	if err != nil {
		return fmt.Errorf("Failed to get client error: %w", err)
	}

	// send email with client data
	err = w.Mail.Send(
		client.Name,
		client.Email,
		client.InvoiceLink,
	)
	if err != nil {
		return fmt.Errorf("Failed to send email error: %w", err)
	}

	// mark job as completed
	err = w.Email.MarkJobAsCompleted(ctx, JobID)
	if err != nil {
		return fmt.Errorf("Failed to mark job completed error: %w", err)
	}

	return nil
}
