package cmd

import (
	"strings"
	"time"

	"github.com/hako/durafmt"
	"github.com/spf13/cobra"

	"github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/AP-Hunt/what-next/m/views"
)

var TodoRootCmd = &cobra.Command{
	Use:     "todo",
	Aliases: []string{"t"},
}

var TodoAddCmd = &cobra.Command{
	Use:                   "add action [--due due] [--duration duration]",
	DisableFlagsInUseLine: true,
	Aliases:               []string{"a"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)
		var dueDate *time.Time = nil

		if cmd.Flags().Lookup("due") != nil {
			dueDateInput, err := cmd.Flags().GetString("due")
			if err != nil {
				return err
			}

			if dueDateInput != "" {
				parsedDate, err := todo.ParseDueDate(dueDateInput)
				if err != nil {
					return err
				}
				dueDate = &parsedDate
			}
		}

		var duration *time.Duration = nil

		if cmd.Flags().Lookup("duration") != nil {
			durationInput, err := cmd.Flags().GetString("duration")
			if err != nil {
				return err
			}
			if durationInput != "" {
				parsedDuration, err := durafmt.ParseString(durationInput)
				if err != nil {
					return err
				}
				pd := parsedDuration.Duration()
				duration = &pd
			}
		}

		itemAction := strings.Join(args, " ")
		item := todo.TodoItem{
			Action:    itemAction,
			Completed: false,
			DueDate:   dueDate,
			Duration:  duration,
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
	TodoAddCmd.Flags().String("due", "", todoAddDueDateHelp)
	TodoAddCmd.Flags().String("duration", "", todoAddDurationHelp)
	TodoRootCmd.AddCommand(TodoAddCmd)
	TodoRootCmd.AddCommand(TodoListCmd)
}

var todoAddDueDateHelp = `Optional. Date and time at which the new item is due.

Due dates can be specified as any valid datetime string, and are assumed to be local time.
For due dates of today or tomorrow, the following shorthand strings can be used:
* @today
* @tod
* @tomorrow
* @tom
* @tmrw
`

var todoAddDurationHelp = `Optional. Duration you expect this item to take. Durations can be provided in a human readable form, e.g. '30m' or '1h10m'.`
