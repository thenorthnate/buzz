package buzz

import (
	"context"
	"errors"
	"testing"
)

type mockTask struct {
	dofunc func(ctx context.Context) error
}

func (task *mockTask) Do(ctx context.Context) error {
	return task.dofunc(ctx)
}

// func TestWorker(t *testing.T) {
// 	waiter := make(chan bool, 1)
// 	worker := NewWorker(&mockTask{dofunc: func(ctx context.Context) error {
// 		select {
// 		case <-waiter:
// 			waiter <- true
// 		case <-ctx.Done():
// 			waiter <- true
// 			return ctx.Err()
// 		}
// 		return nil
// 	}})
// 	ctx, cancel := context.WithCancel(context.Background())
// 	go worker.run(ctx)
// 	waiter <- true
// 	<-waiter
// 	cancel()
// }

func TestWorkerAssembleCallChain(t *testing.T) {
	worker := NewWorker(&mockTask{})
	chain := worker.assembleCallChain()
	if chain.exec == nil {
		t.Fatal("exec was supposed to be defined but was nil instead")
	}
	if chain.next != nil {
		t.Fatal("chain.next was supposed to be nil")
	}
}

func TestWorkerWorkTillError(t *testing.T) {
	worker := NewWorker(&mockTask{
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
	if err := worker.runChainOnce(context.Background(), chain); err != nil {
		t.Fatalf("runChainOnce returned an error: %v", err)
	}
}
