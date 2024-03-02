package buzz

import (
	"context"
	"testing"
)

func TestHive_RemoveDoneWorkers(t *testing.T) {
	hive := New()
	worker := NewWorker(&mockTask{})
	worker.done = true
	hive.colony = append(hive.colony, worker)
	hive.removeDoneWorkers()
	if len(hive.colony) != 0 {
		t.Fatalf("hive is supposed to be empty but has %v workers in it still", len(hive.colony))
	}
}

func TestHive_Middleware(t *testing.T) {
	waiter := make(chan struct{}, 1)
	hive := New()
	hive.Use(RecoveryMiddleware)
	worker := NewWorker(&mockTask{dofunc: func(ctx context.Context) error {
		select {
		case waiter <- struct{}{}:
		default:
		}
		<-ctx.Done()
		return nil
	}})
	hive.Submit(worker)
	<-waiter
	if len(worker.middleware) != 1 {
		t.Fatalf("worker was supposed to have 1 middlefunc but had %v", len(worker.middleware))
	}
	hive.StopAll()
	if len(hive.colony) > 0 {
		t.Fatalf("the hive is supposed to be empty but it still has %v workers in it", len(hive.colony))
	}
}
