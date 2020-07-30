package termui

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
)

func CreateHeaderTextWidget() *text.Text {
	ht, err := text.New()
	if err != nil {
		panic(err)
	}

	if err := ht.Write(versionString); err != nil {
		panic(err)
	}
	// if err := ht.Write("Owncast v0.0.0-linux-blah. Uptime 3:14:12:38."); err != nil {
	// 	panic(err)
	// }

	return ht
}

func CreateHeaderTextWidgetElement(widget *text.Text) grid.Element {
	return grid.RowHeightPerc(5,
		grid.Widget(widget,
			container.Border(linestyle.None),
			container.BorderTitle("Press Esc to quit"),
		),
	)
}
