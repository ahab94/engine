package engine

import (
	"context"
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// Engine - for creating workers and distributing jobs
type Engine struct {
	id      string
	ctx     context.Context
	stop    chan struct{}
	pool    chan chan work
	input   chan work
	workers []*worker
	start   *sync.Once
}

// NewEngine - initializing a new engine
func NewEngine(ctx context.Context) *Engine {
	return &Engine{
		id:    fmt.Sprintf("%s-%s", "dispatcher", uuid.NewV4().String()),
		ctx:   ctx,
		start: new(sync.Once),
	}
}

// Start - starting workers and setting up dispatcher for use
func (e *Engine) Start(workerCount uint) {
	e.start.Do(func() {
		e.stop = make(chan struct{})
		e.pool = make(chan chan work)
		e.input = make(chan work)
		e.workers = make([]*worker, 0)

		for i := 0; i <= int(workerCount); i++ {
			worker := NewWorker(e.ctx, e.pool)
			e.workers = append(e.workers, worker)
			worker.Start()
		}

		go e.dispatch()
	})
}

// Stop - closes channels/goroutines
func (e *Engine) Stop() {
	defer func() { e.start = new(sync.Once) }()
	for _, worker := range e.workers {
		worker.Stop()
	}
	close(e.stop)
}

// Do - executes work
func (e *Engine) Do(executable Executable) <-chan bool {
	done := make(chan bool, 0)

	e.input <- work{
		Executable: executable,
		success:    done,
	}

	return done
}

func (e *Engine) dispatch() {
	for {
		select {
		case work := <-e.input:
			log(e.id).Debugf("dispatching: %v", work)
			worker := <-e.pool
			worker <- work

		case <-e.stop:
			log(e.id).Debugf("stopping...")
			return
		}
	}
}
