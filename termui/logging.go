package termui

import (
	"context"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
)

func CreateLogWidget(ctx context.Context) *text.Text {
	rollT, err := newRollText(ctx)
	if err != nil {
		panic(err)
	}
	return rollT
}

func CreateLogElement(widget *text.Text) grid.Element {
	return grid.RowHeightPerc(10,
		grid.Widget(widget,
			container.Border(linestyle.Round),
			container.BorderTitle("Log"),
		),
	)
}

func Log(text string) {
	if _logText != nil {
		_logText.Write(text)
	}
}
