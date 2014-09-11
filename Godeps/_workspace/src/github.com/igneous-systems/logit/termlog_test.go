package logit

import (
	"io"
	"testing"
	"time"
)

var logRecordWriteTests = []struct {
	Test    string
	Record  *LogRecord
	Console string
}{
	{
		Test: "Normal message",
		Record: &LogRecord{
			Level:   CRITICAL,
			Source:  "source",
			Message: "message",
			Created: time.Unix(0, 1234567890123456789).In(time.UTC),
		},
		Console: "[02/13/09 23:31:30] [CRIT] message\n",
	},
}

func TestConsoleLogWriter(t *testing.T) {
	r, w := io.Pipe()
	console := NewConsoleLogWriterWithTarget(w)
	c := make(chan *LogRecord, 32)
	console.StartLogger(c)
	defer func() {
		helper_closeit(c)
	}()

	buf := make([]byte, 1024)
	for _, test := range logRecordWriteTests {
		name := test.Test

		c <- test.Record
		n, _ := r.Read(buf)

		if got, want := string(buf[:n]), test.Console; got != want {
			t.Errorf("%s:  got %q", name, got)
			t.Errorf("%s: want %q", name, want)
		}
	}
}
