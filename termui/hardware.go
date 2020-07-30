package termui

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/linechart"
)

func CreateCPUWidget() *gauge.Gauge {
	ht, err := gauge.New()
	if err != nil {
		panic(err)
	}
	ht.Percent(0)
	return ht
}

func CreateHardwareGaugeWidgetElement(cpuWidget *gauge.Gauge, memoryWidget *gauge.Gauge) grid.Element {
	cpu := grid.Widget(cpuWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("CPU"),
	)

	ram := grid.Widget(memoryWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("Memory"),
	)

	return grid.RowHeightPerc(7,
		grid.ColWidthPerc(50,
			cpu,
		),
		grid.ColWidthPerc(50,
			ram,
		),
	)
}

func CreateMemoryWidget() *gauge.Gauge {
	ht, err := gauge.New()
	if err != nil {
		panic(err)
	}
	ht.Percent(0)
	return ht
}

var _cpuGraphValues = make([]float64, 0)
var _memoryGraphValues = make([]float64, 0)

func setCPUGraphValue(graph *linechart.LineChart, value float64) {
	if len(_cpuGraphValues) > 500 {
		_cpuGraphValues = _cpuGraphValues[1:]
	}
	_cpuGraphValues = append(_cpuGraphValues, value)
	graph.Series("Y", _cpuGraphValues)
}

func CreateCPUGraphWidget() *linechart.LineChart {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)

	if err != nil {
		panic(err)
	}
	return lc
}

func setMemoryGraphValue(graph *linechart.LineChart, value int) {
	if len(_memoryGraphValues) > 500 {
		_memoryGraphValues = _memoryGraphValues[1:]
	}
	_memoryGraphValues = append(_memoryGraphValues, float64(value))
	graph.Series("Y", _memoryGraphValues)
}

func CreateMemoryGraphWidget() *linechart.LineChart {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if err != nil {
		panic(err)
	}
	return lc
}

func CreateHardwareGraphWidgetElement(cpuWidget *linechart.LineChart, memoryWidget *linechart.LineChart) grid.Element {
	cpu := grid.Widget(cpuWidget,
		container.Border(linestyle.Round),
	)

	ram := grid.Widget(memoryWidget,
		container.Border(linestyle.Round),
	)

	return grid.RowHeightPerc(25,
		grid.ColWidthPerc(50,
			cpu,
		),
		grid.ColWidthPerc(50,
			ram,
		),
	)
}
