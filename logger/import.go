package logger

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/argcv/stork/log"
)

// importMetricsData import data from log file
func importMetricsData() ([]byte, error) {
	// check if log file exists, if it does not exist, we return empty data object
	logExists, err := accessLogFile()
	if err != nil {
		return nil, err
	}

	workPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// if log file does not exist - we do not have anything to import
	logFile := filepath.Join(workPath, PackageName, LogsLocation, LogFilename)
	if !logExists {
		log.Infof(OutputLogsNotExist(logFile))
		return nil, nil
	}

	// log file exists -> open it
	jsonSrc, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer jsonSrc.Close()

	metrics, err := ioutil.ReadAll(jsonSrc)
	if err != nil {
		return nil, err
	}
	log.Info(OuputMetricsImport())
	return metrics, nil
}
