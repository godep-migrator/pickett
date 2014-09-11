// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.
// Copyright (C) 2014, Igneous Systems.  All rights reserved.

// Package logit provides level-based and highly configurable logging, with an optional
// hieracrchy of loggers for fine grained log level control.  It was derived
// from the log4go package but is now substantially different from it.
//
// Enhanced Logging
//
// This is inspired by the logging functionality in Java.  Essentially, you create a Logger
// object and create output filters for it.  You can send whatever you want to the Logger,
// and it will filter that based on your settings and send it to the outputs.  This way, you
// can put as much debug code in your program as you want, and when you're done you can filter
// out the mundane messages so only the important ones show up.
//
//
// Logging functions
//
// Logit provides three forms of most of the logging functions.  There is a *f form that takes
// a Printf style argument list.  There is an *ln form that takes a series of objects and formats
// using Println style formating.  And finally there is a *c format that takes a function closure
// that returns a string.  This function is evaluated at most once (and only if the log will be used).
// This last form is intended as a way of avoiding computation of expensive debugging output unless
// it will actually be used.
//	logit.Infof("A Printf style log of the number %d", 10)
//	logit.Infoln("A Println", "style log of the number", 10)
//	logit.Infoc(func() string { return fmt.Sprintf("A closure style log of the number %d", 10) })
//
// Logging Levels
//
// Logit provides the following logger levels (in order).  These can be directly specified to the Log*
// functions but are more commonly implicitly specified via a Debugf(), Infof(), etc function.  There
// are no logger functions for the NONE level and it is special because if it is set as a log level
// no logging will occur (even for log requests explicitly set to that level)
//
//	FINEST Level = iota
//	FINE
//	DEBUG
//	TRACE
//	INFO
//	WARNING
//	ERROR
//	CRITICAL
//	NONE
//
// Sample Logger
//
// Simple use of local, non nested loggers.
//
// 	log := logit.NewLogger()
// 	log.AddFilter("stdout", logit.DEBUG, nil, nil, logit.NewConsoleLogWriter())
// 	log.AddFilter("log",    logit.FINE, nil, nil, logit.NewFileLogWriter("example.log", true))
// 	log.Infof("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
// 	The first two lines can be combined with the utility NewDefaultLogger:
// 	log := logit.NewDefaultLogger(logit.DEBUG)
// 	log.AddFilter("log",    logit.FINE, nil, nil, logit.NewFileLogWriter("example.log", true))
// 	log.Infof("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// Nested Loggers
//
// Simple use of nested local loggers.  These loggers expose a package hierarchy that can be
// used to modify logging levels for arbitrary portions of the hierarchy.  All settings for
// the logger must still occur on the underlying non-tested logger wrapped by the nested
// loggers.
//
//	main := logit.NewDefaultLogger(logit.INFO)
//	nested := logit.NewNestedLogger("github.com/this/is/my/project", main)
//	nested2 := logit.NewNestedLogger("github.com/this/is", main)
//	nd := NewNestedData()
//	nd.AddEntry("github.com/this/is/my", logit.DEBUG, nil)
//	main.ModifyFilterLvl("stdout", logit.INFO, nil, nd)
//	main.Debugf("This doesn't show up")
//	nested2.Debugf("Nor does this")
//	nested.Debugf("But this does")
//
// It is also possible to have the system infer the package name based on runtime
// caller information.  In this case, the runtime package of the immediate caller
// is used as the name.
//
//	nested := logit.NewNestedLoggerFromCaller(main)
package logit

