package todo

import (
	"context"
	"errors"

	"github.com/AP-Hunt/what-next/m/db"
	"github.com/jmoiron/sqlx"
)

var (
	ItemNotFoundError = errors.New("not found")
)

//counterfeiter:generate -o fakes/ . TodoRepositoryInterface
type TodoRepositoryInterface interface {
	Add(item TodoItem) (TodoItem, error)
	Get(id int) (TodoItem, error)
	List() (*TodoItemCollection, error)
	Update(item TodoItem) (TodoItem, error)
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
	val, err := db.InTransaction(
		func(tx *sqlx.Tx) (*TodoItem, error) {
			var duration *int = nil
			if item.Duration != nil {
				d := int(*item.Duration)
				duration = &d
			}
			row := tx.QueryRowx(
				`
				INSERT INTO todo_items
					(action, due_date, duration, completed, completed_at)
				VALUES
					(?, ?, ?, ?, ?)
		
				RETURNING *
				`,
				item.Action,
				item.DueDate,
				duration,
				item.Completed,
				item.CompletedAt,
			)

			newItem := TodoItem{}
			err := row.StructScan(&newItem)
			return &newItem, err
		},
		repo.conn,
		repo.ctx,
	)

	if err != nil {
		return TodoItem{}, err
	}

	return *val, nil
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

func (repo *TodoSQLRepository) Update(item TodoItem) (TodoItem, error) {
	updated, err := db.InTransaction(
		func(tx *sqlx.Tx) (*TodoItem, error) {

			var updatedItem TodoItem

			err := tx.QueryRowx(
				`
					UPDATE todo_items
					SET
						action = ?,
						due_date = ?,
						duration = ?,
						completed = ?,
                         completed_at = ?
					WHERE 
						id = ?

					RETURNING *
				`,
				item.Action,
				item.DueDate,
				item.Duration,
				item.Completed,
				item.CompletedAt,
				item.Id,
			).StructScan(&updatedItem)

			return &updatedItem, err
		},
		repo.conn,
		repo.ctx,
	)

	if err != nil {
		return TodoItem{}, err
	}

	return *updated, nil
}
