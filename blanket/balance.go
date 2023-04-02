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
	"github.com/mum4k/termdash/widgets/text"
)

const priceURL = "https://api-osmosis.imperator.co/tokens/v2/price/jkl"

func updateBalance(ctx context.Context, t *text.Text, delay time.Duration, logger *logger.Logger, url string) {
	ticker := time.NewTicker(delay)
	running := false

	run := func() {
		var balance BalanceResponse
		var price PriceResponse

		if running {
			return
		}
		running = true
		// logger.Info("Requesting Metrics...")
		r, err := http.Get(fmt.Sprintf("%s/api/network/balance", url))
		if err != nil {
			logger.Error(err)
			return
		}
		running = false

		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(body, &balance)
		if err != nil {
			logger.Error(err)
			return
		}

		priceRes, err := http.Get(priceURL)
		if err != nil {
			logger.Error(err)
			return
		}
		running = false

		priceBody, readErr := ioutil.ReadAll(priceRes.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(priceBody, &price)
		if err != nil {
			logger.Error(err)
			return
		}

		bal, err := strconv.ParseInt(balance.Balance.Amount, 10, 64)
		if err != nil {
			logger.Error(err)
			return
		}

		jklBalance := float64(bal) / 1000000

		tokenPrice := jklBalance * price.Price

		t.Reset()
		t.Write(fmt.Sprintf("%.2fjkl ≈ $%.2f", jklBalance, tokenPrice))

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

func buildBalance() *text.Text {
	borderless, err := text.New()
	if err != nil {
		panic(err)
	}

	return borderless
}
