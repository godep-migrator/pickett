package logit

import (
	"time"
)

// These are the integer logging levels used by the logger
type Level int

const (
	FLUSH = -1
	CLOSE = -2
)

const (
	FINEST Level = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
	NONE
)

// Logging level strings
var (
	levelStrings = [...]string{"FNST", "FINE", "DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT", "NONE"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

// This is an interface for anything that should be able to write logs
type LogWriter interface {
	StartLogger(chan *LogRecord)
	Closed() bool
}

// A LogRecord contains all of the pertinent information for each message
type LogRecord struct {
	Level      Level      // The log level
	Created    time.Time  // The time at which the log message was created (nanoseconds)
	Source     string     // The message source
	NestedName string     // The name of the nested logger
	Message    string     // The log message
	GoRoutine  int32      // An identifier for this goroutine (if possible)
	reply      chan error //If this is a "special" flush message a channel for announcing completion
}

func (rec *LogRecord) GetReplyChan() chan error {
	return rec.reply
}

func (rec *LogRecord) SetReplyChan(c chan error) {
	rec.reply = c
}
