package termui

import (
	"fmt"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/textinput"
)

func CreateMessageInputWidget() *textinput.TextInput {
	input, err := textinput.New(
		textinput.Label("Send chat message as Owncast: ", cell.FgColor(cell.ColorBlue)),
		textinput.MaxWidthCells(50),
		textinput.PlaceHolder("enter any text"),
		textinput.OnSubmit(func(text string) error {
			fmt.Println(text)
			return nil
		}),
		textinput.ClearOnSubmit(),
	)

	if err != nil {
		panic(err)
	}

	return input
}

func CreateMessageInputWidgetElement(widget *textinput.TextInput) grid.Element {
	return grid.RowHeightPerc(4,
		grid.Widget(widget,
			container.Border(linestyle.None),
		),
	)
}
