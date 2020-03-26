package logger

import (
	"os"
	"path/filepath"
)

// exportMetricsData exports metrics data to json file inside logs directory
func exportMetricsData(data []byte) error {
	if err := getLogsDir(); err != nil {
		return err
	}

	workPath, err := os.Getwd()
	if err != nil {
		return err
	}

	// we do not care to check if any file exists inside the logs directory
	// because exporting and importing is with representational purpose here
	logFile, err := os.OpenFile(filepath.Join(workPath, PackageName, LogsLocation, LogFilename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer logFile.Close()

	logFile.Write(data)
	return nil
}
