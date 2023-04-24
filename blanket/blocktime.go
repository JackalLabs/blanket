package blanket

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/JackalLabs/blanket/logger"
	"github.com/mum4k/termdash/widgets/donut"
)

func updateBlockTime(ctx context.Context, d *donut.Donut, delay time.Duration, logger *logger.Logger, url string) {
	ticker := time.NewTicker(delay)
	running := false

	run := func() {
		var block BlockResponse
		var params ParamResponse
		if running {
			return
		}
		running = true
		// logger.Info("Requesting Metrics...")
		r, err := http.Get(fmt.Sprintf("%s/blocks/latest", url))
		if err != nil {
			logger.Error(err)
			return
		}

		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(body, &block)
		if err != nil {
			logger.Error(err)
			return
		}

		r, err = http.Get(fmt.Sprintf("%s/jackal-dao/canine-chain/storage/params", url))
		if err != nil {
			logger.Error(err)
			return
		}
		body, readErr = ioutil.ReadAll(r.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(body, &params)
		if err != nil {
			logger.Error(err)
			return
		}

		running = false

		window, err := strconv.ParseInt(params.Params.ProofWindow, 10, 64)
		if err != nil {
			logger.Error(err)
			return
		}

		height, err := strconv.ParseInt(block.Block.Header.Height, 10, 64)
		if err != nil {
			logger.Error(err)
			return
		}

		timing := height % window

		amount := int(100 * (float64(timing) / float64(window)))
		if amount == 0 {
			amount = 1
		}
		err = d.Percent(amount)
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
