package termui

import (
	"fmt"
	"time"

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
			container.ID("streamInfo"),
			container.Border(linestyle.Round),
			container.BorderTitle("Live Stream Disconnected"),
		),
	)
}

func updateStreamInfoTitle() {
	title := "Live Stream Disconnected"

	if StreamConnectedAt.Valid {
		totalSeconds := int(time.Since(StreamConnectedAt.Time).Seconds())
		hours := int(totalSeconds / (60 * 60) % 24)
		minutes := int(totalSeconds/60) % 60
		seconds := int(totalSeconds % 60)
		title = fmt.Sprintf("Live Stream Connected %02d:%02d:%02d", hours, minutes, seconds)
	}

	_container.Update("streamInfo", container.BorderTitle(title))

}