import (
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

/****** Constants ******/

/****** Variables ******/
var (
	// logBufferLength specifies how many log messages a particular logit
	// logger can buffer at a time before writing them.
	logBufferLength = 32
)

/****** Logger ******/

type Expensive struct {
	srcRe []*regexp.Regexp
}

func NewExpensive() *Expensive {
	var exp Expensive
	return &exp
}

func (expensive *Expensive) AddSrcRegex(re *regexp.Regexp) {
	expensive.srcRe = append(expensive.srcRe, re)
}

// A Filter represents the log level below which no log records are written to
// the associated LogWriter.
type Filter struct {
	writer LogWriter
	c      chan *LogRecord
}

// A named logger entry.  Names are used for replacing existing
// filters with new settings.
type filterEntry struct {
	name      string
	lvl       Level
	expensive *Expensive
	nd        *NestedData
	filter    *Filter
}

type filterList []filterEntry

// A Logger represents a collection of Filters through which log messages are
// written.
type Logger struct {
	// The name for the nested logger
	nestedName string
	// the version number of the settings
	version int64
	// the master logger if this is nested
	master *Logger
	// The current list of filter entries
	fList filterList
	// Grabbed when reading or setting the filter list.  Shared with the master.
	swapMutex *sync.Mutex
	// Grabbed over the entire process of modifying the filter list
	modMutex sync.Mutex
}

// Return a snapshot of the filter list.  This slice directly references the
// same array as the internal list, so it should not be modified.
func (logger *Logger) getFilterRef(force bool) filterList {
	if logger.master == nil {
		logger.swapMutex.Lock()
		defer logger.swapMutex.Unlock()
		return logger.fList
	} else {
		logger.swapMutex.Lock()
		if !force && logger.master.version == logger.version {
			defer logger.swapMutex.Unlock()
			return logger.fList
		}
		mustCopy := false
		for i := range logger.master.fList {
			if logger.master.fList[i].nd != nil {
				mustCopy = true
				break
			}
		}
		if !mustCopy {
			defer logger.swapMutex.Unlock()
			logger.fList = logger.master.fList
			logger.version = logger.master.version
			return logger.fList
		}
		list := logger.master.fList
		version := logger.master.version
		logger.swapMutex.Unlock()

		copyList := make(filterList, len(list))
		copy(copyList, list)

		for i := range copyList {
			if copyList[i].nd != nil {
				ne := copyList[i].nd.closestEntry(logger.nestedName)
				copyList[i].nd = nil
				if ne != nil {
					copyList[i].lvl = ne.lvl
					copyList[i].expensive = ne.expensive
				}
			}
		}

		logger.swapMutex.Lock()
		if logger.master.version == version {
			logger.fList = copyList
			logger.version = version
		}
		logger.swapMutex.Unlock()
		return copyList
	}

}

// Return a snapshot of the filter list.  This slice is a full slice copy
// of the the array in internal list, so the slice can be modified.  Referenced
// data structures are NOT copied, so should NOT be modified.
func (logger *Logger) getFilterCopy() filterList {
	refList := logger.getFilterRef(false)
	copyList := make(filterList, len(refList), len(refList)+1)
	copy(copyList, refList)
	return copyList
}

// Set a new filter list.  It is safe to call this function concurrently.
func (logger *Logger) setFilterRef(newList filterList) {
	logger.swapMutex.Lock()
	defer logger.swapMutex.Unlock()
	logger.fList = newList
	logger.version++
}

func NewLogger() *Logger {
	var logger Logger
	logger.swapMutex = new(sync.Mutex)
	return &logger
}

func (log *Logger) GetName() string {
	return log.nestedName
}

// Create a nested logger with the specified name. The name should be
// a golang style package.
func NewNestedLogger(name string, parent *Logger) *Logger {
	master := parent
	if parent.master != nil {
		master = parent.master
	}

	logger := new(Logger)
	logger.nestedName = name
	logger.swapMutex = parent.swapMutex
	logger.master = master
	logger.getFilterRef(true)
	return logger
}

// Create a nested logger, implicitly naming based on the package name of the
// immediate calling function.
func NewNestedLoggerFromCaller(parent *Logger) *Logger {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	name := f.Name()
	pos := strings.Index(name, "(")
	if pos != -1 {
		name = name[:pos]
	}
	pos = strings.LastIndex(name, ".")
	if 0 >= pos {
		parent.Errorf("Unable to find package name, skipping nested logger")
		return parent
	} else {
		name = name[:pos]
		return NewNestedLogger(name, parent)
	}
}

// Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
func NewDefaultLoggerWithTarget(lvl Level, target io.Writer) *Logger {
	logger := NewLogger()
	err := logger.AddFilter("stdout", lvl, nil, nil, NewConsoleLogWriterWithTarget(target))
	if err != nil {
		panic("Error creating default logger")
	}
	return logger
}

// Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
func NewDefaultLogger(lvl Level) *Logger {
	logger := NewLogger()
	err := logger.AddFilter("stdout", lvl, nil, nil, NewConsoleLogWriter())
	if err != nil {
		panic("Error creating default logger")
	}
	return logger
}

func sendToFilter(filter *Filter, rec *LogRecord) (success bool) {
	defer func() {
		r := recover()
		if r != nil {
			if filter.writer.Closed() {
				success = false
			} else {
				panic(r)
			}
		}
	}()
	filter.c <- rec
	return true
}

// Add a new LogWriter to the Logger which will only log messages at lvl or
// higher. This will refuse to replace existing filters.  It is safe to call
// this from go routines.
func (log *Logger) AddFilter(name string, lvl Level, expensive *Expensive, nd *NestedData, writer LogWriter) error {
	if nd != nil {
		nd.lockEntries()
	}

	c := make(chan *LogRecord, logBufferLength)
	writer.StartLogger(c)
	filter := &Filter{writer, c}

	log.modMutex.Lock()
	defer log.modMutex.Unlock()
	if log.master != nil {
		return fmt.Errorf("Attempt to modify nested filter settings!")
	}

	fe := filterEntry{name, lvl, expensive, nd, filter}
	list := log.getFilterCopy()
	for i := range list {
		if list[i].name == name {
			return fmt.Errorf("Tried to replace existing loger %s", name)
		}
	}

	list = append(list, fe)
	log.setFilterRef(list)
	return nil
}

// Modify the level of the named filter while maintaining the specified lvl.
// Returns true on success, false on failure.
func (log *Logger) ModifyFilterLvl(name string, lvl Level, expensive *Expensive, nd *NestedData) error {
	if nd != nil {
		nd.lockEntries()
	}

	log.modMutex.Lock()
	defer log.modMutex.Unlock()
	if log.master != nil {
		return fmt.Errorf("Attempt to modify nested filter settings!")
	}
	list := log.getFilterCopy()
	for i := range list {
		if list[i].name == name {
			list[i].lvl = lvl
			list[i].expensive = expensive
			list[i].nd = nd
			log.setFilterRef(list)
			return nil
		}
	}
	return fmt.Errorf("Couldn't find filter name %v", name)
}

func (exp *Expensive) matchExpensive(str string) bool {
	if exp == nil {
		return false
	}
	for _, re := range exp.srcRe {
		if re.MatchString(str) {
			return true
		}
	}
	return false
}

// Determine if any logging will be done
func (list filterList) skipLog(lvl Level, haveSrc bool, srcIn string) (string, bool) {
	hadExpensive := false
	// Determine if any logging will be done
	for _, fe := range list {
		if fe.lvl != NONE && lvl >= fe.lvl {
			if !haveSrc {
				return callerToSourceString(), false
			} else {
				return srcIn, false
			}
		}
		if fe.expensive != nil {
			hadExpensive = true
		}
	}
	if !hadExpensive {
		return "", true
	}
	if !haveSrc {
		srcIn = callerToSourceString()
	}
	for _, fe := range list {
		if fe.expensive.matchExpensive(srcIn) {
			return srcIn, false
		}
	}
	return "", true
}

func (list filterList) dispatchRecord(lvl Level, nestedName string, src string, msg string) {
	// Make the log record
	rec := &LogRecord{
		Level:      lvl,
		Created:    time.Now(),
		GoRoutine:  GoID(),
		Source:     src,
		NestedName: nestedName,
		Message:    msg,
	}

	// Dispatch the logs
	for _, fe := range list {
		if lvl < fe.lvl || fe.lvl == NONE {
			if fe.expensive == nil {
				continue
			}
			if !fe.expensive.matchExpensive(src) {
				continue
			}
		}
		sendToFilter(fe.filter, rec)
	}
}

/******* Logging *******/
// Send a formatted log message internally
func (log *Logger) intLogf(lvl Level, format string, args ...interface{}) {
	list := log.getFilterRef(false)

	// Determine if any logging will be done
	src, skip := list.skipLog(lvl, false, "")
	if skip {
		return
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Make the record and dispatch to loggers
	list.dispatchRecord(lvl, log.nestedName, src, msg)
}

func callerToSourceString() string {
	_, fname, line, ok := runtime.Caller(4)
	if ok {
		return fmt.Sprintf("%s:%d", fname, line)
	} else {
		return "UNKNOWN"
	}
}

/******* Logging *******/
// Send a parameter log message internally
func (log *Logger) intLogln(lvl Level, args ...interface{}) {
	list := log.getFilterRef(false)

	// Determine if any logging will be done
	src, skip := list.skipLog(lvl, false, "")
	if skip {
		return
	}

	msg := fmt.Sprintln(args...)

	// Make the record and dispatch to loggers
	list.dispatchRecord(lvl, log.nestedName, src, msg)
}

// Send a closure log message internally
func (log *Logger) intLogc(lvl Level, closure func() string) {
	list := log.getFilterRef(false)

	// Determine if any logging will be done
	src, skip := list.skipLog(lvl, false, "")
	if skip {
		return
	}

	msg := closure()

	// Make the record and dispatch to loggers
	list.dispatchRecord(lvl, log.nestedName, src, msg)
}

// Send a log message with manual level, source, and message.
func (log *Logger) Log(lvl Level, src, msg string) {
	list := log.getFilterRef(false)

	// Determine if any logging will be done
	_, skip := list.skipLog(lvl, true, src)
	if skip {
		return
	}

	// Make the record and dispatch to loggers
	list.dispatchRecord(lvl, log.nestedName, src, msg)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func (log *Logger) Logf(lvl Level, format string, args ...interface{}) {
	log.intLogf(lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func (log *Logger) Logc(lvl Level, closure func() string) {
	log.intLogc(lvl, closure)
}

// Logs the given message and crashes the program
// Waits for up to 10 seconds for the logger to close before giving up and finishing
// panicing.
func (log *Logger) Panicf(format string, args ...interface{}) {
	log.intLogf(CRITICAL, format, args...)
	log.Close(10 * time.Second) // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Logs the given message and crashes the program.
// Waits for up to 10 seconds for the logger to close before giving up and finishing
// panicing.
func (log *Logger) Panicln(args ...interface{}) {
	log.intLogln(CRITICAL, args...)
	log.Close(10 * time.Second) // so that hopefully the messages get logged
	panic(fmt.Sprintln(args...))
}

// Send a fake loke expecting a reply (waiting for that reply)
// Currently this is just Flush and Close requests
func sendReplyExpected(list filterList, lev Level, d time.Duration) error {
	var rec LogRecord
	var err error
	c := make(chan error)
	sent := 0

	for _, fe := range list {
		rec.SetReplyChan(c)
		rec.Level = lev
		success := sendToFilter(fe.filter, &rec)
		if success {
			sent++
		}
	}
	// Treat forever as 10 years...
	if d < 0 {
		d = 10 * 365 * 24 * time.Hour
	}

	reply := 0
	start := time.Now()
	elapsed := time.Duration(0)

	for d != 0 && reply < sent {
		select {
		case err2 := <-c:
			if err == nil {
				err = err2
			}
			reply++
		case <-time.After(d):
		}
		elapsed = time.Since(start)
		if elapsed >= d {
			break
		} else {
			d -= elapsed
		}
	}

	if err == nil && reply < sent {
		err = fmt.Errorf("Didn't receive reply from all loggers")
	}

	return err
}

//Flush all previously sent messages to the log, returning when this is
//complete (or duration expires).  Returns with an error if any errors are
//received or if the duration expires.  Note that the timeout currently only
//applies to waiting for completion.  Blocking close is called on each channel.
//Negative duration blocks until the first error is received or the flush is
//complete.  This function does not wait for writes that occur after the call
//is made.
func (log *Logger) Flush(d time.Duration) error {
	list := log.getFilterRef(false)
	return sendReplyExpected(list, FLUSH, d)
}

//Flush all previously sent messages to the log, and close all loggers such that they throw away any
//more log messages ent to them.  Returns when this is complete (or duration
//expires).  Returns with an error if any errors are received or if the
//duration expires.  Note that the timeout currently only applies to waiting
//for completion.  Blocking close is called on each channel.  Negative duration
//blocks until the first error is received or the flush is complete.  This
//function does not wait for writes that occur after the call is made.
func (log *Logger) Close(d time.Duration) error {
	log.modMutex.Lock()
	if log.master != nil {
		log.modMutex.Unlock()
		return fmt.Errorf("Attempt to close nested filter!")
	}

	list := log.getFilterRef(false)
	log.setFilterRef(nil)
	log.modMutex.Unlock()
	return sendReplyExpected(list, CLOSE, d)
}
