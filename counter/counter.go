package counter

import (
	"time"
	"sync"
)

/// Functions and methods for Metric object manipulation

// createMetric creates a new metric object without setting values and time
// Values and recoding time will be set in consequent methods
func createMetric(name string) *Metric {
	return &Metric{
		Lock: &sync.Mutex{},
		Data: &MetricData {
			Name:     name,
			ValueSet: make(timeValues),
		},
	}
}

// # Metric object implements metricManagerInterface

// addValue adds a new value to a metric
func (m *Metric) addValue(v int, time int64) error {
	// Same times can occure if more than one values from different routines
	// are requested within one second
	if _, ok := m.Data.ValueSet[time]; ok {
		return ErrSameSecondRecord(m.Data.Name)
	}

	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Data.LastRecord = time
	if m.Data.FirstRecord == 0 { // int64 zero value is 0
		m.Data.FirstRecord = time
	}
	m.Data.Sum += v
	m.Data.ValueSet[time] = v
	return nil
}

func (m *Metric) removeValue(t int64) {
	if _, ok := m.Data.ValueSet[t]; !ok {
		return
	}

	delete(m.Data.ValueSet, t)
}

func (m *Metric) update() {}

/// Functions and methods for Metric object manipulation

// loadMetricsObject initializes MetricData manager object
func loadMetricsObject() *Metrics {
	return &Metrics{
		Lock:            &sync.Mutex{},
		SignalRefresher: make(chan bool),
		
		Data: &MetricsData{
			List: make(metricsListData),
		},
	}
}

// # Metrics object implements metricsCounterInterface

// AddMetric adds a new metric record to the list of metrics
func (ms *Metrics) AddMetric(name string) error {
	if _, ok := ms.Data.List[name]; ok {
		return ErrMetricExists(name)
	}

	// if the metric does not exist -> create it
	m := createMetric(name)

	ms.Lock.Lock()
	defer ms.Lock.Unlock()

	// we do not update LastRecord value here because there is no value set
	// creating a new metric object does not sommunicate with the time range condition
	// of keeping records
	ms.Data.List[name] = m
	return nil
}

// GetMetric returns a metric
func (ms *Metrics) GetMetric(name string) (int, error) {
	if _, ok := ms.Data.List[name]; !ok {
		return 0, ErrMetricNotExisting(name)
	}

	return ms.Data.List[name].Data.Sum, nil
}

// AddMetricValue add a new value to the metric object
func (ms *Metrics) AddMetricValue(name string, value int) error {
	if _, ok := ms.Data.List[name]; !ok {
		return ErrMetricNotExisting(name)
	}

	time := time.Now().Unix()

	ms.Lock.Lock()
	ms.Data.LastRecord = time
	if ms.Data.FirstRecord == 0 { // int64 zero value is 0
		ms.Data.FirstRecord = time
	}
	ms.Lock.Unlock()

	return ms.Data.List[name].addValue(value, time)
}

// UpdateMetricName updates a name for a specific metric
func (ms *Metrics) UpdateMetricName(name, update string) error {
	if _, ok := ms.Data.List[name]; !ok {
		return ErrMetricNotExisting(name)
	}

	ms.Lock.Lock()
	defer ms.Lock.Unlock()
	ms.Data.List[update] = ms.Data.List[name]
	delete(ms.Data.List, name)
	return nil
}

// ListMetrics returns a list of all metrics and their sums
func (ms *Metrics) ListMetrics() metricsList {
	list := make(metricsList)

	for mname, data := range ms.Data.List {
		list[mname] = data.Data.Sum
	}

	return list
}

// removeMetric removes a metric from existing list of metrics
// next time current record are logged - the metric name will not be present
func (ms *Metrics) removeMetric(name string) error {
	if _, ok := ms.Data.List[name]; !ok {
		return ErrMetricNotExisting(name)
	}

	ms.Lock.Lock()
	defer ms.Lock.Unlock()

	delete(ms.Data.List, name)
	return nil
}
