package server

import (
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/argcv/stork/log"

	"monitor/counter"
)


// Default handler is Info
func info(w http.ResponseWriter, req *http.Request) {
	fmt.Println("\n")
	fmt.Println("*****")

	fmt.Printf("Server running at: %s\n", addr)
	fmt.Println("To create a new metric use: \n /create/<metric-name>/\n")
	fmt.Println("To set a new value for a metric use: \n /metric/<metric-name>/\n")
	fmt.Println("To update a metric name use: \n /metric/<metric-name>/update/\n")

	fmt.Println("To request a metric cumulative value for the last hour use:\n /metric/<metric-name>/sum\n")
	fmt.Println("To request a list of all existing metrics and their values for the last hour use:\n /metrics/list\n")
}

// postCreateMetric handles requests for creating a new metric
// POST: /create/<metric-name>
func postCreateMetric(manager counter.MetricsCounterSetterInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars  := mux.Vars(req)
		mname := vars["key"]
		
		if err := manager.AddMetric(mname); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(ErrMetricCreate(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(OutputSuccess(fmt.Sprintf("Metric %s created", mname)))
	}
}

// postAddMetricValue handles requests for setting a new value for a metric
// POST: /metric/<metric-name>/ {"value" : int }
func postAddMetricValue(manager counter.MetricsCounterSetterInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		// Check query input for `value` value
		if len(vars["value"]) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(ErrMetricAddValue("Value is not provided"))
			return
		}

		value, err := strconv.Atoi(vars["value"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(ErrMetricAddValue(err.Error()))
			return
		}

		// Value must be greater than 0
		if value == 0 {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(ErrMetricAddValue("Value must be a sufficient number greater than 0"))
			return
		}

		// Input is sufficient -> proceed
		if err := manager.AddMetricValue(vars["key"], value); err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Error(ErrMetricAddValue(err.Error()))
			return
		}
	
		w.WriteHeader(http.StatusOK)
		log.Info(OutputSuccess(fmt.Sprintf("Added value %s to metric %s", vars["value"], vars["key"])))
	}
}

// postAddMetricValueNotProvided handles missing query in post request for adding a value
func postAddMetricValueNotProvided() http.HandlerFunc {
	return func(w http.ResponseWriter, red *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(ErrMetricAddValue("Value is not provided"))
	}
}

// postUpdateMetric handles requests for updating a metric name
// POST: /metric/<metric-name>/update
func postUpdateMetric(manager counter.MetricsCounterSetterInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		// Check query input for `value` value
		if len(vars["value"]) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(ErrMetricUpdate("Update value is not provided"))
			return
		}

		// Input is sufficient -> proceed
		if err := manager.UpdateMetricName(vars["key"], vars["value"]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error(ErrMetricUpdate(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(OutputSuccess(fmt.Sprintf("Metric name %s updated to %s", vars["key"], vars["value"])))
	}
}

// postUpdateMetricValueNotProvided handles missing query in post request for updating a metric
func postUpdateMetricValueNotProvided() http.HandlerFunc {
	return func(w http.ResponseWriter, red *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		log.Info(ErrMetricUpdate("Update value is not provided"))
	}
}

// getMetricSum handles requests for getting a metric value sum
// GET: /metric/{key}/sum
func getMetricSum(manager counter.MetricsCounterGetterInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		sum, err := manager.GetMetric(vars["key"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info(ErrGetMetricSum(err.Error()))
			return
		}

		payload := MetricSumGet{vars["key"]:sum}
		json.NewEncoder(w).Encode(payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		log.Info(OutputSuccess(fmt.Sprintf("Metric %s sum: %s", vars["key"], sum)))
	}
}

// getMetricsList handles requests for listing all existing metrics
// GET: /metrics/list
func getMetricsList(manager counter.MetricsCounterGetterInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list    := manager.ListMetrics()
		payload := MetricsListGet{"metrics":list}

		json.NewEncoder(w).Encode(payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		log.Info(OutputSuccess("metric listed"))
	}
}
