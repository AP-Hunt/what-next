package views

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/alexeyco/simpletable"
	"github.com/golang-module/carbon/v2"
	"github.com/hako/durafmt"
	"github.com/isbm/textwrap"
)

type TodoListView struct {
	todoItems *todo.TodoItemCollection
}

func (v *TodoListView) Draw(out io.Writer) error {
	todoCompletedSymbolMap := map[bool]string{
		true:  "✓",
		false: "",
	}

	tbl := simpletable.New()
	tbl.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "✓?"},
			{Align: simpletable.AlignLeft, Text: "Due"},
			{Align: simpletable.AlignLeft, Text: "Duration"},
			{Align: simpletable.AlignLeft, Text: "Action"},
		},
	}
	tbl.SetStyle(simpletable.StyleCompactLite)

	textWrapper := textwrap.NewTextWrap()
	textWrapper.SetWidth(30)

	durationWrapper := textwrap.NewTextWrap()
	durationWrapper.SetWidth(30)

	for _, item := range v.todoItems.Enumerate() {
		textLines := textWrapper.Wrap(item.Action)
		multiLineText := strings.Join(textLines, "\n")

		due := ""
		if item.DueDate != nil {
			carbonDate := carbon.Time2Carbon(*item.DueDate)

			if carbonDate.IsYesterday() {
				due = "Yesterday"
			} else if carbonDate.IsToday() {
				due = "Today"
			} else if carbonDate.IsTomorrow() {
				due = "Tomorrow"
			} else {
				due = carbonDate.Format("dS M y H:i")
			}
		}

		duration := ""
		if item.Duration != nil {
			duration = durafmt.Parse(*item.Duration).String()
		}

		durationLines := durationWrapper.Wrap(duration)
		multiLineDuration := strings.Join(durationLines, "\n")

		row := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(item.Id)},
			{Align: simpletable.AlignCenter, Text: todoCompletedSymbolMap[item.Completed]},
			{Align: simpletable.AlignLeft, Text: due},
			{Align: simpletable.AlignLeft, Text: multiLineDuration},
			{Align: simpletable.AlignLeft, Text: multiLineText},
		}

		tbl.Body.Cells = append(tbl.Body.Cells, row)
	}

	fmt.Fprintln(out, tbl.String())

	return nil
}

func (v *TodoListView) SetData(data interface{}) {
	v.todoItems = data.(*todo.TodoItemCollection)
}
func (v *TodoListView) Data() interface{} {
	return v.todoItems
}
