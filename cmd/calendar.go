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
	"golang.org/x/exp/slices"
)

var CalendarRootCmd = &cobra.Command{
	Use:     "calendar",
	Aliases: []string{"cal", "c"},
}

var CalendarViewCmd = &cobra.Command{
	Use:     "view display_name",
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
		parsedUrl, err := url.Parse(calendarUrl)
		if err != nil {
			return err
		}

		if !slices.Contains([]string{"http", "https", "file"}, parsedUrl.Scheme) {
			return fmt.Errorf("unsupported scheme '%s'", parsedUrl.Scheme)
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

		fmt.Printf("Added calendar %s (%s)", displayName, url)
		return nil
	},
}

func init() {
	CalendarRootCmd.AddCommand(CalendarViewCmd)
	CalendarRootCmd.AddCommand(CalendarAddCmd)
}
