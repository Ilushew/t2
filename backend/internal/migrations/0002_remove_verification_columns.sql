-- +goose Up
-- +goose StatementBegin

ALTER TABLE users
DROP COLUMN IF EXISTS email_verification_code,
DROP COLUMN IF EXISTS email_verification_expires_at,
DROP COLUMN IF EXISTS password_hash,
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users
ADD COLUMN email_verification_code VARCHAR(6),
ADD COLUMN email_verification_expires_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN password_hash VARCHAR(255),
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- +goose StatementEnd
