package blanket

import (
	"context"
	"encoding/json"
	"fmt"
	//nolint
	"io/ioutil"
	"net/http"
	"time"

	"github.com/JackalLabs/blanket/logger"
	"github.com/mum4k/termdash/widgets/text"
)

func updateBurns(ctx context.Context, t *text.Text, delay time.Duration, logger *logger.Logger, url string, api string) {
	ticker := time.NewTicker(delay)
	running := false

	run := func() {
		var index IndexResponse
		var provider Providers

		if running {
			return
		}
		running = true
		r, err := http.Get(url)
		if err != nil {
			logger.Error(err)
			return
		}
		running = false

		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(body, &index)
		if err != nil {
			logger.Error(err)
			return
		}

		dealRes, err := http.Get(fmt.Sprintf("%s%s/providers/%s", api, storageAPI, index.Address))
		if err != nil {
			logger.Error(err)
			return
		}
		running = false

		dealBody, readErr := ioutil.ReadAll(dealRes.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(dealBody, &provider)
		if err != nil {
			logger.Error(err)
			return
		}

		t.Reset()
		_ = t.Write(provider.Providers.BurnedContracts)
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

func buildBurns() *text.Text {
	borderless, err := text.New()
	if err != nil {
		panic(err)
	}

	return borderless
}
