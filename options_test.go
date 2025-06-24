package pgrx

import (
	"context"
	"testing"
	"time"

	"github.com/godepo/groat"
	"github.com/jaswdr/faker/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type optDeps struct {
	Faker faker.Faker
}

type OptState struct {
	RandomString   string
	Faker          faker.Faker
	RandomDeadline time.Duration
	RandomInt32    int32
	migratorConfig MigratorConfig
	Migrator       func(ctx context.Context, migratorConfig MigratorConfig) error
	ctx            context.Context
	FS             afero.Fs
}

func newOptsCase(t *testing.T) *groat.Case[optDeps, OptState, *config] {
	t.Helper()
	tcs := groat.New[optDeps, OptState, *config](t, func(t *testing.T, deps optDeps) *config {
		return &config{}
	})
	tcs.Before(func(t *testing.T, deps optDeps) optDeps {
		deps.Faker = faker.New()
		return deps
	})
	tcs.Go()
	tcs.State.Faker = tcs.Deps.Faker
	return tcs
}

func ArrangeRandomString(t *testing.T, state OptState) OptState {
	t.Helper()
	state.RandomString = state.Faker.RandomStringWithLength(10)
	return state
}

func ArrangeRandomDeadline(t *testing.T, state OptState) OptState {
	t.Helper()
	now := time.Now()
	state.RandomDeadline = now.Sub(state.Faker.Time().TimeBetween(now.Truncate(time.Hour), now))
	return state
}

func ArrangeRandomInt32(t *testing.T, state OptState) OptState {
	t.Helper()
	state.RandomInt32 = state.Faker.Int32()
	return state
}

func ArrangeMigrator(t *testing.T, state OptState) OptState {
	t.Helper()
	state.Migrator = func(ctx context.Context, migratorConfig MigratorConfig) error {
		return nil
	}
	return state
}

func ArrangeContext(t *testing.T, state OptState) OptState {
	t.Helper()
	state.ctx = context.Background()
	return state
}

func TestWithContainerImage(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithContainerImage(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.containerImage)
}

func TestWithContainerImageEnv(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithContainerImageEnv(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.imageEnvValue)
}

func TestWithDBName(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithDBName(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.dbName)
}

func TestWithMigrationsPath(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithMigrationsPath(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.migrationsPath)
}

func TestWithPassword(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithPassword(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.password)
}

func TestWithUserName(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithUserName(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.userName)
}

func TestWithDeadline(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomDeadline)
	WithDeadline(tcs.State.RandomDeadline)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomDeadline, tcs.SUT.deadline)
}

func TestWithPoolMaxConnections(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomInt32)
	WithPoolMaxConnections(tcs.State.RandomInt32)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomInt32, tcs.SUT.poolMaxConns)
}

func TestWithPoolMinConnections(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomInt32)
	WithPoolMinConnections(tcs.State.RandomInt32)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomInt32, tcs.SUT.poolMinConns)
}

func TestWithPoolMaxIdleTime(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomDeadline)
	WithPoolMaxIdleTime(tcs.State.RandomDeadline)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomDeadline, tcs.SUT.poolMaxConnIdleTime)
}

func TestWithMigrator(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeMigrator, ArrangeMigrator, ArrangeContext)
	WithMigrator(tcs.State.Migrator)(tcs.SUT)
	assert.NoError(t, tcs.SUT.migrator(tcs.State.ctx, tcs.State.migratorConfig))
	assert.True(t, tcs.SUT.hasSetMigrator)
}

func TestWithFileSystem(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeMigrator, ArrangeFileSystem)
	WithFileSystem(tcs.State.FS)(tcs.SUT)
	assert.Equal(t, tcs.State.FS, tcs.SUT.fs)

}

func ArrangeFileSystem(t *testing.T, state OptState) OptState {
	state.FS = afero.NewMemMapFs()
	return state
}

func TestWithPoolConfigInjectLabel(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithPoolConfigInjectLabel(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.injectConfigLabel)
}

func TestWithPoolInjectLabel(t *testing.T) {
	tcs := newOptsCase(t).Given(ArrangeRandomString)
	WithPoolInjectLabel(tcs.State.RandomString)(tcs.SUT)
	assert.Equal(t, tcs.State.RandomString, tcs.SUT.injectPoolLabel)
}
