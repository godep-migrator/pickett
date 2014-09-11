// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package logit

import (
	"fmt"
	"sync"
)

// The interface for anything using this skeleton writer.
type SmallWriter interface {
	LogNormal(lr *LogRecord) error
	Flush() error
	Close() error
	Wake(obj interface{}) error
	Cleanup()
}

// This is the skeleton writer that wraps a real writer.
type SkelLog struct {
	closed      bool
	closedMutex sync.Mutex
	c           chan *LogRecord
	sw          SmallWriter
	wake        chan interface{}
}

// This creates a new SkelLog
func NewSkelLog(sw SmallWriter, wake chan interface{}) *SkelLog {
	var cw SkelLog
	cw.sw = sw
	cw.wake = wake
	return &cw
}

func (cw *SkelLog) run() {
	var err error
	c := cw.c
	sw := cw.sw
	wake := cw.wake
Outer:
	for {
		select {
		case obj := <-wake:
			err = sw.Wake(obj)
			if err != nil {
				break Outer
			}
		case rec, ok := <-c:
			if !ok {
				break Outer
			}
			if rec.Level < 0 {
				switch rec.Level {
				case FLUSH:
					err = sw.Flush()
					rec.GetReplyChan() <- err
				case CLOSE:
					err = sw.Close()
					rec.GetReplyChan() <- err
					break Outer
				}
				continue Outer
			}

			err = sw.LogNormal(rec)
			if err != nil {
				break Outer
			}
		}
	}
	cw.closedMutex.Lock()
	cw.closed = true
	cw.closedMutex.Unlock()
	close(c)
	for rec := range c {
		switch rec.Level {
		case FLUSH:
			rec.GetReplyChan() <- fmt.Errorf("Closed channel")
		case CLOSE:
			rec.GetReplyChan() <- fmt.Errorf("Closed channel")
		}
		continue
	}
	sw.Cleanup()

}

func (cw *SkelLog) Closed() bool {
	cw.closedMutex.Lock()
	defer cw.closedMutex.Unlock()
	return cw.closed
}

func (cw *SkelLog) StartLogger(c chan *LogRecord) {
	cw.c = c
	go cw.run()
}
