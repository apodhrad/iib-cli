package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setTestLogFile(t *testing.T) string {
	logFile := filepath.Join(t.TempDir(), "test.log")
	os.Setenv(ENV_KEY_LOG_FILE, logFile)
	return logFile
}

func unsetTestLogFile(t *testing.T) {
	os.Unsetenv(ENV_KEY_LOG_FILE)
}

func TestGetDefaultLogFile(t *testing.T) {
	expectedLogFile := "/tmp/" + DEFAULT_LOG_FILE_NAME
	logFile := getLogFile()
	assert.Equal(t, expectedLogFile, logFile)
}

func TestGetCustomLogFile(t *testing.T) {
	testLogFile := setTestLogFile(t)
	defer unsetTestLogFile(t)

	logFile := getLogFile()
	assert.Equal(t, testLogFile, logFile)

}

func testLevelLog(t *testing.T, level string, loggerFunc func() *log.Logger) {
	testLogFile := setTestLogFile(t)
	defer unsetTestLogFile(t)

	testData := fmt.Sprintf("Some %s info", level)
	loggerFunc().Println(testData)

	data, err := os.ReadFile(testLogFile)
	assert.Nil(t, err)
	logData := string(data)
	regex := fmt.Sprintf("%s: [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} logger_test.go:[0-9]+: %s\n", level, testData)
	assert.Regexp(t, regex, logData)
}
func TestDebugLog(t *testing.T) {
	testLevelLog(t, "DEBUG", DEBUG)
}

func TestInfoLog(t *testing.T) {
	testLevelLog(t, "INFO", INFO)
}

func TestWarningLog(t *testing.T) {
	testLevelLog(t, "WARNING", WARNING)
}

func TestErrorLog(t *testing.T) {
	testLevelLog(t, "ERROR", ERROR)
}
