package buzz

import (
	"context"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	worker := NewWorker(&mockTask{
		dofunc: func(ctx context.Context) error {
			panic("darn")
		},
	}).Use(RecoveryMiddleware)
	chain := worker.assembleCallChain()
	if chain.exec == nil {
		t.Fatal("exec was supposed to be defined but was nil instead")
	}
	if chain.next == nil {
		t.Fatal("chain.next was not supposed to be nil")
	}
	if err := worker.runChainOnce(context.Background(), chain); err == nil {
		t.Fatal("runChainOnce was supposed to return an error but did not")
	}
}
