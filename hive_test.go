package buzz

import (
	"context"
	"log"
	"testing"
	"time"
)

func logMiddleware(ctx context.Context, chain *CallChain) error {
	log.Println("task starting")
	err := chain.Next(ctx)
	log.Println("task complete")
	return err
}

func TestHive_RemoveDoneWorkers(t *testing.T) {
	hive := New()
	worker := NewWorker(&mockTask{})
	worker.done.Store(true)
	hive.colony = append(hive.colony, worker)
	hive.removeDoneWorkers()
	if len(hive.colony) != 0 {
		t.Fatalf("hive is supposed to be empty but has %v workers in it still", len(hive.colony))
	}
}

func TestHive_Middleware(t *testing.T) {
	waiter := make(chan struct{}, 1)
	hive := New()
	worker := NewWorker(&mockTask{dofunc: func(ctx context.Context) error {
		select {
		case waiter <- struct{}{}:
		default:
		}
		<-ctx.Done()
		return nil
	}}).Use(logMiddleware)
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

func TestHive_MultipleWorkers(t *testing.T) {
	waiter1 := make(chan struct{}, 1)
	waiter2 := make(chan struct{}, 1)
	hive := New()
	hive.Use(logMiddleware)
	worker1 := NewWorker(&mockTask{dofunc: func(ctx context.Context) error {
		select {
		case waiter1 <- struct{}{}:
		default:
		}
		<-ctx.Done()
		return nil
	}})
	hive.Submit(worker1)
	worker2 := NewWorker(&mockTask{dofunc: func(ctx context.Context) error {
		select {
		case waiter2 <- struct{}{}:
		default:
		}
		<-ctx.Done()
		return nil
	}}).Tick(time.Microsecond)
	hive.Submit(worker2)

	<-waiter1
	if len(worker1.middleware) != 1 {
		t.Fatalf("worker1 was supposed to have 1 middlefunc but had %v", len(worker1.middleware))
	}
	<-waiter2
	if len(worker2.middleware) != 1 {
		t.Fatalf("worker2 was supposed to have 1 middlefunc but had %v", len(worker2.middleware))
	}
	hive.StopAll()
	if len(hive.colony) > 0 {
		t.Fatalf("the hive is supposed to be empty but it still has %v workers in it", len(hive.colony))
	}
}
