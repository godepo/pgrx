package pgrx

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/godepo/groat"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MigratorDeps struct {
	FS       *MockFs
	Dir      *MockFile
	File     *MockFile
	FileInfo *MockFileInfo
	DB       *MockDB
}

type MigratorState struct {
	ExpectError    error
	OpenDirError   error
	ResultError    error
	Migration      string
	MigratorConfig MigratorConfig
}

type MigratorSUT func(fs afero.Fs, path string) (Migrator, error)

func newMigratorTestCase(t *testing.T) *groat.Case[MigratorDeps, MigratorState, MigratorSUT] {
	t.Helper()
	t.Parallel()
	tcs := groat.New[MigratorDeps, MigratorState, MigratorSUT](t, func(t *testing.T, deps MigratorDeps) MigratorSUT {
		return PlainMigrator
	})
	tcs.Before(func(t *testing.T, deps MigratorDeps) MigratorDeps {
		deps.FS = NewMockFs(t)
		deps.FileInfo = NewMockFileInfo(t)
		deps.Dir = NewMockFile(t)
		deps.File = NewMockFile(t)
		deps.DB = NewMockDB(t)
		return deps
	})
	tcs.Go()
	return tcs
}

func TestPlainMigrator(t *testing.T) {
	t.Run("should be able to be able", func(t *testing.T) {
		tcs := newMigratorTestCase(t)
		tcs.Given(ArrangeExpectError, ArrangePlainMigratorConfig(tcs.Deps.DB)).
			When(
				ActDirOpen(nil),
				ActDirRead(nil),
				ActFileInfoName,
				ActFileOpen(nil),
				ActFileRead(nil),
				ActFileClose(nil),
				ActDirClose(nil),
				ActMigrationExecute(nil),
			).
			Then(ExpectNoError)
		var migrations Migrator
		migrations, err := tcs.SUT(tcs.Deps.FS, "./sql")
		require.NoError(t, err)
		tcs.State.ResultError = migrations(context.Background(), tcs.State.MigratorConfig)
	})
	t.Run("should be able return error", func(t *testing.T) {
		t.Run("when given not exists path", func(t *testing.T) {
			tcs := newMigratorTestCase(t).
				Given(ArrangeExpectError)
			tcs.
				When(ActDirOpen(tcs.State.ExpectError)).
				Then(ExpectError)
			_, tcs.State.ResultError = tcs.SUT(tcs.Deps.FS, "./sql")
		})
		t.Run("when can't read exist dir", func(t *testing.T) {
			tcs := newMigratorTestCase(t)
			tcs.Given(ArrangeExpectError).
				When(
					ActDirOpen(nil),
					ActDirRead(tcs.State.ExpectError),
					ActDirClose(nil),
				).
				Then(ExpectError)
			_, tcs.State.ResultError = tcs.SUT(tcs.Deps.FS, "./sql")
		})

		t.Run("when can't open exist file", func(t *testing.T) {
			tcs := newMigratorTestCase(t)
			tcs.Given(ArrangeExpectError).
				When(
					ActDirOpen(nil),
					ActDirRead(nil),
					ActFileInfoName,
					ActFileOpen(tcs.State.ExpectError),
					ActDirClose(nil),
				).
				Then(ExpectError)
			_, tcs.State.ResultError = tcs.SUT(tcs.Deps.FS, "./sql")
		})

		t.Run("when can't read exist file", func(t *testing.T) {
			tcs := newMigratorTestCase(t)
			tcs.Given(ArrangeExpectError).
				When(
					ActDirOpen(nil),
					ActDirRead(nil),
					ActFileInfoName,
					ActFileOpen(nil),
					ActFileRead(tcs.State.ExpectError),
					ActFileClose(nil),
					ActDirClose(nil),
				).
				Then(ExpectError)
			_, tcs.State.ResultError = tcs.SUT(tcs.Deps.FS, "./sql")
		})

		t.Run("when failed sql migration request", func(t *testing.T) {
			tcs := newMigratorTestCase(t)
			tcs.Given(ArrangeExpectError, ArrangePlainMigratorConfig(tcs.Deps.DB)).
				When(
					ActDirOpen(nil),
					ActDirRead(nil),
					ActFileInfoName,
					ActFileOpen(nil),
					ActFileRead(nil),
					ActFileClose(nil),
					ActDirClose(nil),
					ActMigrationExecute(tcs.State.ExpectError),
				).
				Then(ExpectError)
			var migrations Migrator
			migrations, err := tcs.SUT(tcs.Deps.FS, "./sql")
			require.NoError(t, err)
			tcs.State.ResultError = migrations(context.Background(), tcs.State.MigratorConfig)
		})
	})
}

