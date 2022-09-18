package views

import (
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	ical "github.com/arran4/golang-ical"
	"github.com/isbm/textwrap"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

type CalendarView struct {
	data *CalendarViewData
}

type CalendarViewData struct {
	Calendar   *ical.Calendar
	TargetDate time.Time
}

func (c *CalendarView) Draw(out io.Writer) error {
	type tableRow struct {
		duration string
		title    string
		location string
	}

	events := c.data.Calendar.Events()
	slices.SortFunc(events, func(evtA *ical.VEvent, evtB *ical.VEvent) bool {
		aStart, err := evtA.GetStartAt()
		if err != nil {
			panic(fmt.Sprintf("sorting calendar entries for view: %s", err))
		}

		bStart, err := evtB.GetStartAt()
		if err != nil {
			panic(fmt.Sprintf("sorting calendar entries for view: %s", err))
		}

		if aStart.Equal(bStart) {
			aEnd, err := evtA.GetEndAt()
			if err != nil {
				panic(fmt.Sprintf("sorting calendar entries for view: %s", err))
			}

			bEnd, err := evtB.GetEndAt()
			if err != nil {
				panic(fmt.Sprintf("sorting calendar entries for view: %s", err))
			}

			aDuration := aEnd.Sub(aStart)
			bDuration := bEnd.Sub(bStart)

			if aDuration == bDuration || aDuration > bDuration {
				return true
			} else {
				return false
			}
		} else if aStart.Before(bStart) {
			return true
		} else {
			return false
		}
	})

	boldWhite := color.New(color.FgWhite, color.Bold)
	date := boldWhite.Sprint(c.data.TargetDate.Format("Monday January _2 2006"))

	// Build list of rows before drawing so we can work out the widest column
	rows := []tableRow{}
	longestDurationStrLen := 0
	for _, evt := range events {
		row := tableRow{}

		entryId := evt.Id()
		row.title = evt.GetProperty(ical.ComponentProperty(ical.PropertyName)).Value
		row.location = evt.GetProperty(ical.ComponentProperty(ical.PropertyLocation)).Value

		startTime, err := evt.GetStartAt()
		if err != nil {
			return fmt.Errorf("failed to get entry start time for calendar entry %s: %s", entryId, err)
		}

		endTime, err := evt.GetEndAt()
		if err != nil {
			return fmt.Errorf("failed to get entry end time for calendar entry %s: %s", entryId, err)
		}

		startsToday, err := calendar.EventStartsToday(evt)
		if err != nil {
			return err
		}

		endsToday, err := calendar.EventEndsToday(evt)
		if err != nil {
			return err
		}

		startStr := fmt.Sprintf("%02d%02d", startTime.Hour(), startTime.Minute())
		endStr := fmt.Sprintf("%02d%02d", endTime.Hour(), endTime.Minute())
		dateMarkers := " "
		if !startsToday {
			dateMarkers = "*" + dateMarkers
		}

		if !endsToday {
			dateMarkers = "#" + dateMarkers
		}

		row.duration = fmt.Sprintf("%-3s%s - %s", dateMarkers, startStr, endStr)

		if len(row.duration) > longestDurationStrLen {
			longestDurationStrLen = len(row.duration)
		}

		rows = append(rows, row)
	}

	durationCols := colsRequiredToFitChars(longestDurationStrLen)
	remainingCols := 12 - durationCols
	titleCols := int(math.Floor(float64(remainingCols) * 0.66))
	titleWidth := layoutColCharWidth(titleCols)
	roomCols := int(math.Floor(float64(remainingCols) * 0.33))

	out.Write([]byte(fmt.Sprintf("Showing calendar entries for %s\n", date)))
	out.Write([]byte(fmt.Sprintf("%s * = event started yesterday, # = event ends tomorrow\n", boldWhite.Sprint("Key:"))))
	out.Write([]byte("\n"))

	rowFormatter, err := threeColRowFormatter([3]int{durationCols * -1, titleCols * -1, roomCols * -1})
	if err != nil {
		return err
	}

	headerRow := boldWhite.Sprint(rowFormatter([3]string{"time", "meeting", "location"}))
	out.Write([]byte(headerRow))

	out.Write([]byte(strings.Repeat("-", termWidth) + "\n"))

	wrapper := textwrap.NewTextWrap()
	wrapper.SetWidth(titleWidth)

	for _, r := range rows {
		wrappedTitle := wrapper.Wrap(r.title)

		for i, titleLine := range wrappedTitle {
			switch i {
			case 0:
				out.Write([]byte(rowFormatter([3]string{r.duration, titleLine, r.location})))
			default:
				out.Write([]byte(rowFormatter([3]string{"", titleLine, ""})))
			}

		}
	}

	return nil
}

func (c *CalendarView) SetData(data interface{}) {
	c.data = data.(*CalendarViewData)
}

func (c *CalendarView) Data() interface{} {
	return c.data
}
