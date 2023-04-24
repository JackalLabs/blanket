package blanket

import (
	"context"
	"fmt"
	"time"

	"github.com/JackalLabs/blanket/logger"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/donut"
)

func build(ctx context.Context, t terminalapi.Terminal, url string, api string) (*container.Container, error) {
	spaceDonut, err := donut.New(
		donut.CellOpts(cell.FgColor(cell.ColorRed)),
		donut.TextCellOpts(cell.FgColor(cell.ColorRed)),
		donut.Label("Space Used", cell.FgColor(cell.ColorRed)),
	)
	blockTimeDonut, err := donut.New(
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
		donut.TextCellOpts(cell.FgColor(cell.ColorBlue)),
		donut.Label("Proof Window", cell.FgColor(cell.ColorBlue)),
	)

	if err != nil {
		panic(err)
	}

	bLogger := logger.NewLogger()
	balance := buildBalance()
	ratio := buildDealRatio()
	burns := buildBurns()

	go updateSpaceUsage(ctx, spaceDonut, time.Second*30, bLogger, url)
	go updateBalance(ctx, balance, time.Second*60, bLogger, url)
	go updateRatio(ctx, ratio, time.Second*60, bLogger, api)
	go updateBurns(ctx, burns, time.Second*60, bLogger, url, api)
	go updateBlockTime(ctx, blockTimeDonut, time.Second*10, bLogger, api)

	cc := container.SplitHorizontal(
		container.Top(
			container.SplitVertical(
				container.Left(
					container.SplitHorizontal(
						container.Top(
							container.PlaceWidget(burns),
							container.Border(linestyle.Light),
							container.BorderTitle("Burned Contracts"),
						),
						container.Bottom(
							container.PlaceWidget(ratio),
							container.Border(linestyle.Light),
							container.BorderTitle("Deal/Stray Ratio"),
						),
					),
				),
				container.Right(
					container.SplitHorizontal(
						container.Top(
							container.PlaceWidget(balance),
							container.Border(linestyle.Light),
							container.BorderTitle("Balance"),
						),
						container.Bottom(),
					),
				),
			),
		),
		container.Bottom(
			container.SplitVertical(

				container.Left(
					container.PlaceWidget(bLogger.GetWidget()),
					container.Border(linestyle.Light),
					container.BorderTitle("Logger"),
				),
				container.Right(
					container.SplitVertical(
						container.Left(
							container.PlaceWidget(spaceDonut),
							container.Border(linestyle.Light),
							container.BorderColor(cell.ColorRed),
							container.BorderTitle("Space"),
							container.FocusedColor(cell.ColorRed),
						),
						container.Right(
							container.PlaceWidget(blockTimeDonut),
							container.Border(linestyle.Light),
							container.BorderColor(cell.ColorBlue),
							container.BorderTitle("Proof Window"),
							container.FocusedColor(cell.ColorBlue),
						)),
				),
			)),
	)

	c, err := container.New(
		t,
		cc,
		container.Border(linestyle.Light),
		container.BorderTitle(fmt.Sprintf("Blanket - %s", url)),
	)

	return c, err
}

func CmdRunBlanket(url string, api string) {
	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())

	c, err := build(ctx, t, url, api)
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
