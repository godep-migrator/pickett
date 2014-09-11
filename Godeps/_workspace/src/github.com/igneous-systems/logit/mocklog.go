package logit

import (
	"sync"
)

type MockSkellLogEntry struct {
	Call int
	Rec  *LogRecord
	Wake interface{}
}

const (
	INVALID = iota
	MOCK_FLUSH
	MOCK_CLOSE
	MOCK_CLEANUP
	MOCK_WAKE
	MOCK_LOG
)

// This is the standard writer that prints to standard output.
type MockSkelLog struct {
	lock       sync.Mutex
	CallRecord chan *MockSkellLogEntry
	BlockWake  chan error

	shouldBlock                                 bool
	flushError, closeError, logError, wakeError error

	SkelLog *SkelLog
}

func NewMockSkelLog(wakeChan chan interface{}) *MockSkelLog {
	var ms MockSkelLog
	ms.CallRecord = make(chan *MockSkellLogEntry, 200)
	ms.BlockWake = make(chan error)
	ms.SkelLog = NewSkelLog(&ms, wakeChan)
	return &ms
}

func (ms *MockSkelLog) Flush() error {
	ms.CallRecord <- &MockSkellLogEntry{MOCK_FLUSH, nil, nil}
	if ms.shouldBlock {
		err, ok := <-ms.BlockWake
		if !ok {
			panic("Unexpected closed channel")
		}
		return err
	}
	ms.lock.Lock()
	defer ms.lock.Unlock()
	return ms.flushError
}

func (ms *MockSkelLog) Close() error {
	ms.CallRecord <- &MockSkellLogEntry{MOCK_CLOSE, nil, nil}
	if ms.shouldBlock {
		err, ok := <-ms.BlockWake
		if !ok {
			panic("Unexpected closed channel")
		}
		return err
	}

	ms.lock.Lock()
	defer ms.lock.Unlock()
	return ms.closeError
}

func (ms *MockSkelLog) Cleanup() {
	ms.CallRecord <- &MockSkellLogEntry{MOCK_CLEANUP, nil, nil}
	if ms.shouldBlock {
		_, ok := <-ms.BlockWake
		if !ok {
			panic("Unexpected closed channel")
		}
		return
	}

	ms.lock.Lock()
	defer ms.lock.Unlock()
}

func (ms *MockSkelLog) Wake(obj interface{}) error {
	ms.CallRecord <- &MockSkellLogEntry{MOCK_WAKE, nil, obj}
	if ms.shouldBlock {
		err, ok := <-ms.BlockWake
		if !ok {
			panic("Unexpected closed channel")
		}
		return err
	}

	ms.lock.Lock()
	defer ms.lock.Unlock()
	return ms.wakeError
}

func (ms *MockSkelLog) LogNormal(rec *LogRecord) error {
	ms.CallRecord <- &MockSkellLogEntry{MOCK_LOG, rec, nil}
	if ms.shouldBlock {
		err, ok := <-ms.BlockWake
		if !ok {
			panic("Unexpected closed channel")
		}
		return err
	}

	ms.lock.Lock()
	defer ms.lock.Unlock()
	return ms.logError
}

func (ms *MockSkelLog) SetErrors(flushE error, closeE error, logE error, wakeE error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.flushError = flushE
	ms.closeError = closeE
	ms.logError = logE
	ms.wakeError = wakeE
}

func (ms *MockSkelLog) ShouldBlock(block bool) {
	ms.shouldBlock = block
}
