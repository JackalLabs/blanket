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
	"github.com/mum4k/termdash/widgets/donut"
)

func (s *SpaceResponse) GetPercentUsed() int {
	if s.TotalSpace == 0 {
		return 1
	}
	return int((s.UsedSpace * 100) / s.TotalSpace)
}

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