func ActMigrationExecute(err error) groat.When[MigratorDeps, MigratorState] {
	return func(t *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
		deps.DB.EXPECT().Exec(mock.Anything, state.Migration).Return(pgconn.CommandTag{}, err)
		return state
	}
}

func ArrangePlainMigratorConfig(db *MockDB) groat.Given[MigratorState] {
	return func(t *testing.T, state MigratorState) MigratorState {
		state.MigratorConfig = MigratorConfig{
			DBName:   "",
			Pool:     db,
			Path:     "",
			UserName: "",
		}
		return state
	}

}

func ExpectNoError(t *testing.T, state MigratorState) {
	require.NoError(t, state.ResultError)
}

func ActFileInfoName(_ *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
	deps.FileInfo.EXPECT().Name().Return("test.sql")
	return state
}

func ActFileClose(err error) groat.When[MigratorDeps, MigratorState] {
	return func(t *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
		deps.File.EXPECT().Close().Return(err)
		return state
	}
}

func ActFileOpen(err error) groat.When[MigratorDeps, MigratorState] {
	return func(t *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
		var f afero.File
		if err == nil {
			f = deps.File
		}
		deps.FS.EXPECT().Open(filepath.Join("sql", "test.sql")).Return(f, err)
		return state
	}
}

func ActFileRead(err error) groat.When[MigratorDeps, MigratorState] {
	return func(t *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
		if err != nil {
			deps.File.EXPECT().Read(mock.Anything).Return(0, err)
			return state
		}
		state.Migration = uuid.NewString()
		attempts := atomic.Int32{}
		deps.File.EXPECT().Read(mock.Anything).RunAndReturn(func(bytes []byte) (int, error) {
			cnt := attempts.Add(1)
			if cnt != 1 {
				return 0, io.EOF
			}
			return copy(bytes, state.Migration), nil
		})

		return state
	}
}

func ActDirClose(err error) groat.When[MigratorDeps, MigratorState] {
	return func(t *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
		deps.Dir.EXPECT().Close().Return(err)
		return state
	}
}

func ActDirRead(err error) groat.When[MigratorDeps, MigratorState] {
	return func(t *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
		var infos []os.FileInfo
		if err == nil {
			infos = []os.FileInfo{deps.FileInfo}
			deps.FileInfo.EXPECT().IsDir().Return(false)
			deps.FileInfo.EXPECT().Name().Return("test.sql")
		}
		deps.Dir.EXPECT().Readdir(0).Return(infos, err)
		return state
	}
}

func ActDirOpen(err error) groat.When[MigratorDeps, MigratorState] {
	return func(t *testing.T, deps MigratorDeps, state MigratorState) MigratorState {
		deps.FS.EXPECT().Open("./sql").Return(deps.Dir, err)
		return state
	}
}

func ExpectError(t *testing.T, state MigratorState) {
	assert.ErrorIs(t, state.ResultError, state.ExpectError)
}

func ArrangeExpectError(_ *testing.T, state MigratorState) MigratorState {
	state.ExpectError = errors.New(uuid.NewString())
	return state
}
