package todo_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/AP-Hunt/what-next/m/todo"
)

var _ = Describe("Model", func() {
	Describe("TodoItem", func() {
		Describe("Complete", func() {
			It("sets the completed at date at the same time as the completed flag", func(){
				item := todo.TodoItem{}

				Expect(item.Completed).To(BeFalse())
				Expect(item.CompletedAt).To(BeNil())

				item.Complete()

				Expect(item.Completed).To(BeTrue())
				Expect(item.CompletedAt).ToNot(BeNil())
			})
		})
	})
})
