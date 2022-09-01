-- +goose Up
CREATE TABLE todo_items (
    id INTEGER PRIMARY KEY,
    action TEXT NOT NULL,
    due_date DATETIME,
    completed BOOLEAN DEFAULT 0
);

-- +goose Down
DROP TABLE todo_items;