-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sensors (
    id TEXT PRIMARY KEY,
    dry_reference INTEGER,
    wet_reference INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS plants (
    id SERIAL PRIMARY KEY,
    sensor_id TEXT UNIQUE REFERENCES sensors(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    date_planted TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS readings (
    id SERIAL PRIMARY KEY,
    plant_id INTEGER REFERENCES plants(id) ON DELETE CASCADE,
    moisture_percentage INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS readings;
DROP TABLE IF EXISTS plants;
DROP TABLE IF EXISTS sensors;
-- +goose StatementEnd
