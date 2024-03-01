package buzz

import (
	"context"

	"github.com/thenorthnate/evs"
)

func RecoveryMiddleware(ctx context.Context, chain *CallChain) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = evs.Newf("worker panic'd: %v", r).Err()
		}
	}()
	return chain.Next(ctx)
}
