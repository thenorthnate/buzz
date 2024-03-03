//go:build integration

package buzz_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/thenorthnate/buzz"
)

func logMiddleware(ctx context.Context, chain *buzz.CallChain) error {
	log.Println("task starting")
	err := chain.Next(ctx)
	log.Println("task complete")
	return err
}

func TestIntegrationWorker(t *testing.T) {
	hive := buzz.New()
	worker := buzz.
		NewWorker(&logTask{}).
		Tick(time.Second).
		Use(logMiddleware)
	hive.Submit(worker)
	time.Sleep(5 * time.Second)
	hive.StopAll()
}
