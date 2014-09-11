// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package logit

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
	"time"
)

var now time.Time = time.Unix(0, 1234567890123456789).In(time.UTC)

func newLogRecord(lvl Level, src string, msg string) *LogRecord {
	return &LogRecord{
		Level:   lvl,
		Source:  src,
		Created: now,
		Message: msg,
	}
}

func helper_flush(c chan *LogRecord) error {
	var rec LogRecord
	rec.Level = FLUSH
	rec.SetReplyChan(make(chan error))
	c <- &rec
	err := <-rec.GetReplyChan()
	return err
}

func helper_closeit(c chan *LogRecord) error {
	var rec LogRecord
	rec.Level = CLOSE
	rec.SetReplyChan(make(chan error))
	c <- &rec
	err := <-rec.GetReplyChan()
	return err
}

func helper_closeitHideErr(c chan *LogRecord) error {
	defer func() {
		_ = recover()
	}()
	var rec LogRecord
	rec.Level = CLOSE
	rec.SetReplyChan(make(chan error))
	c <- &rec
	err := <-rec.GetReplyChan()
	return err
}

func TestELog(t *testing.T) {
	lr := newLogRecord(CRITICAL, "source", "message")
	if lr.Level != CRITICAL {
		t.Errorf("Incorrect level: %d should be %d", lr.Level, CRITICAL)
	}
	if lr.Source != "source" {
		t.Errorf("Incorrect source: %s should be %s", lr.Source, "source")
	}
	if lr.Message != "message" {
		t.Errorf("Incorrect message: %s should be %s", lr.Source, "message")
	}
}
func TestLogger(t *testing.T) {
	sl := NewDefaultLogger(WARNING)
	defer sl.Close(-1)
	if sl == nil {
		t.Fatalf("NewDefaultLogger should never return nil")
	}
	list := sl.getFilterCopy()
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("NewDefaultLogger produced invalid logger (incorrect map count) %v", len(list))
	}
	if list[0].name != "stdout" {
		t.Fatalf("NewDefaultLogger produced invalid default name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("NewDefaultLogger produced nil logger ")
	}
	if list[0].lvl != WARNING {
		t.Fatalf("NewDefaultLogger produced invalid logger (incorrect level) %d", list[0].lvl)
	}

	l := NewLogger()
	l.AddFilter("stdout", DEBUG, nil, nil, NewConsoleLogWriter())

	list = l.getFilterCopy()
	if list == nil {
		t.Fatalf("Nil filter list")
	}
	if len(list) != 1 {
		t.Fatalf("Add produced invalid logger (incorrect map count)")
	}
	if list[0].name != "stdout" {
		t.Fatalf("Add produced invalid name %v", list[0].name)
	}
	if list[0].filter == nil {
		t.Fatalf("Add produced nil logger ")
	}
	if list[0].lvl != DEBUG {
		t.Fatalf("Add produced invalid logger (incorrect level) %d", list[0].lvl)
	}
}

func TestModify(t *testing.T) {
	sl := NewDefaultLogger(WARNING)
	defer sl.Close(-1)
	origList := sl.getFilterRef(false)

	sl.AddFilter("debug", DEBUG, nil, nil, NewConsoleLogWriter())
	addList := sl.getFilterRef(false)

	err := sl.ModifyFilterLvl("stdout", DEBUG, nil, nil)
	if err != nil {
		t.Fatalf("Modify of entry returned failure")
	}
	modList := sl.getFilterRef(false)
	err = sl.ModifyFilterLvl("stdout", DEBUG, nil, nil)
	if err != nil {
		t.Fatalf("Modify of entry returned failure")
	}

	err = sl.ModifyFilterLvl("not", DEBUG, nil, nil)
	if err == nil {
		t.Fatalf("Modify of non existent entry returned success")
	}

	if len(origList) != 1 {
		t.Fatalf("Wrong length of origList: %v", len(origList))
	}
	if len(addList) != 2 {
		t.Fatalf("Wrong length of addList: %v", len(addList))
	}
	if len(modList) != 2 {
		t.Fatalf("Wrong length of modList: %v", len(modList))
	}
	if origList[0].name != "stdout" {
		t.Fatalf("Unexpected array layout %v", origList[0].name)
	}
	if addList[0].name != "stdout" {
		t.Fatalf("Unexpected array layout %v", addList[0].name)
	}
	if modList[0].name != "stdout" {
		t.Fatalf("Unexpected array layout %v", modList[0].name)
	}
	if origList[0].lvl != WARNING {
		t.Fatalf("Unexpected lvl %v", origList[0].lvl)
	}
	if addList[0].lvl != WARNING {
		t.Fatalf("Unexpected lvl %v", addList[0].lvl)
	}
	if modList[0].lvl != DEBUG {
		t.Fatalf("Unexpected lvl %v", modList[0].lvl)
	}
}

