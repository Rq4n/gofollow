CREATE TABLE IF NOT EXISTS email_jobs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id    UUID REFERENCES clients(id) ON DELETE CASCADE,
    status       TEXT NOT NULL CHECK(status IN('pending', 'processing', 'failed', 'completed')),

    send_at      TIMESTAMPTZ NOT NULL,
    sent_at      TIMESTAMPTZ,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
