CREATE TABLE IF NOT EXISTS public.urls (
  slug VARCHAR(7) PRIMARY KEY,
  original TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  clicked INT DEFAULT 0
);