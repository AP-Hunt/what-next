package main

import (
	"fmt"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "demo-cal-generator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cal := ical.NewCalendar()

			now := time.Now()

			currentEvent := cal.AddEvent("current")
			currentEvent.SetProperty(ical.ComponentProperty(ical.PropertySummary), "Watch a demo of what-next")
			currentEvent.SetStartAt(now.Truncate(5 * time.Minute))
			currentEvent.SetDuration(10 * time.Minute)

			nextEvent := cal.AddEvent("next")
			nextEvent.SetProperty(ical.ComponentProperty(ical.PropertySummary), "Install what-next")
			nextEvent.SetStartAt(now.Truncate(1 * time.Hour).Add(1 * time.Hour))
			nextEvent.SetDuration(1 * time.Hour)

			fmt.Println(cal.Serialize())

			return nil
		},
	}

	rootCmd.Execute()
}
