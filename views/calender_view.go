package views

import (
	"fmt"
	"io"
	"strconv"

	ical "github.com/arran4/golang-ical"
)

type CalendarView struct {
	calendarEntries *ical.Calendar
}

func (c *CalendarView) Draw(out io.Writer) error {

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

		entryDuration := entryEndTime.Sub(entryStartTime)

		startTime := fmt.Sprintf("%02d%02d", entryStartTime.Hour(), entryStartTime.Minute())
		duration := fmt.Sprintf("%.0f minutes", entryDuration.Minutes())
		durationStrLen := len(duration)
		titleLineFmtString := "%-" + strconv.Itoa(durationStrLen+1) + "s %s\n"

		entryTitle := entry.GetProperty(ical.ComponentProperty(ical.PropertyName)).Value
		entryLocation := entry.GetProperty(ical.ComponentProperty(ical.PropertyLocation)).Value

		out.Write([]byte(fmt.Sprintf(titleLineFmtString, startTime, entryTitle)))
		out.Write([]byte(fmt.Sprintf("%s  %s\n", duration, entryLocation)))
		out.Write([]byte("\n"))
	}

	return nil
}

func (c *CalendarView) SetData(data interface{}) {
	c.calendarEntries = data.(*ical.Calendar)
}

func (c *CalendarView) Data() interface{} {
	return c.calendarEntries
}
