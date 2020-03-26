package server

import (
	"fmt"
	"sync"
	"net/http"
	"monitor/counter"
)

// port sets the URL and port for receiving requests to the standard 127.0.0.1:8080
var port = 8080
var addr = "localhost:8080"

// HttpSrv object wraps around a net/http server and keeps additional 
// data, as well as a sync.Mutex
type HttpSrv struct {
	port      int

	server    *http.Server
	isStarted bool
	mtx       *sync.Mutex
}

// requests URL data types that are used for json marshalling
// for making and receiving requests
type (
	MetricValuePost map[string]int

	MetricSumGet   map[string]int
	MetricsListGet map[string]map[string]int
)

// Errors
var (
	ErrImportMetrics = func(err string) string {
		return fmt.Sprintf("Error importing metrics data: %s", err)
	}
	ErrExportMetrics = func(err string) string {
		return fmt.Sprintf("Error exporting metrics data: %s", err)
	}

	ErrMetricCreate = func(err string) string {
		return fmt.Sprintf("Error creating new metric: '%s'", err)
	}
	ErrMetricAddValue = func(err string) string {
		return fmt.Sprintf("Error adding metric value: '%s'", err)
	}
	ErrMetricUpdate = func(err string) string {
		return fmt.Sprintf("Error updating metric name: '%s'", err)
	}
	ErrGetMetricSum = func(err string) string {
		return fmt.Sprintf("Error obtaining metric sum: %s", err)
	}
)

// Oputput success
var (
	OutputSuccess = func (action string) string {
		return fmt.Sprintf("Operation successful: %s", action)
	}
)

// Metrics objects instance used across handlers
var metrics *counter.Metrics
