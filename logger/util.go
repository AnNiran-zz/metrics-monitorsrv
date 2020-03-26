package logger

import(
	"fmt"
	"os"
	"path/filepath"
)

// accessLogFile is used to access /logs/metrics_data.json file and check if it exists
func accessLogFile() (bool, error) {
	// check if log file exists
	workPath, err := os.Getwd()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(filepath.Join(workPath, PackageName, LogsLocation, LogFilename)); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

// getLogsDir checks if /logs directory exists and if not - creates it
func getLogsDir() error {
	// check if /logs directory exists, if not - create it
	workPath, err := os.Getwd()
	if err != nil {
		return err
	}

	logsPath := filepath.Join(workPath, PackageName, LogsLocation)
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		fmt.Println("/logs directory does not exist. Creating ...")
		if err = os.MkdirAll(logsPath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
