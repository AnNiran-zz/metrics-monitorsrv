package server

import (
	"fmt"
	"github.com/argcv/stork/log"
	"net/http"
	"time"
	"sync"
	"errors"
	"context"

	"github.com/gorilla/mux"
)


func newHttpSrv() *HttpSrv {
	return &HttpSrv{
		port:      8080,
		server:    nil,
		isStarted: false,
		mtx:       &sync.Mutex{},
	}
}

// Start initializes data at rest import, initializing counter variables
// and starting handling HTTP requests to the API
func (srv *HttpSrv) start() (err error) {
	srv.mtx.Lock()
	defer srv.mtx.Unlock()

	if srv.isStarted {
		return errors.New("Server is already running")
	}
	srv.isStarted = true

	initMetrics()

	// Initialize mux router
	muxr := mux.NewRouter()

	// Register routes
	registerRoutes(muxr)

	// Initialize net/http server
	srv.server = &http.Server{
		Handler :     muxr,
		Addr:         addr, 
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run server
	go func() {
		if err = srv.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				// server is closed upon request
				log.Infof("Server closed upon request: %v", err)

			} else {
				log.Fatalf("Server closed unexpect: %v", err)
			}
			srv.isStarted = false
		}
	}()

	time.Sleep(10 * time.Millisecond)
	return
}

func (m *HttpSrv) shutdown(ctx context.Context) (err error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if !m.isStarted || m.server == nil {
		return errors.New("Server is not started")
	}

	stop := make(chan bool)
	go func() {
		// Export metrics to file at graceful shutdown
		logMetrics()

		// Stop refresher go routine
		metrics.SignalRefresher <- true

		err = m.server.Shutdown(ctx)
		stop <- true
	}()

	select {
	case <-ctx.Done():
		log.Errorf("Timeout: %v", ctx.Err())
	case <-stop:
		log.Infof("Finished")
	}

	return
}

// registerRoutes registers the mux router handlers
func registerRoutes(muxr *mux.Router) {
	// Default handler route
	muxr.HandleFunc("/", info)

	// POST:
	// Create a new metric
	muxr.HandleFunc("/create/{key}", postCreateMetric(metrics)).Methods("POST")
	// Add metric value
	muxr.HandleFunc("/metric/{key}", postAddMetricValue(metrics)).Queries("value", "{value}").Methods("POST")
	muxr.HandleFunc("/metric/{key}", postAddMetricValueNotProvided()).Methods("POST")

	// Update metric name
	muxr.HandleFunc("/metric/{key}/update", postUpdateMetric(metrics)).Queries("value", "{value}").Methods("POST")
	muxr.HandleFunc("/metric/{key}/update", postUpdateMetricValueNotProvided()).Methods("POST")
	
	// GET:
	// Get metric sum
	muxr.HandleFunc("/metric/{key}/sum", getMetricSum(metrics)).Methods("GET")
	// Get list of metric
	muxr.HandleFunc("/metrics/list", getMetricsList(metrics)).Methods("GET")
}

// initMetrics starts the metrics monitor by initializing the metrics obejct
// and runs the refresher for it
func initMetrics() {
	// Initialize counter metrics
	if err := startMetricsMonitor(); err != nil {
		fmt.Println(ErrImportMetrics(err.Error()))
	}

	// Start metrics refresher
	metrics.RunRefresher()
}
