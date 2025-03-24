package containersync

import (
	"context"
	"log"

	"github.com/godepo/groat/pkg/ctxgroup"
	"github.com/testcontainers/testcontainers-go"
)

func Terminator(ctx context.Context, terminate func(context.Context, ...testcontainers.TerminateOption) error) func() {
	return func() {
		<-ctx.Done()
		defer func() {
			ctxgroup.DoneFrom(ctx)
		}()
		err := terminate(context.Background()) //nolint:contextcheck
		if err != nil {
			log.Printf("---[GOAT]: error terminating postgres container: %v\n", err)
		}
	}
}
