package views

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	ical "github.com/arran4/golang-ical"
	"golang.org/x/exp/slices"
)

type CalendarView struct {
	calendarEntries *ical.Calendar
}

func (c *CalendarView) Draw(out io.Writer) error {

	durationCols := 1
	for len("0000 - 0000") > layoutColCharWidth(durationCols) {
		durationCols = durationCols + 1
	}
	durationWidth := layoutColCharWidth(durationCols)

	remainingCols := 12 - durationCols
	titleCols := int(math.Floor(float64(remainingCols) * 0.66))
	titleWidth := layoutColCharWidth(titleCols)
	roomCols := int(math.Floor(float64(remainingCols) * 0.33))
	roomWidth := layoutColCharWidth(roomCols)

	lineFormatString := "%-" + strconv.Itoa(durationWidth) + "s%-" + strconv.Itoa(titleWidth) + "s%-" + strconv.Itoa(roomWidth) + "s\n"
	out.Write([]byte(fmt.Sprintf(lineFormatString, "time", "meeting", "location")))
	out.Write([]byte(strings.Repeat("-", termWidth) + "\n"))

	events := c.calendarEntries.Events()
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

	for _, entry := range events {

		entryId := entry.Id()

		entryStartTime, err := entry.GetStartAt()
		if err != nil {
			return fmt.Errorf("failed to get entry start time for calendar entry %s: %s", entryId, err)
		}

		entryEndTime, err := entry.GetEndAt()
		if err != nil {
			return fmt.Errorf("failed to get entry end time for calendar entry %s: %s", entryId, err)
		}

		startTime := fmt.Sprintf("%02d%02d", entryStartTime.Hour(), entryStartTime.Minute())
		endTime := fmt.Sprintf("%02d%02d", entryEndTime.Hour(), entryEndTime.Minute())
		formattedTime := fmt.Sprintf("%s - %s", startTime, endTime)

		entryTitle := entry.GetProperty(ical.ComponentProperty(ical.PropertyName)).Value
		entryLocation := entry.GetProperty(ical.ComponentProperty(ical.PropertyLocation)).Value

		out.Write([]byte(fmt.Sprintf(lineFormatString, formattedTime, entryTitle, entryLocation)))
	}

	return nil
}

func (c *CalendarView) SetData(data interface{}) {
	c.calendarEntries = data.(*ical.Calendar)
}

func (c *CalendarView) Data() interface{} {
	return c.calendarEntries
}
