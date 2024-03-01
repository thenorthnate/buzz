package buzz

import (
	"context"
	"errors"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	bee := NewWorker(&mockTask{
		dofunc: func(ctx context.Context) error {
			return errors.New("darn")
		},
	}).Use(RecoveryMiddleware)
	chain := bee.assembleCallChain()
	if chain.exec == nil {
		t.Fatal("exec was supposed to be defined but was nil instead")
	}
	if chain.next == nil {
		t.Fatal("chain.next was not supposed to be nil")
	}
	if err := bee.runChainOnce(context.Background(), chain); err == nil {
		t.Fatal("runChainOnce was supposed to return an error but did not")
	}
}
