package counter

import (
	"sync"
)

// metricManagerInterface defines behavior towards:
// adding a new value for the specific metric
// removing a value
// obtaining cumulative sum of values for the metric - always within the time range
type metricManagerInterface interface {
	addValue(v int) error
	removeValue(t int64)
	update()
}

// MetricData contains data about a single metric type `mt#`
// Name is used for metric recognition
// LastRecord shows the last time a value was set for the metric in Unixtime
// used as a caching value for a faster recognition if data resetting is done
// ValueSet is the set of values in the format - [timeSet]value
type MetricData struct {
	Name        string     `json:"name"`
	FirstRecord int64      `json:"firstRecord"`
	LastRecord  int64      `json:"lastRecord"`
	Sum         int        `json:"sum"`
	ValueSet    timeValues `json:"valueSet"`
}

// Metric holds the metric data and the lock option
type Metric struct {
	Lock *sync.Mutex
	Data *MetricData `json:"data"`
}

// MetricsCounterSetterInterface defines behavior towards:
// creating a new metric object for monitoring and reporting values
// adding a metric value
// removing a metric object from the records
// updating a metric object data accoring to the time range
type MetricsCounterSetterInterface interface {
	AddMetric(name string) error
	AddMetricValue(name string, v int) error
	UpdateMetricName(name, update string) error
	removeMetric(name string) error
}

// MetricCounterGetterInterface defines behavior towards:
// getting a metric sum
type MetricsCounterGetterInterface interface {
	GetMetric(name string) (int, error)
	ListMetrics() metricsList
}

// MetricsData contains information about all metrics currently having values
// this is positive for the last 60 minutes 
// LastRecord is Unixtime and contains the time of the last record
// used as a caching value for faster checks of the time range
type MetricsData struct {
	FirstRecord int64           `json:"firstRecord"`
	LastRecord  int64           `json:"lastRecord"`
	List        metricsListData `json:"list"`
}

// Metrics holds the metrics data and the lock option
// LastRefreshTime keeps the time in Unix for the last metrics refresh performed
type Metrics struct {
	Lock            *sync.Mutex
	SignalRefresher chan bool
	LastRefreshTime int64
	Data            *MetricsData `json:"data"`
}

// Metrics contains records for metrics data - if any
var MetricsRecord Metrics
