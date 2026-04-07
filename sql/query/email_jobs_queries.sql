-- name: CreateEmailJob :exec
INSERT INTO email_jobs (client_id, send_at)
VALUES ($1, $2);

-- name: GetEmailJobByID :one
SELECT id, client_id, status, send_at, sent_at, created_at
FROM email_jobs
WHERE id = $1;

-- name: GetPendingEmailJobs :many
SELECT id, client_id, status, send_at, sent_at, created_at
FROM email_jobs
WHERE status = 'pending'
AND send_at <= NOW();

-- name: TryMarkJobAsProcessing :execrows
UPDATE email_jobs
SET status = 'processing'
WHERE id = $1
AND status = 'pending';

-- name: MarkJobAsCompleted :exec
UPDATE email_jobs
SET status = 'completed', sent_at = NOW()
WHERE id = $1 
AND status = 'processing';

-- name: MarkJobAsFailed :exec
UPDATE email_jobs
SET status = 'failed'
WHERE id = $1
AND (status = 'pending' OR status = 'processing');

