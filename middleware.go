package buzz

import (
	"context"
	"fmt"
)

func RecoveryMiddleware(ctx context.Context, chain *CallChain) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// TODO : add in the stack trace to this error
			err = fmt.Errorf("panic'd: %v", r)
		}
	}()
	return chain.Next(ctx)
}
