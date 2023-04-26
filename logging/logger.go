package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const ENV_KEY_LOG_FILE string = "IIB_CLI_LOG_FILE"
const DEFAULT_LOG_FILE_NAME string = "iib-cli.log"

var loggerMap map[string]*log.Logger

func getLogFile() string {
	logFile := os.Getenv(ENV_KEY_LOG_FILE)
	if logFile == "" {
		tmpDir := os.TempDir()
		logFile = filepath.Join(tmpDir, DEFAULT_LOG_FILE_NAME)
	}
	return logFile
}

func newFileLogger(level string, logFile string) *log.Logger {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	prefix := fmt.Sprintf("%s: ", level)
	return log.New(file, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}

func getLogger(level string) *log.Logger {
	if loggerMap == nil {
		loggerMap = make(map[string]*log.Logger)
	}

	logger := loggerMap[level]
	if logger == nil {
		logger = newFileLogger(level, getLogFile())
		loggerMap[level] = logger
	}

	return logger
}

func DEBUG() *log.Logger {
	return getLogger("DEBUG")
}

func INFO() *log.Logger {
	return getLogger("INFO")
}

func WARNING() *log.Logger {
	return getLogger("WARNING")
}

func ERROR() *log.Logger {
	return getLogger("ERROR")
}
