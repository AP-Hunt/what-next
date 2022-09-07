package cmd

import (
	"time"

	"github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/views"
	"github.com/spf13/cobra"

	ical "github.com/arran4/golang-ical"
)

var CalendarRootCmd = &cobra.Command{
	Use:     "calendar",
	Aliases: []string{"cal", "c"},
}

var CalendarViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"v"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx context.CommandContext = cmd.Context().(context.CommandContext)

		viewEngine := ctx.ViewEngine()

		calendarView := &views.CalendarView{}

		midnight := time.Now().Truncate(24 * time.Hour)

		cal := ical.NewCalendar()
		evt1 := cal.AddEvent("evt-1")
		evt1.SetProperty(ical.ComponentProperty(ical.PropertyName), "Event 1")
		evt1.SetStartAt(midnight.Add(10 * time.Hour))
		evt1.SetDuration(30 * time.Minute)
		evt1.SetLocation("outside")

		evt2 := cal.AddEvent("evt-2")
		evt2.SetProperty(ical.ComponentProperty(ical.PropertyName), "Event 2")
		evt2.SetStartAt(midnight.Add(12 * time.Hour))
		evt2.SetDuration(90 * time.Minute)
		evt2.SetLocation("inside")

		calendarView.SetData(cal)

		return viewEngine.Draw(calendarView)
	},
}

func init() {
	CalendarRootCmd.AddCommand(CalendarViewCmd)
}
