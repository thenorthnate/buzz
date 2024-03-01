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

func TestBee(t *testing.T) {
	waiter := make(chan bool, 1)
	bee := NewWorker(&mockTask{dofunc: func(ctx context.Context) error {
		select {
		case <-waiter:
			waiter <- true
		case <-ctx.Done():
			waiter <- true
			return ctx.Err()
		}
		return nil
	}})
	ctx, cancel := context.WithCancel(context.Background())
	go bee.run(ctx)
	waiter <- true
	<-waiter
	cancel()
}

func TestBeeAssembleCallChain(t *testing.T) {
	bee := NewWorker(&mockTask{})
	chain := bee.assembleCallChain()
	if chain.exec == nil {
		t.Fatal("exec was supposed to be defined but was nil instead")
	}
	if chain.next != nil {
		t.Fatal("chain.next was supposed to be nil")
	}
}

func TestBeeWorkTillError(t *testing.T) {
	bee := NewWorker(&mockTask{
		dofunc: func(ctx context.Context) error {
			return errors.New("darn")
		},
	})
	chain := bee.assembleCallChain()
	if chain.exec == nil {
		t.Fatal("exec was supposed to be defined but was nil instead")
	}
	if chain.next != nil {
		t.Fatal("chain.next was supposed to be nil")
	}
	if err := bee.runChainOnce(context.Background(), chain); err != nil {
		t.Fatalf("runChainOnce returned an error: %v", err)
	}
}
