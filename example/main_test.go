package example

import (
	"context"
	"os"
	"testing"

	"github.com/godepo/groat"
	"github.com/godepo/groat/integration"
	"github.com/godepo/pgrx"
	"github.com/jaswdr/faker/v2"
)

func mainProvider(t *testing.T) *groat.Case[RepoDeps, State, *Repository] {
	tcs := groat.New[RepoDeps, State, *Repository](t, func(t *testing.T, deps RepoDeps) *Repository {
		return New(deps.DB)
	})
	tcs.Before(func(t *testing.T, deps RepoDeps) RepoDeps {
		deps.Faker = faker.New()
		deps.MockDB = NewMockDB(t)
		deps.MockRows = NewMockRows(t)
		return deps
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
			pgrx.WithMigrationsPath("../sql"),
		),
	)
	os.Exit(suite.Go())
}
