package internal

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Addr        string
	logger      Logger
	listener    net.Listener
	handlerJobs chan net.Conn
}

func NewServer(addr string, logger Logger, handlerJobs chan net.Conn) *Server {
	return &Server{Addr: addr, logger: logger, handlerJobs: handlerJobs, listener: &net.TCPListener{}}
}

func (srv *Server) Listen() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":8080"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		srv.logger.Fatal(err)
		return err
	}
	srv.listener = listener
	return nil
}

func (srv *Server) Serve(ctx context.Context) {
	for {
		conn, err := srv.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				srv.logger.Error(err)
			}
		}

		srv.handlerJobs <- conn
	}
}

func (srv *Server) Close() {
	err := srv.listener.Close()
	if err != nil {
		srv.logger.Error(err)
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	defer close(srv.handlerJobs)
	for {
		if len(srv.handlerJobs) == 0 {
			return
		}
		select {
		case <-ticker.C:
			srv.logger.Info(fmt.Sprintf("waiting on %v connections", len(srv.handlerJobs)))
		}
	}

}
