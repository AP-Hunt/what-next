package views

import (
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
	textWidth := layoutColCharWidth(textCols)

	rowFormatter, err := threeColRowFormatter([3]int{idCols * -1, completionCols, textCols * -1})
	if err != nil {
		return err
	}

	out.Write([]byte(rowFormatter([3]string{"id", "", "action"})))

	wrapper := textwrap.NewTextWrap()
	wrapper.SetWidth(textWidth)

	for _, item := range v.todoItems.Enumerate() {
		textLines := wrapper.Wrap(item.Action)
		for i, line := range textLines {
			if i == 0 {
				out.Write([]byte(
					rowFormatter([3]string{
						strconv.Itoa(item.Id),
						todoCompletedSymbolMap[item.Completed],
						line,
					}),
				))
			} else {
				out.Write([]byte(
					rowFormatter([3]string{
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
