-- +goose Up
-- +goose StatementBegin

DROP TABLE IF EXISTS places;

CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price TEXT DEFAULT 'mid',
    time FLOAT,
    types_of_movement TEXT,
    category TEXT,
    lat_start FLOAT,
    lon_start FLOAT,
    lat_end FLOAT,
    lon_end FLOAT,
    is_indoor BOOLEAN DEFAULT FALSE,
    with_child BOOLEAN DEFAULT FALSE,
    with_pets BOOLEAN DEFAULT FALSE,
    description TEXT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS places;

-- +goose StatementEnd
