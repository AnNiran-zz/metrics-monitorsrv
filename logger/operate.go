package logger

import (
	"encoding/json"
	"github.com/argcv/stork/log"

	"monitor/counter"
)

// Import and Export functions are used for calling outside the module
// and keeping symmetry

// Import calls the internal function for importing data from json log
// called from server module at starting
func Import() ([]byte, error) {
	// metricsData is a byte slice of the data
	metricsData, err := importMetricsData()
	if err != nil {
		return nil, err
	}
	return metricsData, nil
}

// Export export metric data to a file
// called from the server at graceful shutdown
func Export(data *counter.Metrics) error {
	if data == nil || (data != nil && len(data.Data.List) == 0) {
		log.Infof(OutputEmptyMetrics())
		return nil
	}

	jsonData, err := json.Marshal(data.Data)
	if err != nil {
		return ErrDataExport(err.Error())
	}

	return exportMetricsData(jsonData)
}
