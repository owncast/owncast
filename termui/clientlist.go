package termui

import (
	"context"
	"fmt"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

func CreateClientListWidget(ctx context.Context) *text.Text {
	rollT, err := newRollText(ctx)
	if err != nil {
		panic(err)
	}
	return rollT
}

func CreateClientListElement(widget *text.Text, graph *linechart.LineChart) grid.Element {
	list := grid.Widget(widget,
		container.Border(linestyle.Light),
		container.BorderTitle("Clients"),
		container.BorderColor(cell.ColorBlue),
	)

	clientsGraph := grid.Widget(graph,
		container.Border(linestyle.Light),
		container.BorderColor(cell.ColorBlue),
	)

	return grid.RowHeightPerc(17,
		grid.ColWidthPerc(60,
			list,
		),
		grid.ColWidthPerc(40,
			clientsGraph,
		),
	)
}

func newRollText(ctx context.Context) (*text.Text, error) {
	t, err := text.New(text.RollContent())
	if err != nil {
		return nil, err
	}

	return t, nil
}

// SetClientList draws the list of clients in the termui
func SetClientList(list map[string]time.Time) {
	if _clientListText == nil {
		return
	}
	_clientListText.Reset()
	for key, _ := range list {
		_clientListText.Write(key + "\n")
	}

	count := len(list)
	title := fmt.Sprintf("Clients(%d)", count)
	_container.Update("clientListContainer", container.BorderTitle(title))

	_clientList = list
}

func CreateClientGraphWidget() *linechart.LineChart {
	lc, err := linechart.New(
		linechart.AxesCellOpts(
			cell.FgColor(cell.ColorBlue),
		),

		linechart.YLabelCellOpts(
			cell.FgColor(cell.ColorWhite),
		),

		linechart.YAxisFormattedValues(linechart.ValueFormatterRound),
	)

	if err != nil {
		panic(err)
	}
	return lc
}

var _clientList map[string]time.Time
var _clientListCounterValues = make([]float64, 0)

func updateClientGraphValues(graph *linechart.LineChart) {
	count := len(_clientList)
	_clientListCounterValues = append(_clientListCounterValues, float64(count))
	graph.Series("Clients", _clientListCounterValues)
}
