package todo

import "time"

type TodoItem struct {
	Id        int
	Action    string
	DueDate   time.Time `db:"due_date"`
	Completed bool
}
