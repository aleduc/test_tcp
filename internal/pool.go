package internal

import (
	"context"
	"net"
)

// Pool struct for control our worker pool.
type Pool struct {
	jobs chan net.Conn
}

type Listener interface {
	Serve(ctx context.Context)
	Close()
}

// Worker process jobs from channel.
type Worker interface {
	Process(jobs <-chan net.Conn)
}

// NewPool init new Pool. maxJobs решил передать в init, чтобы время ее жизни было покороче.
func NewPool(maxJobs int) (*Pool, chan net.Conn) {
	jobs := make(chan net.Conn, maxJobs)
	return &Pool{
		jobs: jobs,
	}, jobs
}

// Start run all workers/consumers
func (p *Pool) Start(maxWorkers int, w Worker) {
	for i := 1; i <= maxWorkers; i++ {
		go w.Process(p.jobs)
	}
}
