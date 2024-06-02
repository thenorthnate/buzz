package buzz_test

import (
	"context"
	"log"
	"testing"

	"github.com/thenorthnate/buzz"
)

type logTask struct{}

func (t *logTask) Do(ctx context.Context) error {
	log.Println("message here")
	return nil
}

func Example() {
	// This defines some middleware that logs before and after the task runs
	logger := func(ctx context.Context, chain *buzz.CallChain) error {
		// This happens before the task runs
		log.Println("Starting!")
		// This call runs the rest of the middleware and the task
		err := chain.Next(ctx)
		// This runs after the task has completed
		log.Printf("Finished with err=[%v]\n", err)
		return err
	}
	hive := buzz.New()
	worker := buzz.NewWorker(&logTask{}).Use(logger)
	hive.Submit(worker)
	// Some time later... during shutdown
	hive.StopAll()
}

func TestNewTestCallChain(t *testing.T) {
	middleware := func(ctx context.Context, chain *buzz.CallChain) error {
		return nil
	}
	chain := buzz.NewTestCallChain(middleware)
	if err := chain.Next(context.Background()); err != nil {
		t.Fatal("got unexpected error ", err)
	}
}
