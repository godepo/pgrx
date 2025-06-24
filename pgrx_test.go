package pgrx

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/godepo/groat"
	"github.com/godepo/groat/integration"
	"github.com/godepo/groat/pkg/ctxgroup"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

type SystemUnderTest struct {
}

type State struct {
}

type Deps struct {
	DB *pgxpool.Pool `groat:"pgxpool"`
}

var suite *integration.Container[Deps, State, *SystemUnderTest]

func TestMain(m *testing.M) {
	_ = os.Setenv("GROAT_I9N_PG_IMAGE", "postgresql:17")
	suite = integration.New[Deps, State, *SystemUnderTest](m,
		func(t *testing.T) *groat.Case[Deps, State, *SystemUnderTest] {
			tcs := groat.New[Deps, State, *SystemUnderTest](t, func(t *testing.T, deps Deps) *SystemUnderTest {
				return &SystemUnderTest{}
			})
			return tcs
		},
		New[Deps](
			WithMigrationsPath("./sql"),
			WithDBName("test"),
			WithUserName("test"),
			WithPassword("test"),
			WithContainerImageEnv("postgresql:16"),
			WithDeadline(time.Minute),
		),
	)
	os.Exit(suite.Go())
}

func TestContainer_Injector(t *testing.T) {
	t.Run("should be able to be able", func(t *testing.T) {
		suite.Case(t)
	})
}

func TestNew(t *testing.T) {
	bt := New[Deps](WithContainerImage("postgresnosql:16"), WithMigrationsPath("./sql"))
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	ctx = ctxgroup.WithWaitGroup(ctx, wg)
	_, err := bt(ctx)
	require.Error(t, err)
	cancel()
	fmt.Println("cancelled")
	wg.Wait()
}
