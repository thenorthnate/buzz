package buzz

import (
	"context"
)

// Task represents the thing that you want accomplished.
type Task interface {
	// Do should perform the desired work of the Worker. If the context is cancelled, it should
	// return an error. If no error is returned, [Do] is called repeatedly in a loop.
	Do(ctx context.Context) error
}

// CallChain represents a linked list that provides the mechanism through which middleware
// can be implemented.
type CallChain struct {
	next *CallChain
	exec MiddleFunc
}

// Next is used to allow the chain to proceed processing. When it returns, you can assume
// that all middleware as well as the task itself executed and returned.
func (chain *CallChain) Next(ctx context.Context) error {
	return chain.exec(ctx, chain.next)
}

// MiddleFunc defines the type of any middleware that can be used in the hive.
type MiddleFunc func(ctx context.Context, chain *CallChain) error

// NewTestCallChain creates a new [CallChain] that simply executes the given [MiddleFunc].
// The provided [MiddleFunc] will recieve a nil [CallChain]. This function is a utility to
// make it easy to test your own middleware.
func NewTestCallChain(exec MiddleFunc) *CallChain {
	return &CallChain{
		exec: exec,
	}
}
