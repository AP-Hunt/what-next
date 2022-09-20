package views

import (
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/golang-module/carbon/v2"
	"github.com/hako/durafmt"
	"github.com/isbm/textwrap"
)

type TodoListView struct {
	todoItems *todo.TodoItemCollection
}

func (v *TodoListView) Draw(out io.Writer) error {
	todoCompletedSymbolMap := map[bool]string{
		true:  "[✓]",
		false: "[ ]",
	}

	idCols := 2
	completionCols := 1
	dueCols := 3
	durationCols := 2
	textCols := 4
	textWidth := layoutColCharWidth(textCols)

	rowFormatter, err := fiveColRowFormatter([5]int{
		colLAlign(idCols),
		colLAlign(completionCols),
		colLAlign(dueCols),
		colLAlign(durationCols),
		colLAlign(textCols),
	})
	if err != nil {
		return err
	}

	out.Write([]byte(rowFormatter([5]string{"id", "✓?", "due", "duration", "action"})))
	out.Write([]byte(strings.Repeat("-", termWidth)))

	textWrapper := textwrap.NewTextWrap()
	textWrapper.SetWidth(textWidth)

	durationWrapper := textwrap.NewTextWrap()
	durationWrapper.SetWidth(layoutColCharWidth(durationCols))

	for _, item := range v.todoItems.Enumerate() {
		textLines := textWrapper.Wrap(item.Action)

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
		maxLines := math.Max(float64(len(textLines)), float64(len(durationLines)))

		for i := 0; i < int(maxLines); i++ {
			textLine := ""
			durationLine := ""

			if i < len(textLines) {
				textLine = textLines[i]
			}

			if i < len(durationLines) {
				durationLine = durationLines[i]
			}

			if i == 0 {
				out.Write([]byte(
					rowFormatter([5]string{
						strconv.Itoa(item.Id),
						todoCompletedSymbolMap[item.Completed],
						due,
						durationLine,
						textLine,
					}),
				))
			} else {
				out.Write([]byte(
					rowFormatter([5]string{
						"",
						"",
						"",
						durationLine,
						textLine,
					}),
				))
			}
		}
	}

	return nil
}

func (v *TodoListView) SetData(data interface{}) {
	v.todoItems = data.(*todo.TodoItemCollection)
}
func (v *TodoListView) Data() interface{} {
	return v.todoItems
}
