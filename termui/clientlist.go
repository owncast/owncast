package termui

import (
	"context"
	"fmt"
	"time"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
)

func CreateClientListWidget(ctx context.Context) *text.Text {
	rollT, err := newRollText(ctx)
	if err != nil {
		panic(err)
	}
	return rollT
}

func CreateClientListElement(widget *text.Text) grid.Element {
	return grid.RowHeightPerc(15,
		grid.Widget(widget,
			container.Border(linestyle.Round),
			container.BorderTitle("Clients"),
			container.ID("clientListContainer"),
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
	if _clientList == nil {
		return
	}
	_clientList.Reset()
	for key, _ := range list {
		_clientList.Write(key)
	}

	title := fmt.Sprintf("Clients(%d)", len(list))
	_container.Update("clientListContainer", container.BorderTitle(title))
}
