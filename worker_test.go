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
	worker := NewWorker(task)
	hive := New()
	hive.Submit(worker)
	<-waiter
	hive.StopAll()
	if len(hive.colony) > 0 {
		t.Fatalf("the hive is supposed to be empty but it still has %v workers in it", len(hive.colony))
	}
}

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
	if err := worker.runChainOnce(context.Background(), chain); err == nil {
		t.Fatal("runChainOnce was supposed to return an error but did not")
	}
}
