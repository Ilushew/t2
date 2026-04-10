-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS place_images (
    id SERIAL PRIMARY KEY,
    place_id INT NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    sort_order INT DEFAULT 0
);

CREATE INDEX idx_place_images_place_id ON place_images(place_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS place_images;

-- +goose StatementEnd
