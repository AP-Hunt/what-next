package views

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	"github.com/alexeyco/simpletable"
	ical "github.com/arran4/golang-ical"
	"github.com/isbm/textwrap"

	"github.com/fatih/color"
)

type CalendarView struct {
	data *CalendarViewData
}

type CalendarViewData struct {
	Calendar   *ical.Calendar
	TargetDate time.Time
}

func (c *CalendarView) Draw(out io.Writer) error {
	events := c.data.Calendar.Events()
	err := calendar.SortEventsByStartDateAscending(events)
	if err != nil {
		return err
	}

	boldWhite := color.New(color.FgWhite, color.Bold)
	date := boldWhite.Sprint(c.data.TargetDate.Format("Monday January _2 2006"))

	tbl := simpletable.New()
	tbl.SetStyle(simpletable.StyleCompactLite)
	tbl.Header.Cells = []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "Time"},
		{Align: simpletable.AlignLeft, Text: "Title"},
		{Align: simpletable.AlignLeft, Text: "Room"},
	}

	titleWrapper := textwrap.NewTextWrap()
	titleWrapper.SetWidth(30)

	timeWrapper := textwrap.NewTextWrap()
	timeWrapper.SetWidth(30)

	for _, evt := range events {
		evtId := evt.Id()
		title := evt.GetProperty(ical.ComponentProperty(ical.PropertyName)).Value
		location := evt.GetProperty(ical.ComponentProperty(ical.PropertyLocation)).Value

		startTime, endTime, err := calendar.EventStartAndEnd(evt)
		if err != nil {
			return fmt.Errorf("failed to extract start and end time for evt %s: %s", evtId, err)
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

		timeStr := fmt.Sprintf("%-3s%s - %s", dateMarkers, startStr, endStr)
		timeStr = strings.Join(timeWrapper.Wrap(timeStr), "\n")

		titleStr := strings.Join(
			titleWrapper.Wrap(title),
			"\n",
		)

		tbl.Body.Cells = append(tbl.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: timeStr},
			{Align: simpletable.AlignLeft, Text: titleStr},
			{Align: simpletable.AlignLeft, Text: location},
		})
	}

	fmt.Fprintf(out, "Showing calendar entries for %s\n", date)
	fmt.Fprintf(out, "%s * = event started yesterday, # = event ends tomorrow\n", boldWhite.Sprint("Key:"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, tbl.String())

	return nil
}

func (c *CalendarView) SetData(data interface{}) {
	c.data = data.(*CalendarViewData)
}

func (c *CalendarView) Data() interface{} {
	return c.data
}
