CREATE TABLE IF NOT EXISTS clients (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    email        TEXT NOT NULL,
    invoice_link TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW()
);
