-- +goose Up
ALTER TABLE todo_items
    ADD COLUMN completed_at DATETIME NULL;

UPDATE todo_items
SET completed_at = datetime()
WHERE completed = 1;


-- +goose Down
ALTER TABLE todo_items
    DROP COLUMN completed_at;