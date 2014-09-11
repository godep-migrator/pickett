package filelog

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	. "logit"
)

var now time.Time = time.Unix(0, 1234567890123456789).In(time.UTC)

func closeit(c chan *LogRecord) error {
	var rec LogRecord
	rec.Level = CLOSE
	rec.SetReplyChan(make(chan error))
	c <- &rec
	err := <-rec.GetReplyChan()
	return err
}

func newLogRecord(lvl Level, src string, msg string) *LogRecord {
	return &LogRecord{
		Level:   lvl,
		Source:  src,
		Created: now,
		Message: msg,
	}
}

const testLogFile = "_logtest.log"

func TestFileLogWriter(t *testing.T) {
	defer os.Remove(testLogFile)

	w := NewFileLogWriter(testLogFile, false)
	if w == nil {
		t.Fatalf("Invalid return: w should not be nil")
	}
	c := make(chan *LogRecord, 32)
	w.StartLogger(c)
	c <- newLogRecord(CRITICAL, "source", "message")
	closeit(c)

	if contents, err := ioutil.ReadFile(testLogFile); err != nil {
		t.Errorf("read(%q): %s", testLogFile, err)
	} else if len(contents) != 50 {
		t.Errorf("malformed filelog: %q (%d bytes)", string(contents), len(contents))
	}
}

func TestXMLLogWriter(t *testing.T) {
	defer os.Remove(testLogFile)

	w := NewXMLLogWriter(testLogFile, false)
	if w == nil {
		t.Fatalf("Invalid return: w should not be nil")
	}
	c := make(chan *LogRecord, 32)
	w.StartLogger(c)
	c <- newLogRecord(CRITICAL, "source", "message")
	err := closeit(c)
	if err != nil {
		t.Errorf("Error flushing %v", err)
	}
	if contents, err := ioutil.ReadFile(testLogFile); err != nil {
		t.Errorf("read(%q): %s", testLogFile, err)
	} else if len(contents) != 185 {
		t.Errorf("malformed xmllog: %q (%d bytes)", string(contents), len(contents))
	}
}

func TestLogOutput(t *testing.T) {
	exp := []string{"[CRIT] This message is level 7",
		"[EROR] This message is level EROR",
		"[WARN] This message is level WARN",
		"[INFO] This message is level INFO",
		"[TRAC] This message is level 3",
		"[DEBG] This message is level DEBG",
		"[FINE] This message is level FINE",
		"[FNST] This message is level FNST",
	}

	// Unbuffered output

	l := NewLogger()
	defer os.Remove(testLogFile)

	// Delete and open the output log without a timestamp (for a constant md5sum)
	l.AddFilter("file", FINEST, nil, nil, NewFileLogWriter(testLogFile, false).SetFormat("[%L] %M"))

	// Send some log messages
	l.Log(CRITICAL, "testsrc1", fmt.Sprintf("This message is level %d", int(CRITICAL)))
	l.Logf(ERROR, "This message is level %v", ERROR)
	l.Logf(WARNING, "This message is level %s", WARNING)
	l.Logc(INFO, func() string { return "This message is level INFO" })
	l.Tracef("This message is level %d", int(TRACE))
	l.Debugf("This message is level %s", DEBUG)
	l.Finec(func() string { return fmt.Sprintf("This message is level %v", FINE) })
	l.Finestf("This message is level %v", FINEST)

	err := l.Close(-1)
	if err != nil {
		t.Fatalf("Error closing log: %s", err)
	}

	contents, err := ioutil.ReadFile(testLogFile)
	if err != nil {
		t.Fatalf("Could not read output log: %s", err)
	}

	arr := strings.Split(string(contents), "\n")
	if arr[len(arr)-1] == "" {
		arr = arr[0 : len(arr)-1]
	}
	if len(arr) != len(exp) {
		t.Errorf("Expected:\n%v---\nSeen:\n%v\n---", exp, arr)
		t.Fatalf("Mismatched log length")
	}
	for i := range arr {
		if arr[i] != exp[i] {
			t.Errorf("Expected:\n%v---\nSeen:\n%v\n---", exp[i], arr[i])
			t.Fatalf("Mismatched entry line %d", i)
		}
	}
}

func BenchmarkFileNotLogged(b *testing.B) {
	sl := NewLogger()
	b.StopTimer()
	sl.AddFilter("file", INFO, nil, nil, NewFileLogWriter("benchlog.log", false))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.Log(DEBUG, "here", "This is a log message")
	}
	b.StopTimer()
	os.Remove("benchlog.log")
}

func BenchmarkFileUtilLog(b *testing.B) {
	sl := NewLogger()
	b.StopTimer()
	sl.AddFilter("file", INFO, nil, nil, NewFileLogWriter("benchlog.log", false))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.Infof("%s is a log message", "This")
	}
	b.StopTimer()
	os.Remove("benchlog.log")
}

func BenchmarkFileUtilNotLog(b *testing.B) {
	sl := NewLogger()
	b.StopTimer()
	sl.AddFilter("file", INFO, nil, nil, NewFileLogWriter("benchlog.log", false))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.Debugf("%s is a log message", "This")
	}
	b.StopTimer()
	os.Remove("benchlog.log")
}

func BenchmarkFileLog(b *testing.B) {
	sl := NewLogger()
	b.StopTimer()
	sl.AddFilter("file", INFO, nil, nil, NewFileLogWriter("benchlog.log", false))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.Log(WARNING, "here", "This is a log message")
	}
	b.StopTimer()
	os.Remove("benchlog.log")
}
