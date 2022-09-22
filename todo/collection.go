package todo

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
