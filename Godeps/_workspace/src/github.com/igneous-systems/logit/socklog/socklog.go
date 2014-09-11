// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package socklog

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	. "logit"
)

// This log writer sends output to a socket
type SocketLogWriter struct {
	skelLog  *SkelLog
	proto    string
	hostport string
	sock     net.Conn
}

func (sw *SocketLogWriter) LogNormal(rec *LogRecord) error {
	// Marshall into JSON
	js, err := json.Marshal(rec)
	if err != nil {
		fmt.Fprint(os.Stderr, "SocketLogWriter(%q): %s", sw.hostport, err)
		return err
	}

	_, err = sw.sock.Write(js)
	if err != nil {
		fmt.Fprint(os.Stderr, "SocketLogWriter(%q): %s", sw.hostport, err)
		return err
	}
	return nil
}

func (sw *SocketLogWriter) Flush() error {
	return nil
}

func (sw *SocketLogWriter) Close() error {
	err := sw.sock.Close()
	sw.sock = nil
	return err
}

func (sw *SocketLogWriter) Cleanup() {
	if sw.sock != nil && sw.proto == "tcp" {
		sw.sock.Close()
		sw.sock = nil
	}
}

func (sw *SocketLogWriter) Wake(obj interface{}) error {
	return nil
}

func NewSocketLogWriter(proto, hostport string) *SocketLogWriter {
	var sw SocketLogWriter
	sw.proto = proto
	sw.hostport = hostport

	sock, err := net.Dial(proto, hostport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewSocketLogWriter(%q): %s\n", hostport, err)
		return nil
	}

	sw.sock = sock
	sw.skelLog = NewSkelLog(&sw, nil)
	return &sw
}

func (sw *SocketLogWriter) StartLogger(c chan *LogRecord) {
	sw.skelLog.StartLogger(c)
}

func (sw *SocketLogWriter) Closed() bool {
	return sw.skelLog.Closed()
}
