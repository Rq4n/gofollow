CREATE TABLE IF NOT EXISTS users(
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email       TEXT NOT NULL,
  password    TEXT NOT NULL,
  created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
