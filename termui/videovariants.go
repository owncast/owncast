package termui

import (
	"context"
	"fmt"

	"github.com/gabek/owncast/config"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
)

func CreateVariantListWidget(ctx context.Context) *text.Text {
	variants := config.Config.GetVideoStreamQualities()

	rollT, err := newRollText(ctx)

	for _, variant := range variants {
		rollT.Write(variant.String() + "\n")
	}

	if err != nil {
		panic(err)
	}
	return rollT
}

func CreateConfigListWidget(ctx context.Context) *text.Text {
	rollT, err := newRollText(ctx)
	rollT.Write(fmt.Sprintf("Stream key: %s\n", config.Config.VideoSettings.StreamingKey))
	rollT.Write(fmt.Sprintf("Segment length: %d\n", config.Config.GetVideoSegmentSecondsLength()))
	rollT.Write(fmt.Sprintf("Playlist length: %d\n", config.Config.GetMaxNumberOfReferencedSegmentsInPlaylist()))

	if config.Config.S3.Enabled {
		rollT.Write(fmt.Sprintf("S3 storage: %s\n", config.Config.S3.Endpoint))
	}

	if err != nil {
		panic(err)
	}
	return rollT
}

func CreateConfigElement(videoVariantsList *text.Text, configSettingsList *text.Text) grid.Element {
	video := grid.Widget(videoVariantsList,
		container.Border(linestyle.Round),
		container.BorderTitle("Video Variants"),
	)

	config := grid.Widget(configSettingsList,
		container.Border(linestyle.Round),
		container.BorderTitle("Configuration"),
	)

	return grid.RowHeightPerc(15,
		grid.ColWidthPerc(50,
			video,
		),
		grid.ColWidthPerc(50,
			config,
		),
	)
}
