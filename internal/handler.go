package internal

import (
	"github.com/dropbox/godropbox/errors"
	"io"
	"net"
	"strings"
	"time"
)

const (
	bufferSize = 64 // 7 байт на дату и время + 32 байта на хэш
	timeout    = 1 * time.Second
)

var (
	errNotCorrectDataLen = errors.New("not correct data length")
)

type Validator interface {
	Validate(ipAddress string, data []byte) error
}

type WisdomGetter interface {
	Get() (string, error)
}

type Handle struct {
	logger   Logger
	verifier Validator
	wisdom   WisdomGetter
}

func NewHandle(logger Logger, verifier Validator, wisdom WisdomGetter) *Handle {
	return &Handle{logger: logger, verifier: verifier, wisdom: wisdom}
}

func (h *Handle) Handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	err := conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		h.logger.Error(err)
		return
	}

	buf := make([]byte, bufferSize+1)
	n, err := io.LimitReader(conn, bufferSize+1).Read(buf)
	if err != nil {
		h.logger.Error(err)
		return
	}
	if n != bufferSize {
		h.logger.Error(errNotCorrectDataLen)
		return
	}

	err = h.verifier.Validate(strings.Split(conn.RemoteAddr().String(), ":")[0], buf[:bufferSize])
	if err != nil {
		h.logger.Error(err)
		return
	}

	val, err := h.wisdom.Get()
	if err != nil {
		h.logger.Error(err)
		return
	}

	_, err = conn.Write([]byte(val))
	if err != nil {
		h.logger.Error(err)

	}

	return
}
