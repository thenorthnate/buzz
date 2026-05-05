package buzz

import (
	"context"
	"sync"
	"time"
)

// Worker wraps your task with additional context to provide a robust operational environment.
type Worker struct {
	task        Task
	middleware  []MiddleFunc
	cadence     time.Duration
	cadenceChan <-chan time.Time
	lifetime    time.Duration
}

// New wraps the task and returns a worker that can be started.
func New(task Task) *Worker {
	cadenceChan := make(chan time.Time)
	close(cadenceChan)
	return &Worker{
		task:        task,
		cadenceChan: cadenceChan,
	}
}

// Use adds the given middleware functions to the Worker.
func (w *Worker) Use(middleware ...MiddleFunc) *Worker {
	w.middleware = append(w.middleware, middleware...)
	return w
}

// WithCadence provides a mechanism through which you can schedule your task to get run on a
// regular interval. By default the tick time is zero meaning that the task is called
// repeatedly as fast as the computer executes it.
func (w *Worker) WithCadence(cadence time.Duration) *Worker {
	w.cadence = cadence
	return w
}

// WithLifetime sets a
func (w *Worker) WithLifetime(lifetime time.Duration) *Worker {
	w.lifetime = lifetime
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
	// set "workTillError" as the final middleware in the callchain
	node.exec = w.workTillError
	return root
}

// Run begins the worker running.
func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	cancel := func() {}
	if w.lifetime > 0 {
		ctx, cancel = context.WithTimeout(ctx, w.lifetime)
	}
	if w.cadence > 0 {
		ticker := time.NewTicker(w.cadence)
		// after Go 1.24, no need to call Stop on tickers!
		w.cadenceChan = ticker.C
	}
	callChain := w.assembleCallChain()
	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()
		for {
			// execute chain of middleware funcs where each func is passed the next func
			select {
			case <-ctx.Done():
				return
			default:
				_ = callChain.Next(ctx)
			}
		}
	}()
}

// workTillError should be the final "middleware" called in the call chain. The next call chain
// link will be nil and should not be used hence the underscore.
func (w *Worker) workTillError(ctx context.Context, _ *CallChain) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.cadenceChan:
			if err := w.task.Do(ctx); err != nil {
				return err
			}
		}
	}
}
