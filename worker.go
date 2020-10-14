package engine

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
)

type work struct {
	Executable
	success chan bool
}

// worker - a unit task executor
type worker struct {
	id    string
	ctx   context.Context
	stop  chan struct{}
	input chan work
	pool  chan chan work
}

// NewWorker - initializes a new worker
func NewWorker(ctx context.Context, pool chan chan work) *worker {
	return &worker{
		id:    fmt.Sprintf("%s-%s", "worker", uuid.NewV4().String()),
		ctx:   ctx,
		pool:  pool,
		input: make(chan work),
		stop:  make(chan struct{}),
	}
}

// Start - readies worker for execution
func (w *worker) Start() {
	log(w.id).Debugf("starting...")
	go w.work()
}

// Stop - stops the worker routine
func (w *worker) Stop() {
	close(w.stop)
}

func (w *worker) execute(work work) {
	if !work.IsCompleted() {
		if err := work.Execute(); err != nil {
			log(w.id).Errorf("error while executing: %+v", work)
			work.OnFailure(err)

			go func() {
				work.success <- false
				close(work.success)
			}()

			return
		}

		log(w.id).Infof("completed executing: %+v", work)
		work.OnSuccess()
	}

	go func() {
		work.success <- true
		close(work.success)
	}()
}

func (w *worker) work() {
	for {
		select {
		case w.pool <- w.input:
			log(w.id).Debugf("back In queue")
		case execute := <-w.input:
			log(w.id).Debugf("executing: %+v", execute)
			w.execute(execute)
		case <-w.stop:
			log(w.id).Debugf("stopping...")
			return
		}
	}
}
