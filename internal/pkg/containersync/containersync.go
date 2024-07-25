package containersync

import (
	"context"
	"fmt"

	"github.com/godepo/groat/pkg/ctxgroup"
)

func Terminator(ctx context.Context, terminate func(ctx context.Context) error) func() {
	return func() {
		<-ctx.Done()
		defer func() {
			ctxgroup.DoneFrom(ctx)
		}()
		err := terminate(context.Background())
		if err != nil {
			fmt.Printf("---[GOAT]: error terminating postgres container: %v\n", err)
		}
	}

}
