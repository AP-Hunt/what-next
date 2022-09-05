package views

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/isbm/textwrap"
	"golang.org/x/crypto/ssh/terminal"
)

type TodoListView struct {
	todoItems *todo.TodoItemCollection
}

func (v *TodoListView) Draw(out io.Writer) error {
	todoCompletedSymbolMap := map[bool]string{
		true:  "[✓]",
		false: "[x]",
	}

	width, _, err := terminal.GetSize(0)
	if err != nil {
		return fmt.Errorf("getting size of terminal: %s", err)
	}

	oneColWidth := math.Floor(float64(width / 12))
	idCols := float64(2)
	completionCols := float64(1)
	textCols := float64(9)

	idWidth := int(oneColWidth * idCols)
	completionWidth := int(oneColWidth * completionCols)
	textWidth := int(oneColWidth * textCols)

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
