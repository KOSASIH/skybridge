package logging

import (
	"bytes"
	"testing"
)

func TestNewLogger(t *testing.T) {
	output := &bytes.Buffer{}
	logger := NewLogger(INFO, output)
	if logger.level != INFO {
		t.Errorf("Expected logger level to be INFO, got %d", logger.level)
	}
	if logger.output != output {
		t.Errorf("Expected logger output to be %v, got %v", output, logger.output)
	}
}

func TestLog(t *testing.T) {
	output := &bytes.Buffer{}
	logger := NewLogger(DEBUG, output)
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warning("Warning message")
	logger.Error("Error message")
	logger.Fatal("Fatal message")
	expected := `[DEBUG] Debug message
[INFO] Info message
[WARNING] Warning message
[ERROR] Error message
[FATAL] Fatal message
`
	if output.String() != expected {
		t.Errorf("Expected log output to be %s, got %s", expected, output.String())
	}
}

func TestLogLevel(t *testing.T) {
	output := &bytes.Buffer{}
	logger := NewLogger(INFO, output)
	logger.Log(DEBUG, "Debug message")
	if output.String() != "" {
		t.Errorf("Expected log output to be empty, got %s", output.String())
	}
}
