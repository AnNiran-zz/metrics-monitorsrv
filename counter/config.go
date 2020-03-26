package counter

import (
	"fmt"
)

// timeValues is the way we sort times and values of metrics
type timeValues map[int64]int

// metricsListData is the format in which metrics and their values are saved inside the manager
type metricsListData map[string]*Metric

// metricsList defines the format metrics list is returned
type metricsList map[string]int

// TimeRange is used for comparing times in Unixtime format
var TimeRange = 60*60

// RefreshTimeInterval is the interval on which the ticker is called
var RefreshTimeInterval int64 = int64(1*60)

// CheckRefreshTimeInterval is the interval on which the last refresh time is checked
var CheckRefreshTimeInterval = 2

// Errors 
var (
	ErrMetricExists = func(name string) error {
		return fmt.Errorf("Metric with name: %s already exists", name)
	}
	ErrMetricNotExisting = func(name string) error {
		return fmt.Errorf("Metric with name: %s does not exist", name)
	}
	ErrSameSecondRecord = func(name string) error {
		return fmt.Errorf("Record for metric %s with same time already exists", name)
	} 
)
