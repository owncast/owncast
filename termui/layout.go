package termui

import (
	"context"
	"os"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

const redrawInterval = 5 * time.Second

// Terminal implementations
const (
	termboxTerminal = "termbox"
	tcellTerminal   = "tcell"
)

var versionString string

// rootID is the ID assigned to the root container.
const rootID = "root"

var _headerText *text.Text
var _currentStreamText *text.Text
var _clientList *text.Text
var _videoVariantList *text.Text
var _configSettingsList *text.Text
var _cpuGauge *gauge.Gauge
var _memoryGauge *gauge.Gauge
var _cpuGraph *linechart.LineChart
var _memoryGraph *linechart.LineChart
var _textInput *textinput.TextInput
var _logText *text.Text

var _ctx context.Context

func Setup(version string) error {
	versionString = version

	ctx, cancel := context.WithCancel(context.Background())
	_ctx = ctx

	var t terminalapi.Terminal
	t, err := termbox.New(termbox.ColorMode(terminalapi.ColorMode256))

	if err != nil {
		cancel()
		return err
	}

	defer func() {
		t.Close()
		os.Exit(0)
	}()

	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		panic(err)
	}

	layout, _ := GetLayout()

	if err := c.Update(rootID, layout...); err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			cancel()
		}
	}

	ticker := time.NewTicker(redrawInterval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cpuUsage := getCPUUsage()
				_cpuGauge.Percent(cpuUsage)
				setCPUGraphValue(_cpuGraph, float64(cpuUsage))

				memUsage := getMemoryUsage()
				_memoryGauge.Percent(int(memUsage))
				setMemoryGraphValue(_memoryGraph, memUsage)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}

	return nil
}

// GetLayout creates a new layout container.
func GetLayout() ([]container.Option, error) {
	leftRows := getLeftRows()

	builder := grid.New()
	builder.Add(
		grid.ColWidthPerc(70, leftRows...),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

func getLeftRows() []grid.Element {
	_headerText = CreateHeaderTextWidget()
	headerTextLayout := CreateHeaderTextWidgetElement(_headerText)

	_currentStreamText = CreateStreamInfoTextWidget()
	streamInfoTextWidget := CreateStreamInfoTextWidgetElement(_currentStreamText)

	_clientList = CreateClientListWidget(_ctx)
	clientListBox := CreateClientListElement(_clientList)

	_videoVariantList = CreateVariantListWidget(_ctx)
	_configSettingsList = CreateConfigListWidget(_ctx)

	configurations := CreateConfigElement(_videoVariantList, _configSettingsList)

	_cpuGauge = CreateCPUWidget()
	_memoryGauge = CreateMemoryWidget()
	hardwareGauges := CreateHardwareGaugeWidgetElement(_cpuGauge, _memoryGauge)

	_cpuGraph = CreateCPUGraphWidget()
	_memoryGraph = CreateMemoryGraphWidget()
	hardwareGraphs := CreateHardwareGraphWidgetElement(_cpuGraph, _memoryGraph)

	_logText = CreateLogWidget(_ctx)
	logText := CreateLogElement(_logText)

	// _textInput = CreateMessageInputWidget()
	// textInput := CreateMessageInputWidgetElement(_textInput)

	leftRows := []grid.Element{headerTextLayout, streamInfoTextWidget, clientListBox, configurations, hardwareGauges, hardwareGraphs, logText}

	return leftRows
}
