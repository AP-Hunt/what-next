package cmd

import (
	"time"

	"github.com/AP-Hunt/what-next/m/context"
	. "github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/scheduler"
	"github.com/AP-Hunt/what-next/m/views"
	ical "github.com/arran4/golang-ical"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "what-next",
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)
		viewEngine := ctx.ViewEngine()
		calService := ctx.CalendarService()
		repo := ctx.TodoRepository()

		allCalendarRecords, err := calService.GetAllCalendars()
		if err != nil {
			return err
		}

		calendars := []*ical.Calendar{}
		for _, record := range allCalendarRecords {
			cal, err := calService.OpenCalendar(record.URL)
			if err != nil {
				return err
			}

			calendars = append(calendars, cal)
		}

		todoList, err := repo.List()
		if err != nil {
			return err
		}

		schedule, err := scheduler.GenerateSchedule(time.Now(), calendars, todoList)

		scheduleView := views.ScheduleView{}
		scheduleView.SetData(schedule)

		return viewEngine.Draw(&scheduleView)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(TodoRootCmd)
	RootCmd.AddCommand(CalendarRootCmd)
}

func ExecuteCWithArgs(ctx CommandContext, args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.ExecuteContext(ctx)
}

func ExecuteC(ctx CommandContext) error {
	return RootCmd.ExecuteContext(ctx)
}
