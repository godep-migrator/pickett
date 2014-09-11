package logit

import (
	"errors"
	"testing"
	"time"
)

func helperCleanupClose(ms *MockSkelLog) {
	ms.ShouldBlock(false)
	select {
	case ms.BlockWake <- nil:
	default:
	}

	helper_closeitHideErr(ms.SkelLog.c)
}

func recvNBErr(t *testing.T, ms chan *MockSkellLogEntry) *MockSkellLogEntry {
	select {
	case le := <-ms:
		return le
	default:
		t.Errorf("No message received")
	}
	return nil
}

func recvTOErr(t *testing.T, ms chan *MockSkellLogEntry) *MockSkellLogEntry {
	select {
	case le := <-ms:
		return le
	case <-time.After(time.Second * 30):
		t.Errorf("No message received by timeout")
	}
	return nil
}

func recvExpectTimeout(t *testing.T, ms chan *MockSkellLogEntry) *MockSkellLogEntry {
	select {
	case le := <-ms:
		t.Errorf("Unexpected receive")
		return le
	case <-time.After(time.Second * 1):
	}
	return nil
}

func TestLogs(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)
	results := ms.CallRecord
	rec := newLogRecord(DEBUG, "source", "message")
	rec2 := newLogRecord(DEBUG, "source2", "message2")
	c <- rec
	c <- rec2

	res := recvTOErr(t, results)
	res2 := recvTOErr(t, results)

	if res.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res.Call)
	}
	if res2.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res2.Call)
	}
	if res.Rec != rec {
		t.Errorf("Bad record %v", res.Rec)
	}
	if res2.Rec != rec2 {
		t.Errorf("Bad record %v", res2.Rec)
	}
}

// Basic test that Log, Flush, and Close all pass through
func TestAllNonWakeCalls(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)
	results := ms.CallRecord
	rec := newLogRecord(DEBUG, "source", "message")
	rec2 := newLogRecord(DEBUG, "source2", "message2")

	c <- rec
	err := helper_flush(c)
	if err != nil {
		t.Errorf("Error flushing %v", err)
	}
	c <- rec2
	err = helper_closeit(c)
	if err != nil {
		t.Errorf("Error flushing %v", err)
	}

	res := recvTOErr(t, results)
	res2 := recvTOErr(t, results)
	res3 := recvTOErr(t, results)
	res4 := recvTOErr(t, results)
	if res.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res.Call)
	}
	if res2.Call != MOCK_FLUSH {
		t.Errorf("Wrong call type %v", res2.Call)
	}
	if res3.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res3.Call)
	}
	if res4.Call != MOCK_CLOSE {
		t.Errorf("Wrong call type %v", res4.Call)
	}

}

// Basic test that Wake actually wakes
func TestWake(t *testing.T) {
	c := make(chan *LogRecord, 32)
	wake := make(chan interface{})
	ms := NewMockSkelLog(wake)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)
	results := ms.CallRecord

	wake <- "wake"
	res := recvTOErr(t, results)
	if res.Call != MOCK_WAKE {
		t.Errorf("Wrong call type %v", res.Call)
	}
	swake, ok := res.Wake.(string)
	if !ok || swake != "wake" {
		t.Errorf("Wrong data %v", res.Call)
	}
}

// Test that flush doesn't return until all previous logs finish
func TestFlushWaits(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)
	results := ms.CallRecord
	rec := newLogRecord(DEBUG, "source", "message")
	rec2 := newLogRecord(DEBUG, "source2", "message2")
	rec3 := newLogRecord(DEBUG, "source3", "message3")

	ms.ShouldBlock(true)

	var recf LogRecord
	recf.Level = FLUSH
	recf.SetReplyChan(make(chan error))
	c <- rec
	c <- rec2
	c <- &recf
	c <- rec3
	select {
	case err := <-recf.GetReplyChan():
		t.Errorf("Expected to block %v", err)
	case <-time.After(time.Second * 1):
	}
	ms.BlockWake <- nil
	select {
	case err := <-recf.GetReplyChan():
		t.Errorf("Expected to block %v", err)
	case <-time.After(time.Second * 1):
	}
	ms.BlockWake <- nil
	select {
	case err := <-recf.GetReplyChan():
		t.Errorf("Expected to block %v", err)
	case <-time.After(time.Second * 1):
	}
	ms.BlockWake <- nil
	select {
	case err := <-recf.GetReplyChan():
		if err != nil {
			t.Errorf("Unexpected error")
		}
	case <-time.After(time.Second * 20):
		t.Errorf("Expected not to block")
	}

	ms.BlockWake <- nil

	res := recvTOErr(t, results)
	res2 := recvTOErr(t, results)
	res3 := recvTOErr(t, results)
	res4 := recvTOErr(t, results)
	if res.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res.Call)
	}
	if res2.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res2.Call)
	}
	if res3.Call != MOCK_FLUSH {
		t.Errorf("Wrong call type %v", res3.Call)
	}
	if res4.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res4.Call)
	}
}

