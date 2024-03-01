//go:build integration

package buzz_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/thenorthnate/buzz"
)

type logTask struct{}

func (t *logTask) Do(ctx context.Context) error {
	log.Println("message here")
	return nil
}

func TestWorker(t *testing.T) {
	hive := buzz.New()
	worker := buzz.NewWorker(&logTask{}).Tick(time.Second).Use(buzz.RecoveryMiddleware)
	hive.Submit(worker)
	time.Sleep(5 * time.Second)
	worker.Stop()
}
