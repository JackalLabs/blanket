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

const storageAPI = "/jackal-dao/canine-chain/storage/"

func updateRatio(ctx context.Context, t *text.Text, delay time.Duration, logger *logger.Logger, api string) {
	ticker := time.NewTicker(delay)
	running := false

	run := func() {

		var strays APIResponse
		var deals APIResponse

		if running {
			return
		}
		running = true
		r, err := http.Get(fmt.Sprintf("%s%s/strays", api, storageAPI))
		if err != nil {
			logger.Error(err)
			return
		}
		running = false

		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(body, &strays)
		if err != nil {
			logger.Error(err)
			return
		}

		dealRes, err := http.Get(fmt.Sprintf("%s%s/active_deals", api, storageAPI))
		if err != nil {
			logger.Error(err)
			return
		}
		running = false

		dealBody, readErr := ioutil.ReadAll(dealRes.Body)
		if readErr != nil {
			logger.Error(readErr)
		}

		err = json.Unmarshal(dealBody, &deals)
		if err != nil {
			logger.Error(err)
			return
		}

		fcount, _ := strconv.ParseFloat(deals.Pagination.Count, 64)
		scount, _ := strconv.ParseFloat(strays.Pagination.Count, 64)

		t.Reset()
		t.Write(fmt.Sprintf("%s : %s | %.2f%% healthy", deals.Pagination.Count, strays.Pagination.Count, 100*(fcount/(scount+fcount))))

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

func buildDealRatio() *text.Text {
	borderless, err := text.New()
	if err != nil {
		panic(err)
	}

	return borderless
}
