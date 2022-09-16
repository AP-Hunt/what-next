package views

import (
	"fmt"
	"io"
	"strconv"

	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/isbm/textwrap"
)

type TodoListView struct {
	todoItems *todo.TodoItemCollection
}

func (v *TodoListView) Draw(out io.Writer) error {
	todoCompletedSymbolMap := map[bool]string{
		true:  "[✓]",
		false: "[x]",
	}

	idCols := 2
	completionCols := 1
	textCols := 9

	idWidth := layoutColCharWidth(idCols)
	completionWidth := layoutColCharWidth(int(completionCols))
	textWidth := layoutColCharWidth(textCols)

	formatString := fmt.Sprintf("%%-%ss %%%ss %%-%ss\n", strconv.Itoa(idWidth), strconv.Itoa(completionWidth), strconv.Itoa(textWidth))
	out.Write([]byte(fmt.Sprintf(formatString, "id", "", "action")))

	wrapper := textwrap.NewTextWrap()
	wrapper.SetWidth(textWidth)

	for _, item := range v.todoItems.Enumerate() {
		textLines := wrapper.Wrap(item.Action)
		for i, line := range textLines {
			if i == 0 {
				out.Write([]byte(
					fmt.Sprintf(
						formatString,
						strconv.Itoa(item.Id),
						todoCompletedSymbolMap[item.Completed],
						line,
					),
				))
			} else {
				out.Write([]byte(
					fmt.Sprintf(
						formatString,
						"",
						"",
						line,
					),
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
