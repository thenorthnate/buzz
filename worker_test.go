package buzz

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type mockTask struct {
	dofunc func(ctx context.Context) error
}

func (task *mockTask) Do(ctx context.Context) error {
	return task.dofunc(ctx)
}

func TestWorker(t *testing.T) {
	waiter := make(chan struct{}, 1)
	task := &mockTask{dofunc: func(ctx context.Context) error {
		select {
		case waiter <- struct{}{}:
		default:
		}
		<-ctx.Done()
		return nil
	}}
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	New(task).Run(ctx, wg)
	<-waiter
	cancel()
	wg.Wait()
}

func TestWorkerAssembleCallChain(t *testing.T) {
	worker := New(&mockTask{})
	chain := worker.assembleCallChain()
	if chain.exec == nil {
		t.Fatal("exec was supposed to be defined but was nil instead")
	}
	if chain.next != nil {
		t.Fatal("chain.next was supposed to be nil")
	}
}

func TestWorkerWorkTillError(t *testing.T) {
	worker := New(&mockTask{
		dofunc: func(ctx context.Context) error {
			return errors.New("darn")
		},
	})
	chain := worker.assembleCallChain()
	if chain.exec == nil {
		t.Fatal("exec was supposed to be defined but was nil instead")
	}
	if chain.next != nil {
		t.Fatal("chain.next was supposed to be nil")
	}
	if err := chain.Next(context.Background()); err == nil {
		t.Fatal("Next was supposed to return an error but did not")
	}
}
