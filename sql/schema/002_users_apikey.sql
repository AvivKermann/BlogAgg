-- +goose Up
ALTER TABLE users
ADD COLUMN api_key VARCHAR UNIQUE DEFAULT encode(sha256(random()::text::bytea), 'hex') NOT NULL;

-- +goose Down
ALTER TABLE users
REMOVE COLUMN api_key;