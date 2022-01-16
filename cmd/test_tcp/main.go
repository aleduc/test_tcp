package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test_tcp/internal"
	"test_tcp/internal/validate"
	"time"
)

const (
	maxJobs    = 10
	maxWorkers = 10
	lifeTime   = 3 * time.Second
	zeroCnt    = 20
)

func main() {
	logger := internal.NewWithStackTrace(
		log.New(os.Stdout, "INFO: ", 0),
		log.New(os.Stdout, "WARNING: ", 0),
		log.New(os.Stderr, "ERROR: ", 0),
		log.New(os.Stderr, "FATAL: ", 0))

	// Запускаем пул воркеров.
	pool, jobChan := internal.NewPool(maxJobs)
	parser := &validate.Parse{}
	worker := internal.NewTCPWork(
		internal.NewHandle(
			logger,
			validate.NewConnection([]validate.Validator{
				validate.NewExpire(parser, lifeTime, lifeTime),
				validate.NewHashCash(zeroCnt, parser),
			}),
			internal.NewWisdomHT(internal.NewRandomizerTime()),
		))

	pool.Start(maxWorkers, worker)
	// Стартуем листенер, в случае ошибки выходим и пытаемся закрыть(паники из-за nil интерфейса не будет).
	server := internal.NewServer(":2701", logger, jobChan)
	err := server.Listen()
	defer server.Close()
	if err != nil {
		logger.Fatal(err)
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	go server.Serve(ctx)

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	cancelFunc()

}
