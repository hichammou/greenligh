package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

// We use the iota keyword as a shortcut to assign successive integer values to the constants
const (
	LevelInfo  Level = iota // Has the value 0,
	LevelError              // Has the value 1.
	LevelFatal              // Has the value 2.
	LevelOff                // Has the value 3.
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Define a custom logger. This holds the output destination that the log entries will be written to
// the minimum severity level that log entries will be written for, and a mutex for coordinating the writes
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// Return new Logger instance.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties"`
		Trace      string            `json:"trace"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	// Show error stack if the log is an error
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	// Declare a line variable to hold the actual log
	var line []byte

	// Marcha the annonymous struct to JSON and store it in the line variable
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marchal log message: " + err.Error())
	}

	// Lock the mutex so that no two writes to the output destination can happen concurrently.
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out.Write(append(line, '\n'))
}

// We implement a Write() method so our logger satisfies the io.Writer interface. This writes a log entry at the ERROR level
// with no additional properties
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
