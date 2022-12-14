package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	"github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/views"
	ical "github.com/arran4/golang-ical"
	"github.com/spf13/cobra"
)

var CalendarRootCmd = &cobra.Command{
	Use:     "calendar",
	Aliases: []string{"cal", "c"},
}

var CalendarViewCmd = &cobra.Command{
	Use:     "view display_name",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"v"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)

		viewEngine := ctx.ViewEngine()
		calService := ctx.CalendarService()

		displayName := args[0]
		calRecord, err := calService.GetCalendarByDisplayName(displayName)
		if err != nil {
			if _, ok := err.(calendar.ErrNotFound); ok {
				fmt.Printf("cannot find calendar '%s'\n", displayName)
				return nil
			}

			fmt.Printf("error finding calendar: %s\n", err)
			return nil
		}

		cal, err := calService.OpenCalendar(calRecord.URL)
		if err != nil {
			fmt.Printf("error opening calendar: %s", err)
		}

		startOfTheDay := time.Now().Truncate(24 * time.Hour)

		todayOnlyCal := ical.NewCalendar()
		for _, evt := range cal.Events() {
			startsToday, err := calendar.EventStartsToday(evt)
			if err != nil {
				continue
			}
			endsToday, err := calendar.EventEndsToday(evt)
			if err != nil {
				continue
			}

			if startsToday || endsToday {
				todayOnlyCal.AddVEvent(evt)
			}
		}

		calendarView := &views.CalendarView{}
		calendarViwData := &views.CalendarViewData{
			Calendar:   todayOnlyCal,
			TargetDate: startOfTheDay,
		}
		calendarView.SetData(calendarViwData)

		return viewEngine.Draw(calendarView)
	},
}

var CalendarAddCmd = &cobra.Command{
	Use:     "add display_name url",
	Aliases: []string{"a"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}

		calendarUrl := args[1]
		_, err := url.Parse(calendarUrl)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)
		calService := ctx.CalendarService()

		displayName := args[0]
		url := args[1]

		_, err := calService.AddCalendar(url, displayName)
		if err != nil {
			if _, ok := err.(*calendar.ErrDuplicateCalendarDisplayName); ok {
				fmt.Printf("Calendar with display name '%s' already exists\n", displayName)
				return nil
			}

			return err
		}

		fmt.Printf("Added calendar %s (%s)\n", displayName, url)
		return nil
	},
}

var CalendarRemoveCmd = &cobra.Command{
	Use:     "remove display_name",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)
		calService := ctx.CalendarService()

		displayName := args[0]

		cal, err := calService.GetCalendarByDisplayName(displayName)
		if err != nil {
			if _, ok := err.(*calendar.ErrNotFound); ok {
				fmt.Printf("Cannot find calendar with display name '%s'.\n", displayName)
				return nil
			}

			return err
		}

		err = calService.RemoveById(cal.Id)
		if err != nil {
			return err
		}

		fmt.Printf("Calendar '%s' removed.\n", displayName)
		return nil
	},
}

var CalendarListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)

		viewEngine := ctx.ViewEngine()
		calService := ctx.CalendarService()

		calendars, err := calService.GetAllCalendars()
		if err != nil {
			return err
		}

		calListView := views.CalendarListView{}
		calListView.SetData(calendars)
		return viewEngine.Draw(&calListView)
	},
}

func init() {
	CalendarRootCmd.AddCommand(CalendarViewCmd)
	CalendarRootCmd.AddCommand(CalendarAddCmd)
	CalendarRootCmd.AddCommand(CalendarRemoveCmd)
	CalendarRootCmd.AddCommand(CalendarListCmd)
}
