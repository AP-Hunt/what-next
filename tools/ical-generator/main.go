package main

import (
	"fmt"
	"strconv"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "ical-generator num_entries",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				return err
			}

			numEntries := int(n)

			if numEntries <= 0 {
				return fmt.Errorf("num args must be a positive integer")
			}

			cal := ical.NewCalendar()

			midnight := time.Now().Truncate(24 * time.Hour)
			tomorrow := midnight.Add(24 * time.Hour)

			middayYesterday := midnight.Add(-12 * time.Hour)
			middayTomorrow := tomorrow.Add(12 * time.Hour)

			for i := 1; i <= numEntries; i++ {
				event := cal.AddEvent(fmt.Sprintf("evt-%d", i))

				event.SetProperty(ical.ComponentProperty(ical.PropertySummary), gofakeit.Phrase())
				event.SetProperty(ical.ComponentProperty(ical.PropertyLocation), fmt.Sprintf("Room %d", gofakeit.Number(0, 101)))

				event.SetStartAt(gofakeit.DateRange(middayYesterday, middayTomorrow))
				event.SetDuration(time.Duration(gofakeit.Number(5, 181)) * time.Minute)
			}

			fmt.Println(cal.Serialize())

			return nil
		},
	}

	rootCmd.Execute()
}
