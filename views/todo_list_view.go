package views

import (
	"io"
	"strconv"
	"strings"

	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/golang-module/carbon/v2"
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
	textCols := 6
	textWidth := layoutColCharWidth(textCols)

	rowFormatter, err := fourColRowFormatter([4]int{colLAlign(idCols), colLAlign(completionCols), colLAlign(dueCols), colLAlign(textCols)})
	if err != nil {
		return err
	}

	out.Write([]byte(rowFormatter([4]string{"id", "✓?", "due", "action"})))
	out.Write([]byte(strings.Repeat("-", termWidth)))
	wrapper := textwrap.NewTextWrap()
	wrapper.SetWidth(textWidth)

	for _, item := range v.todoItems.Enumerate() {
		textLines := wrapper.Wrap(item.Action)
		due := ""
		if item.DueDate != nil {
			due = carbon.Time2Carbon(*item.DueDate).Format("dS M y H:i")
		}

		for i, line := range textLines {
			if i == 0 {
				out.Write([]byte(
					rowFormatter([4]string{
						strconv.Itoa(item.Id),
						todoCompletedSymbolMap[item.Completed],
						due,
						line,
					}),
				))
			} else {
				out.Write([]byte(
					rowFormatter([4]string{
						"",
						"",
						"",
						line,
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
