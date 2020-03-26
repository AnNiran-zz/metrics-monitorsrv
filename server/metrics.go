package server

import (
	"encoding/json"
	"github.com/argcv/stork/log"

	"moninot/counter"
	"monitor/logger"
)

// startMetricsMonitor initializes MetricData manager object
// and calls the logger import to populate it with data if any
// called from the server module
//
// loading from the current records corresponds to the fact that
// stopping and starting the server might happen within the time range
// and the logged data is still relevant for requests
func startMetricsMonitor() error {
	metrics = counter.Initialize()
	
	msData, err := logger.Import()
	if err != nil {
		return err
	}

	// if there is data, load it inside the object
	if msData != nil {
		err := json.Unmarshal(msData, metrics.Data)
		if err != nil {
			return err
		}
	}

	// if no error is returned the metrics object contains the imported data now
	return nil
}

// logMetrics logs current metrics data using the logger package
func logMetrics() error {
	err := logger.Export(metrics)
	if err != nil {
		log.Error(ErrExportMetrics(err.Error()))
	}

	return nil
}
