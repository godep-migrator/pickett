// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package logit

import (
	"time"
)

var (
	Global *Logger
)

func init() {
	Global = NewDefaultLogger(INFO)
}

// Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl Level, expensive *Expensive, nd *NestedData, writer LogWriter) error {
	return Global.AddFilter(name, lvl, expensive, nd, writer)
}

// Wrapper for (*Logger).AddFilter
func ModifyFilterLvl(name string, lvl Level, expensive *Expensive, nd *NestedData) error {
	return Global.ModifyFilterLvl(name, lvl, expensive, nd)
}

// Logs the given message and crashes the program
func Panicf(format string, args ...interface{}) {
	Global.Panicf(format, args...)
}

// Logs the given message and crashes the program
func Panicln(args ...interface{}) {
	Global.Panicln(args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl Level, source, message string) {
	Global.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl Level, format string, args ...interface{}) {
	Global.intLogf(lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl Level, closure func() string) {
	Global.intLogc(lvl, closure)
}

// Send a argument list log message
// Wrapper for (*Logger).Logln
func Logln(lvl Level, args ...interface{}) {
	Global.intLogln(lvl, args)
}

func Flush(d time.Duration) error {
	return Global.Flush(d)
}

func Close(d time.Duration) error {
	return Global.Close(d)
}
