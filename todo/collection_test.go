package todo_test

import (
	"time"

	. "github.com/AP-Hunt/what-next/m/todo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TodoItemCollection", func() {
	Describe("Enumerate", func() {
		It("exposes the internally held items for iteration", func() {
			now := time.Now()
			items := []*TodoItem{
				{
					Id:        1,
					Action:    "First",
					DueDate:   &now,
					Completed: false,
				},
				{
					Id:        2,
					Action:    "Second",
					DueDate:   &now,
					Completed: false,
				},
			}

			collection := NewTodoItemCollection(items)

			Expect(collection.Enumerate()).To(BeEquivalentTo(items))
		})
	})

	Describe("Filter", func() {
		It("returns a new collection containing the items which met the filter function criteria", func() {
			now := time.Now()
			items := []*TodoItem{
				{
					Id:        1,
					Action:    "First",
					DueDate:   &now,
					Completed: false,
				},
				{
					Id:        2,
					Action:    "Second",
					DueDate:   &now,
					Completed: true,
				},
			}

			collection := NewTodoItemCollection(items)

			filteredCollecton := collection.Filter(func(t *TodoItem) bool {
				return t.Completed
			})

			Expect(filteredCollecton).ToNot(BeIdenticalTo(collection))
			Expect(filteredCollecton.Enumerate()[0].Id).To(Equal(2))
		})
	})
})
