package todo

import "golang.org/x/exp/slices"

type TodoItemCollection struct {
	items []*TodoItem
}

func NewTodoItemCollection(items []*TodoItem) *TodoItemCollection {
	return &TodoItemCollection{
		items: items,
	}
}

func (c *TodoItemCollection) Enumerate() []*TodoItem {
	return c.items
}

func (c *TodoItemCollection) Filter(filter func(*TodoItem) bool) *TodoItemCollection {
	filteredItems := []*TodoItem{}

	for _, item := range c.items {
		if filter(item) {
			filteredItems = append(filteredItems, item)
		}
	}

	return &TodoItemCollection{
		items: filteredItems,
	}
}

func (c *TodoItemCollection) Len() int {
	return len(c.items)
}

func (c *TodoItemCollection) Append(other *TodoItemCollection) *TodoItemCollection {
	aItems := c.items
	bItems := other.items

	for _, i := range bItems {
		aItems = append(aItems, i)
	}

	return NewTodoItemCollection(aItems)
}

func (c *TodoItemCollection) SortByDueDateAsc() *TodoItemCollection {
	items := c.items

	slices.SortFunc(items, func(a *TodoItem, b *TodoItem) bool {
		// Both tasks don't have a due date
		if a.DueDate == nil && b.DueDate == nil {
			// Sort by ID on the assumption that smaller IDs are older tasks
			return a.Id < b.Id
		}

		// Tasks with a due date go higher
		if a.DueDate != nil && b.DueDate == nil {
			return true
		}

		if a.DueDate == nil && b.DueDate != nil {
			return false
		}

		// Both due dates are populated
		return a.DueDate.Before(*b.DueDate) || a.DueDate.Equal(*b.DueDate)
	})

	return NewTodoItemCollection(items)
}
