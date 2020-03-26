package logger

import (
	"fmt"
)

// TimeRange is used for comparing times in Unixtime format
var TimeRange = 60*60

// Outputs and errors
var (
	OutputLogsNotExist = func(filepath string) string {
		return fmt.Sprintf("Log file does not exist inside %s - no metrics data to import \nInitializing empty metrics records", filepath)
	}
	OutputEmptyMetrics = func() string {
		return "No metrics data to export"
	}
	OuputMetricsImport = func() string {
		return "Imported metrics"
	}
	
	ErrDataExport = func(err string) error {
		return fmt.Errorf("Could not export metrics data to json: %s", err)
	}
)

// Files and paths
var PackageName  = "logger"
var LogsLocation = "logs"

// LogFile has the same name for all exporting and importing now,
// because the functionality is for representative purpose
var LogFilename = "metrics_data.json"
