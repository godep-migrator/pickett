// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package logit

import (
	"fmt"
	"io"
	"os"
)

// This is the standard writer that prints to standard output.
type ConsoleLogWriter struct {
	timestr   string
	timestrAt int64

	skelLog *SkelLog
	out     io.Writer
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriterWithTarget(target io.Writer) *ConsoleLogWriter {
	var cw ConsoleLogWriter
	cw.out = target
	cw.skelLog = NewSkelLog(&cw, nil)
	return &cw
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() *ConsoleLogWriter {
	return NewConsoleLogWriterWithTarget(os.Stdout)
}

func (cw *ConsoleLogWriter) Flush() error {
	return nil
}

func (cw *ConsoleLogWriter) Close() error {
	return nil
}

func (cw *ConsoleLogWriter) Cleanup() {
}

func (cw *ConsoleLogWriter) Wake(obj interface{}) error {
	return nil
}

func (cw *ConsoleLogWriter) LogNormal(rec *LogRecord) error {
	if at := rec.Created.UnixNano() / 1e9; at != cw.timestrAt {
		cw.timestr, cw.timestrAt = rec.Created.Format("01/02/06 15:04:05"), at
	}
	fmt.Fprint(cw.out, "[", cw.timestr, "] [", rec.Level.String(), "] ", rec.Message, "\n")

	return nil
}

func (cw *ConsoleLogWriter) StartLogger(c chan *LogRecord) {
	cw.skelLog.StartLogger(c)
}

func (cw *ConsoleLogWriter) Closed() bool {
	return cw.skelLog.Closed()
}
