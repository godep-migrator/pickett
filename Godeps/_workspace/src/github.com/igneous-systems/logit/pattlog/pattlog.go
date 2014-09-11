// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package pattlog

import (
	"fmt"
	"io"
	. "logit"
)

const (
	FORMAT_DEFAULT = "[%D %T] [%L] (%S) %M"
	FORMAT_SHORT   = "[%t %d] [%L] %M"
	FORMAT_ABBREV  = "[%L] %M"
)

// This is the standard writer that prints to standard output.
type FormatLogWriter struct {
	skelLog *SkelLog
	out     io.Writer
	format  string
}

// This creates a new FormatLogWriter
func NewFormatLogWriter(out io.Writer, format string) *FormatLogWriter {
	var fw FormatLogWriter
	fw.out = out
	fw.format = format
	fw.skelLog = NewSkelLog(&fw, nil)
	return &fw
}

func (fw *FormatLogWriter) Flush() error {
	return nil
}

func (fw *FormatLogWriter) Close() error {
	return nil
}

func (fw *FormatLogWriter) Cleanup() {
}

func (fw *FormatLogWriter) Wake(obj interface{}) error {
	return nil
}

func (fw *FormatLogWriter) LogNormal(rec *LogRecord) error {
	fmt.Fprint(fw.out, FormatLogRecord(fw.format, rec))
	return nil
}

func (fw *FormatLogWriter) StartLogger(c chan *LogRecord) {
	fw.skelLog.StartLogger(c)
}

func (fw *FormatLogWriter) Closed() bool {
	return fw.skelLog.Closed()
}
