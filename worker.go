package buzz

import (
	"context"
	"sync"
	"time"
)

// Worker wraps your task with additional context to provide a robust operational environment.
type Worker struct {
	task           Task
	middleware     []MiddleFunc
	cancel         context.CancelFunc
	tick           time.Duration
	tickChan       <-chan time.Time
	notifyComplete chan struct{}
	done           bool
}

// NewWorker wraps the task and returns a worker that can be submitted to the hive.
func NewWorker(task Task) *Worker {
	tickChan := make(chan time.Time)
	close(tickChan)
	return &Worker{
		task:       task,
		middleware: make([]MiddleFunc, 0),
		tickChan:   tickChan,
	}
}

// Use adds the given middleware functions to the Worker.
func (w *Worker) Use(middleware ...MiddleFunc) *Worker {
	w.middleware = append(w.middleware, middleware...)
	return w
}

// Tick provides a mechanism through which you can schedule your task to get run on a
// regular interval. By default the tick time is zero meaning that the task is called
// repeatedly as fast as the computer executes it.
func (w *Worker) Tick(tick time.Duration) *Worker {
	w.tick = tick
	return w
}

func (w *Worker) assembleCallChain() *CallChain {
	root := &CallChain{}
	node := root
	for _, mfunc := range w.middleware {
		node.exec = mfunc
		node.next = &CallChain{}
		node = node.next
	}
	node.exec = w.workTillError
	return root
}

func (w *Worker) run(block *sync.WaitGroup) {
	defer block.Done()
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	if w.tick > 0 {
		ticker := time.NewTicker(w.tick)
		defer ticker.Stop()
		w.tickChan = ticker.C
	}
	callChain := w.assembleCallChain()
	for {
		// execute chain of middleware funcs where each func is passed the next func
		select {
		case <-ctx.Done():
			return
		default:
			_ = w.runChainOnce(ctx, callChain)
		}
	}
}

func (w *Worker) runChainOnce(ctx context.Context, callChain *CallChain) error {
	return callChain.Next(ctx)
}

// workTillError should be the final "middleware" called in the call chain. The next call chain
// link will be nil and should not be used hence the underscore.
func (w *Worker) workTillError(ctx context.Context, _ *CallChain) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.tickChan:
			if err := w.task.Do(ctx); err != nil {
				return err
			}
		}
	}
}

// Stop issues a command to the hive to stop this worker from running and remove it.
func (w *Worker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.done = true
	w.notifyComplete <- struct{}{}
}
