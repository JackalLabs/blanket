package blanket

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/JackalLabs/blanket/logger"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/donut"
)

// playType indicates how to play a donut.
type playType int

const (
	playTypePercent playType = iota
	playTypeAbsolute
)

// playDonut continuously changes the displayed percent value on the donut by the
// step once every delay. Exits when the context expires.
func updateSpaceUsage(ctx context.Context, d *donut.Donut, delay time.Duration, logger *logger.Logger, url string) {
	ticker := time.NewTicker(delay)
	running := false

	run := func() {
		var progress int
		var space SpaceResponse
		if running {
			return
		}
		running = true
		// logger.Info("Requesting Metrics...")
		r, err := http.Get(fmt.Sprintf("%s/api/client/space", url))
		if err != nil {
			logger.Error(err)
			return
		}
		running = false

		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(body, &space)
		if err != nil {
			logger.Error(err)
			return
		}

		progress = space.GetPercentUsed()
		if progress == 0 {
			progress = 1
		}

		err = d.Percent(progress)
		if err != nil {
			logger.Error(err)
		}
	}
	run()
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			run()
		case <-ctx.Done():
			return
		}
	}
}

func CmdRunBlanket(url string) {
	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())
	spaceColor := cell.ColorRed
	green, err := donut.New(
		donut.CellOpts(cell.FgColor(spaceColor)),
		donut.TextCellOpts(cell.FgColor(spaceColor)),
		donut.Label("Space Used", cell.FgColor(spaceColor)),
	)

	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger()

	go updateSpaceUsage(ctx, green, time.Second*10, logger, url)

	builder := grid.New()
	builder.Add(
		grid.ColWidthPerc(50,
			grid.Widget(logger.GetWidget(),
				container.Border(linestyle.Light),
				container.BorderTitle(" Logger "),
			),
		),
	)
	builder.Add(
		grid.RowHeightPerc(50,
			grid.Widget(green,
				container.Border(linestyle.Light),
				container.BorderColor(cell.ColorRed),
				container.BorderTitle(" Space "),
				container.FocusedColor(cell.ColorRed),
			),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		panic(fmt.Errorf("builder.Build => %v", err))
	}

	gridOpts = append(
		gridOpts,
		container.Border(linestyle.Light),
		container.BorderTitle(fmt.Sprintf(" Blanket - %s ", url)),
	)

	c, err := container.New(
		t,
		gridOpts...,
	)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(1*time.Second)); err != nil {
		panic(err)
	}
}
