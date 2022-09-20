-- +goose Up
ALTER TABLE todo_items 
    ADD COLUMN duration INT NULL;

-- +goose Down
ALTER TABLE todo_items 
    DROP COLUMN duration;