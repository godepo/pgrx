package example

import (
	"context"
	"os"
	"testing"

	"github.com/godepo/groat"
	"github.com/godepo/groat/integration"
	"github.com/godepo/pgrx"
)

func mainProvider(t *testing.T) *groat.Case[RepoDeps, State, *Repository] {
	tcs := groat.New[RepoDeps, State, *Repository](t, func(t *testing.T, deps RepoDeps) *Repository {
		return New(deps.FirstDB, deps.SecondDB)
	})
	tcs.Given(func(t *testing.T, state State) State {
		state.Faker = tcs.Deps.Faker
		state.ctx = context.Background()
		return state
	})
	return tcs
}

func TestMain(m *testing.M) {
	suite = integration.New[RepoDeps, State, *Repository](m, mainProvider,
		pgrx.New[RepoDeps](
			pgrx.WithContainerImage("docker.io/postgres:16"),
			pgrx.WithMigrationsPath("./first"),
			pgrx.WithPoolInjectLabel("first_pool"),
			pgrx.WithPoolConfigInjectLabel("first_pool_config"),
		),
		pgrx.New[RepoDeps](
			pgrx.WithContainerImage("docker.io/postgres:16"),
			pgrx.WithMigrationsPath("./second"),
			pgrx.WithPoolInjectLabel("second_pool"),
			pgrx.WithPoolConfigInjectLabel("second_pool_config"),
		),
	)
	os.Exit(suite.Go())
}
