-- +goose Up
CREATE TABLE calendars (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    display_name TEXT NOT NULL UNIQUE,
    calendar_url TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE calendars;