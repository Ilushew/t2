-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS place_comments (
    id SERIAL PRIMARY KEY,
    place_id INT NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    author TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_place_comments_place_id ON place_comments(place_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS place_comments;

-- +goose StatementEnd
