package internal

import (
	"net"
)

type Handler interface {
	Handle(conn net.Conn)
}

type TCPWork struct {
	handler Handler
}

func NewTCPWork(handler Handler) *TCPWork {
	return &TCPWork{handler: handler}
}

func (t *TCPWork) Process(jobs <-chan net.Conn) {
	for k := range jobs {
		t.handler.Handle(k)
	}
}