func TestNone(t *testing.T) {
	var buf bytes.Buffer

	sl := NewDefaultLoggerWithTarget(NONE, &buf)
	defer sl.Close(-1)

	sl.Debugf("This is a test")
	sl.Flush(10*time.Second)
	if buf.Len() != 0 {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}

	sl.Logf(NONE, "This is a test")
	sl.Flush(10*time.Second)
	if buf.Len() != 0 {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}


}


func TestRegex(t *testing.T) {
	var buf bytes.Buffer

	sl := NewDefaultLoggerWithTarget(INFO, &buf)
	defer sl.Close(-1)

	sl.Debugf("This is a test")
	sl.Flush(10*time.Second)
	if buf.Len() != 0 {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}

	// Still doesn't match with an empty expensive list
	exp := NewExpensive()
	sl.ModifyFilterLvl("stdout", INFO, exp, nil)
	sl.Debugf("This is a test")
	sl.Flush(10*time.Second)
	if buf.Len() != 0 {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}

	// Doesn't match with a non-matching expensive list
	exp = NewExpensive()
	re, err := regexp.Compile(".*wontmatch.*")
	if err != nil {
		t.Errorf("Bad regexp compile %v", err)
	}
	exp.AddSrcRegex(re)
	err = sl.ModifyFilterLvl("stdout", INFO, exp, nil)
	if err != nil {
		t.Errorf("Error modifying level")
	}
	sl.Debugf("This is a test")
	sl.Flush(10*time.Second)
	if buf.Len() != 0 {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}

	//Does match with a matching expensive list
	exp = NewExpensive()
	re, err = regexp.Compile(".*logit_test.*")
	if err != nil {
		t.Errorf("Bad regexp compile %v", err)
	}
	exp.AddSrcRegex(re)
	err = sl.ModifyFilterLvl("stdout", INFO, exp, nil)
	if err != nil {
		t.Errorf("Error modifying level")
	}
	sl.Debugf("This is a test")
	sl.Flush(10*time.Second)
	if !strings.Contains(string(buf.Bytes()), "This is a test") {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}

	//Doesn't match with an expensive list that would match the actual src, but not the
	//hand provided src.
	buf.Reset()
	exp = NewExpensive()
	re, err = regexp.Compile(".*logit_test.*")
	if err != nil {
		t.Errorf("Bad regexp compile %v", err)
	}
	exp.AddSrcRegex(re)
	err = sl.ModifyFilterLvl("stdout", INFO, exp, nil)
	if err != nil {
		t.Errorf("Error modifying level")
	}
	sl.Log(DEBUG, "fakesrc", "This is a test")
	sl.Flush(10*time.Second)
	if buf.Len() != 0 {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}

	//Doest match with an expensive list that wouldn't match the actual src, but does the
	//hand provided src.
	exp = NewExpensive()
	re, err = regexp.Compile("fakesr.*")
	if err != nil {
		t.Errorf("Bad regexp compile %v", err)
	}
	exp.AddSrcRegex(re)
	err = sl.ModifyFilterLvl("stdout", INFO, exp, nil)
	if err != nil {
		t.Errorf("Error modifying level")
	}
	sl.Log(DEBUG, "fakesrc", "This is a test")
	sl.Flush(10*time.Second)

	if !strings.Contains(string(buf.Bytes()), "This is a test") {
		t.Errorf("Unexpected data in buff %v", string(buf.Bytes()))
	}
}

func BenchmarkConsoleLog(b *testing.B) {
	sl := NewDefaultLoggerWithTarget(INFO, ioutil.Discard)
	defer sl.Close(-1)
	for i := 0; i < b.N; i++ {
		sl.Log(WARNING, "here", "This is a log message")
	}
}

func BenchmarkConsoleNotLogged(b *testing.B) {
	sl := NewDefaultLogger(INFO)
	defer sl.Close(-1)
	for i := 0; i < b.N; i++ {
		sl.Log(DEBUG, "here", "This is a log message")
	}
}

func BenchmarkConsoleUtilLog(b *testing.B) {
	sl := NewDefaultLogger(INFO)
	defer sl.Close(-1)
	for i := 0; i < b.N; i++ {
		sl.Infof("%s is a log message", "This")
	}
}

func BenchmarkConsoleUtilNotLog(b *testing.B) {
	sl := NewDefaultLogger(INFO)
	defer sl.Close(-1)
	for i := 0; i < b.N; i++ {
		sl.Debugf("%s is a log message", "This")
	}
}