// Test that close doesn't return until all previous logs finish and throws
// away later logs.
func TestCloseWaits(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)
	results := ms.CallRecord
	rec := newLogRecord(DEBUG, "source", "message")
	rec2 := newLogRecord(DEBUG, "source2", "message2")
	rec3 := newLogRecord(DEBUG, "source3", "message3")

	ms.ShouldBlock(true)

	var recf LogRecord
	recf.Level = CLOSE
	recf.SetReplyChan(make(chan error))
	c <- rec
	c <- rec2
	c <- &recf
	c <- rec3
	select {
	case err := <-recf.GetReplyChan():
		t.Errorf("Expected to block %v", err)
	case <-time.After(time.Second * 1):
	}
	ms.BlockWake <- nil
	select {
	case err := <-recf.GetReplyChan():
		t.Errorf("Expected to block %v", err)
	case <-time.After(time.Second * 1):
	}
	ms.BlockWake <- nil
	select {
	case err := <-recf.GetReplyChan():
		t.Errorf("Expected to block %v", err)
	case <-time.After(time.Second * 1):
	}
	ms.BlockWake <- nil
	select {
	case err := <-recf.GetReplyChan():
		if err != nil {
			t.Errorf("Unexpected error")
		}
	case <-time.After(time.Second * 20):
		t.Errorf("Expected not to block")
	}

	ms.BlockWake <- nil

	res := recvTOErr(t, results)
	res2 := recvTOErr(t, results)
	res3 := recvTOErr(t, results)
	res4 := recvTOErr(t, results)
	if res.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res.Call)
	}
	if res2.Call != MOCK_LOG {
		t.Errorf("Wrong call type %v", res2.Call)
	}
	if res3.Call != MOCK_CLOSE {
		t.Errorf("Wrong call type %v", res3.Call)
	}
	if res4.Call != MOCK_CLEANUP {
		t.Errorf("Wrong call type %v", res4.Call)
	}
}

// Test that close or flush after close errors.
func TestDoubleClose(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)
	ms.ShouldBlock(true)

	var recf1 LogRecord
	recf1.Level = CLOSE
	recf1.SetReplyChan(make(chan error))
	var recf2 LogRecord
	recf2.Level = CLOSE
	recf2.SetReplyChan(make(chan error))
	var recf3 LogRecord
	recf3.Level = FLUSH
	recf3.SetReplyChan(make(chan error))
	c <- &recf1
	c <- &recf2
	c <- &recf3
	ms.ShouldBlock(false)
	select {
	case ms.BlockWake <- nil:
	default:
	}
	var err1, err2, err3 error

	select {
	case err1 = <-recf1.GetReplyChan():
	case <-time.After(time.Second * 20):
		t.Errorf("Expected not to block")
	}
	select {
	case err2 = <-recf2.GetReplyChan():
	case <-time.After(time.Second * 20):
		t.Errorf("Expected not to block")
	}
	select {
	case err3 = <-recf3.GetReplyChan():
	case <-time.After(time.Second * 20):
		t.Errorf("Expected not to block")
	}

	if err1 != nil {
		t.Errorf("Unexpected error %v", err1)
	}
	if err2 == nil {
		t.Errorf("Unexpected success")
	}
	if err3 == nil {
		t.Errorf("Unexpected success")
	}
}

func helper_ignorepanic(c chan *LogRecord, rec *LogRecord) {
	defer func() {
		_ = recover()
	}()
	c <- rec
}

// Test that log error terminates.
func TestLogError(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	results := ms.CallRecord
	defer helperCleanupClose(ms)
	rec := newLogRecord(DEBUG, "source", "message")
	rec2 := newLogRecord(DEBUG, "source2", "message2")

	ms.SetErrors(nil, nil, errors.New("Log error"), nil)
	c <- rec
	helper_ignorepanic(c, rec2)

	res := recvTOErr(t, results)
	res2 := recvTOErr(t, results)
	ms.SetErrors(nil, nil, nil, nil)
	if res.Call != MOCK_LOG {
		t.Errorf("Unexpected type", res.Call)
	}
	if res2.Call != MOCK_CLEANUP {
		t.Errorf("Unexpected type %v", res2.Call)
	}
	recvExpectTimeout(t, results)
}

// Test that errors out of close return
func TestCloseError(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)

	ms.SetErrors(nil, errors.New("Close error"), nil, nil)
	err := helper_closeit(c)
	if err == nil {
		t.Errorf("Didn't see expected error")
	}
}

// Test that errors out of flush return, but don't terminate
func TestFlushError(t *testing.T) {
	c := make(chan *LogRecord, 32)
	ms := NewMockSkelLog(nil)
	ms.SkelLog.StartLogger(c)
	results := ms.CallRecord
	defer helperCleanupClose(ms)
	rec := newLogRecord(DEBUG, "source", "message")

	ms.SetErrors(errors.New("Flush error"), nil, nil, nil)
	err := helper_flush(c)
	if err == nil {
		t.Errorf("Didn't see expected error")
	}
	ms.SetErrors(nil, nil, nil, nil)
	c <- rec
	_ = recvTOErr(t, results)
	res2 := recvTOErr(t, results)
	if res2.Call != MOCK_LOG {
		t.Errorf("Unexpected type %v", res2.Call)
	}
}

// Basic test errors out of Wake terminate
func TestWakeError(t *testing.T) {
	c := make(chan *LogRecord, 32)
	wake := make(chan interface{})
	ms := NewMockSkelLog(wake)
	ms.SkelLog.StartLogger(c)
	defer helperCleanupClose(ms)
	results := ms.CallRecord
	rec := newLogRecord(DEBUG, "source", "message")

	ms.SetErrors(nil, nil, nil, errors.New("Flush error"))
	wake <- "wake"
	_ = recvTOErr(t, results)
	ms.SetErrors(nil, nil, nil, nil)
	helper_ignorepanic(c, rec)
	res2 := recvTOErr(t, results)
	if res2.Call != MOCK_CLEANUP {
		t.Errorf("Unexpected type %v", res2.Call)
	}

	recvExpectTimeout(t, results)
}
