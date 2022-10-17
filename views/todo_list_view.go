package views

import (
	"fmt"
	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	"github.com/golang-module/carbon/v2"
	"github.com/hako/durafmt"
	"github.com/isbm/textwrap"
	"io"
	"strconv"
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

	overdueStyle := color.New(color.FgRed, color.Bold)

	textWrapper := textwrap.NewTextWrap()
	textWrapper.SetWidth(75)

	durationWrapper := textwrap.NewTextWrap()
	durationWrapper.SetWidth(30)

	for _, item := range v.todoItems.Enumerate() {
		formattedAction := textWrapper.Fill(item.Action)

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

		if item.IsOverdue() {
			due = overdueStyle.Sprint(due)
		}

		duration := ""
		if item.Duration != nil {
			duration = durafmt.Parse(*item.Duration).String()
		}

		formattedDuration := durationWrapper.Fill(duration)

		row := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(item.Id)},
			{Align: simpletable.AlignCenter, Text: todoCompletedSymbolMap[item.Completed]},
			{Align: simpletable.AlignLeft, Text: due},
			{Align: simpletable.AlignLeft, Text: formattedDuration},
			{Align: simpletable.AlignLeft, Text: formattedAction},
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
