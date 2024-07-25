package pgrx

import (
	"context"
	"testing"
)

func TestBootstrapper(t *testing.T) {
	tcs := newMigratorTestCase(t).Given(ArrangeExpectError)
	tcs.When(ActDirOpen(tcs.State.ExpectError)).Then(ExpectError)
	_, tcs.State.ResultError = bootstrapper[any](config{
		fs:             tcs.Deps.FS,
		migrationsPath: "./sql",
	})(context.Background())

}
