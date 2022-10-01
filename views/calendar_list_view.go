package views

import (
	"fmt"
	"io"

	"github.com/AP-Hunt/what-next/m/calendar"
	"github.com/alexeyco/simpletable"
)

type CalendarListView struct {
	calendars []calendar.CalendarRecord
}

func (cl *CalendarListView) Draw(out io.Writer) error {
	tbl := simpletable.New()
	tbl.SetStyle(simpletable.StyleCompactLite)

	tbl.Header.Cells = []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: "name"},
		{Align: simpletable.AlignLeft, Text: "URL"},
	}

	for _, cal := range cl.calendars {
		row := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: cal.DisplayName},
			{Align: simpletable.AlignLeft, Text: cal.URL},
		}

		tbl.Body.Cells = append(tbl.Body.Cells, row)
	}

	fmt.Fprintln(out, tbl.String())
	return nil
}

func (cl *CalendarListView) SetData(data interface{}) {
	cl.calendars = data.([]calendar.CalendarRecord)
}

func (cl *CalendarListView) Data() interface{} {
	return cl.calendars
}
