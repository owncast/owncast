package termui

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
)

func CreateStreamInfoTextWidget() *text.Text {
	ht, err := text.New(text.WrapAtWords())
	if err != nil {
		panic(err)
	}

	return ht
}

func SetCurrentInboundStream(data string) {
	if _currentStreamText == nil {
		return
	}

	_currentStreamText.Reset()
	_currentStreamText.Write(data)
}

func CreateStreamInfoTextWidgetElement(widget *text.Text) grid.Element {
	return grid.RowHeightPerc(13,
		grid.Widget(widget,
			container.Border(linestyle.Round),
			container.BorderTitle("Live Stream"),
		),
	)
}
