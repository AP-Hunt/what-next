package views

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	ical "github.com/arran4/golang-ical"
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

	for _, entry := range c.calendarEntries.Events() {

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
