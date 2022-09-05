package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/AP-Hunt/what-next/m/views"
)

var todoCompletedSymbolMap = map[bool]string{
	true:  "✓",
	false: "x",
}

var TodoRootCmd = &cobra.Command{
	Use:     "todo",
	Aliases: []string{"t"},
}

var TodoAddCmd = &cobra.Command{
	Use:     "add item",
	Aliases: []string{"a"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)

		itemAction := strings.Join(args, " ")
		item := todo.TodoItem{
			Action:    itemAction,
			Completed: false,
			DueDate:   time.Now(),
		}

		repo := ctx.TodoRepository()

		addedItem, err := repo.Add(item)
		if err != nil {
			return err
		}

		viewEngine := ctx.ViewEngine()
		view := views.TodoListView{}
		view.SetData(todo.NewTodoItemCollection([]*todo.TodoItem{&addedItem}))

		return viewEngine.Draw(&view)
	},
}

var TodoListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)
		repo := ctx.TodoRepository()

		items, err := repo.List()
		if err != nil {
			return err
		}

		viewEngine := ctx.ViewEngine()
		view := views.TodoListView{}
		view.SetData(items)

		return viewEngine.Draw(&view)
	},
}

func init() {
	TodoRootCmd.AddCommand(TodoAddCmd)
	TodoRootCmd.AddCommand(TodoListCmd)
}
