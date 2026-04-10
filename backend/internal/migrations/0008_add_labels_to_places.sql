-- +goose Up
-- +goose StatementBegin

ALTER TABLE places ADD COLUMN IF NOT EXISTS name_label TEXT;
ALTER TABLE places ADD COLUMN IF NOT EXISTS description_label TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE places DROP COLUMN IF EXISTS name_label;
ALTER TABLE places DROP COLUMN IF EXISTS description_label;

-- +goose StatementEnd
