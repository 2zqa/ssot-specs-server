CREATE TABLE IF NOT EXISTS devices (
    uuid uuid PRIMARY KEY,
    last_updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);