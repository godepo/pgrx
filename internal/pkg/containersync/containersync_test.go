package containersync

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/godepo/groat"
	"github.com/godepo/groat/pkg/ctxgroup"
	"github.com/testcontainers/testcontainers-go"
)

type Deps struct {
	WG     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

type State struct {
	cancel  context.CancelFunc
	handler func(ctx context.Context) error
}

func newTerminatorCase(t *testing.T) *groat.Case[Deps, State, func()] {
	var handler func() error
	tcs := groat.New[Deps, State, func()](t, func(t *testing.T, deps Deps) func() {
		return func() {
			_ = handler()
		}
	})
	tcs.Before(func(t *testing.T, deps Deps) Deps {
		deps.WG = &sync.WaitGroup{}
		deps.WG.Add(2)
		deps.ctx = context.Background()
		deps.ctx, deps.cancel = context.WithCancel(deps.ctx)
		deps.ctx = ctxgroup.WithWaitGroup(deps.ctx, deps.WG)
		return deps
	})
	tcs.Given(func(t *testing.T, state State) State {
		state.cancel = tcs.Deps.cancel
		return state
	})
	tcs.After(func(t *testing.T, deps Deps) {
		deps.cancel()
		deps.WG.Wait()
	})
	tcs.Go()
	tcs.SUT = Terminator(tcs.Deps.ctx, func(ctx context.Context, _ ...testcontainers.TerminateOption) error {
		err := tcs.State.handler(ctx)
		tcs.Deps.WG.Done()
		return err
	})
	return tcs
}

func TestTerminator(t *testing.T) {
	t.Run("should be able terminate without error", func(t *testing.T) {
		tcs := newTerminatorCase(t)
		tcs.Given(TriggeredHandler(nil))
		go tcs.SUT()

	})
	t.Run("should be able terminate with error", func(t *testing.T) {
		tcs := newTerminatorCase(t)
		tcs.Given(TriggeredHandler(errors.New("unexpected error")))
		go tcs.SUT()
	})
}

func TriggeredHandler(err error) groat.Given[State] {
	return func(t *testing.T, state State) State {
		state.handler = func(ctx context.Context) error {
			return err
		}
		return state
	}
}
