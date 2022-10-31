package todo

import "time"

type TodoItem struct {
	Id          int
	Action      string
	DueDate     *time.Time `db:"due_date"`
	Duration    *time.Duration
	Completed   bool
	CompletedAt *time.Time `db:"completed_at"`
}

func (t *TodoItem) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}

	return t.DueDate.Before(time.Now())
}

func (t *TodoItem) Complete() {
	t.Completed = true

	now := time.Now()
	t.CompletedAt = &now
}