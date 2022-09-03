package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/todo"
)

var todoCompletedSymbolMap = map[bool]string{
	true:  "✓",
	false: "x",
}

var todoCmd = &cobra.Command{
	Use:     "todo",
	Aliases: []string{"t"},
}

var todoAddCmd = &cobra.Command{
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

		cmd.Printf("%-12d [%s] %s\n", addedItem.Id, todoCompletedSymbolMap[addedItem.Completed], addedItem.Action)
		return nil
	},
}

var todoListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)
		repo := ctx.TodoRepository()

		items, err := repo.List()
		if err != nil {
			return err
		}

		for _, item := range items.Enumerate() {
			cmd.Printf("%-12d [%s] %s\n", item.Id, todoCompletedSymbolMap[item.Completed], item.Action)
		}

		return nil
	},
}

func init() {
	todoCmd.AddCommand(todoAddCmd)
	todoCmd.AddCommand(todoListCmd)
}
