package counter

import (
	"time"
	"github.com/argcv/stork/log"
)

// Initialize initializes MetricData manager object
// it is called from the server module at start to populate it
// with imported data
func Initialize() *Metrics {
	// Set up metric object that will receive imported data inside its Data field
	return loadMetricsObject()
}

func (m *Metrics) RunRefresher() {

	go func() {
		for {
			// Sleep for set interval of time between checks
			time.Sleep(time.Duration(CheckRefreshTimeInterval) * time.Second)

			select {
			case <-m.SignalRefresher:
				log.Infof("Refresher stopped")
				return
			default:
				m.refresher()
			}
		}
	}()
}
