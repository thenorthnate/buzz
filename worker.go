package buzz

import "context"

type Worker struct {
	task       Task
	n          int
	middleware []MiddleFunc
}

func NewWorker(task Task) *Worker {
	return &Worker{
		task:       task,
		middleware: make([]MiddleFunc, 0),
	}
}

// SetN provides a mechanism for you to set a certain count of the number of iterations that the worker
// will run. Once it has run the given number of times (successfully), it will exit and clean up.
func (w *Worker) SetN(n int) *Worker {
	w.n = n
	return w
}

// Use adds the given middleware functions to the Bee.
func (w *Worker) Use(middleware ...MiddleFunc) *Worker {
	w.middleware = append(w.middleware, middleware...)
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

func (w *Worker) run(ctx context.Context) {
	callChain := w.assembleCallChain()
	for {
		// execute chain of middleware funcs where each func is passed the next func
		select {
		case <-ctx.Done():
			return
		default:
			if err := w.runChainOnce(ctx, callChain); err != nil {
				return
			}
		}
	}
}

func (w *Worker) runChainOnce(ctx context.Context, callChain *CallChain) error {
	// TODO: depending on settings... maybe do different things depending on the outcome here
	return callChain.Next(ctx)
}

// workTillError should be the final "middleware" called in the call chain. The next call chain
// link will be nil and should not be used hence the underscore.
func (w *Worker) workTillError(ctx context.Context, _ *CallChain) error {
	if w.n > 0 {
		for count := 0; count < w.n; count++ {
			if err := w.task.Do(ctx); err != nil {
				return err
			}
		}
		return nil
	}
	for {
		if err := w.task.Do(ctx); err != nil {
			return err
		}
	}
}
