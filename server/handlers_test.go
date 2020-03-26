package server

import (
	"testing"
	"net/http"

	"monitor/counter"
)

func TestGetMetricSum(t *testing.T) {
	Run()
	

	// run the handler
	req := &http.Request{}
	w := http.ResponseWriter

	handler(w), req)
}
