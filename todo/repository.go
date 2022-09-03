package todo

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ItemNotFoundError = errors.New("not found")
)

type TodoRepositoryInterface interface {
	Add(item TodoItem) (TodoItem, error)
	Get(id int) (TodoItem, error)
	List() (*TodoItemCollection, error)
}

type TodoSQLRepository struct {
	conn *sqlx.DB
	ctx  context.Context
}

func NewTodoSQLRepository(conn *sqlx.DB, ctx context.Context) *TodoSQLRepository {
	return &TodoSQLRepository{
		conn: conn,
		ctx:  ctx,
	}
}

func (repo *TodoSQLRepository) Add(item TodoItem) (TodoItem, error) {
	tx, err := repo.conn.BeginTxx(repo.ctx, nil)
	if err != nil {
		return TodoItem{}, err
	}

	row := tx.QueryRowx(
		`
		INSERT INTO todo_items
			(action, due_date, completed)
		VALUES
			(?, ?, ?)

		RETURNING *
		`,
		item.Action,
		item.DueDate,
		item.Completed,
	)

	newItem := TodoItem{}
	err = row.StructScan(&newItem)

	if err != nil {
		if e := tx.Rollback(); e != nil {
			return TodoItem{}, e
		}

		return TodoItem{}, err
	}

	err = tx.Commit()
	if err != nil {
		return TodoItem{}, err
	}

	return newItem, nil
}

func (repo *TodoSQLRepository) Get(id int) (TodoItem, error) {

	item := TodoItem{}
	err := repo.conn.GetContext(repo.ctx, &item, "SELECT * FROM todo_items WHERE id = ?", id)
	if err != nil {
		return TodoItem{}, err
	}

	return item, nil
}

func (repo *TodoSQLRepository) List() (*TodoItemCollection, error) {
	items := []*TodoItem{}
	err := repo.conn.Select(&items, "SELECT * FROM todo_items")
	if err != nil {
		return NewTodoItemCollection(nil), err
	}

	return NewTodoItemCollection(items), nil
}
