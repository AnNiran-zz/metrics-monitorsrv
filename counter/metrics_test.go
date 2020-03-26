package counter

import (
	//"sync"
)

type testTimeValues map[int64]int
type testMetricsList map[string]int

type TestMetricData struct {
	Name        string     `json:"name"`
	FirstRecord int64      `json:"firstRecord"`
	LastRecord  int64      `json:"lastRecord"`
	Sum         int        `json:"sum"`
	ValueSet    timeValues `json:"valueSet"`
}

type TestMetric struct {
	Data *TestMetricData
}

type TestMetricsData struct {
	FirstRecord int64           `json:"firstRecord"`
	LastRecord  int64           `json:"lastRecord"`
	List        metricsListData `json:"list"`
}

type TestMetrics struct {
	LastRefreshTime int64
	Data            *TestMetricsData `json:"data"`
}
