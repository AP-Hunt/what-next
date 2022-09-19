package cmd_test

import (
	"context"
	"time"

	"github.com/AP-Hunt/what-next/m/cmd"
	commandContext "github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/todo"
	. "github.com/AP-Hunt/what-next/m/todo/fakes"
	"github.com/AP-Hunt/what-next/m/views"
	. "github.com/AP-Hunt/what-next/m/views/fakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Todo", func() {
	var (
		viewEngine *FakeViewEngineInterface
		todoRepo   *FakeTodoRepositoryInterface
		cmdContext commandContext.CommandContext
	)

	BeforeEach(func() {
		viewEngine = &FakeViewEngineInterface{}
		todoRepo = &FakeTodoRepositoryInterface{}

		cmdContext = commandContext.NewCommandContext(context.Background()).
			WithTodoRepository(todoRepo).
			WithViewEngine(viewEngine)

	})

	Describe("Add", func() {
		It("adds a new item and renders the TodoList view", func() {
			now := time.Now()
			todoRepo.AddReturns(
				todo.TodoItem{
					Id:        1,
					Action:    "foo",
					DueDate:   &now,
					Completed: false,
				},
				nil,
			)

			PrepareCommandForTest(cmd.TodoAddCmd, []string{"foo", "bar"})

			err := cmd.TodoAddCmd.ExecuteContext(cmdContext)
			Expect(err).ToNot(HaveOccurred())

			Expect(viewEngine.DrawCallCount()).To(Equal(1))
			drawnView := viewEngine.DrawArgsForCall(0)

			Expect(drawnView).To(BeAssignableToTypeOf(&views.TodoListView{}))

		})
	})

	Describe("List", func() {
		It("gets existing items and renders a TodoList view", func() {
			now := time.Now()
			todoRepo.ListReturns(
				todo.NewTodoItemCollection([]*todo.TodoItem{
					{
						Id:        1,
						Action:    "foo",
						DueDate:   &now,
						Completed: false,
					},
				}),
				nil,
			)

			PrepareCommandForTest(cmd.TodoListCmd, []string{})

			err := cmd.TodoListCmd.ExecuteContext(cmdContext)
			Expect(err).ToNot(HaveOccurred())

			Expect(viewEngine.DrawCallCount()).To(Equal(1))
			drawnView := viewEngine.DrawArgsForCall(0)

			Expect(drawnView).To(BeAssignableToTypeOf(&views.TodoListView{}))

		})
	})
})
