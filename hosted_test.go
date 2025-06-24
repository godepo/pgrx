package pgrx

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHostedPostgres(t *testing.T) {
	envHost := "http://locaghost:8123/?dial_timeout=200ms&max_execution_time=60"
	t.Run("should be able to be able", func(t *testing.T) {
		tc := suite.Case(t)

		cfg := config{
			hostedDSN:         envHost,
			fs:                afero.OsFs{},
			migrationsPath:    "./sql",
			hostedDBNamespace: uuid.NewString(),
			poolConstructor: func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
				return tc.Deps.DB, nil
			},
		}

		res, err := hostedBootstrapper[Deps](cfg)(t.Context())
		require.NoError(t, err)
		require.NotNil(t, res)
	})
	t.Run("should be able failed", func(t *testing.T) {
		t.Run("when given not exists namespace", func(t *testing.T) {
			t.Setenv("GROAT_I9N_PG_DSN",
				uuid.NewString())

			res, err := New[Deps](
				WithMigrator(
					func(ctx context.Context, migratorConfig MigratorConfig) error {
						return nil
					},
				),
			)(t.Context())

			require.ErrorIs(t, ErrRequireNamespacePrefixForHostedDB, err)
			require.Nil(t, res)
		})

		t.Run("when dsn can't be parsed", func(t *testing.T) {
			t.Setenv("GROAT_I9N_PG_DSN", uuid.NewString())

			res, err := New[Deps](
				WithHostedNamespace(uuid.NewString()),
				WithMigrationsPath("./sql"),
			)(t.Context())
			require.Error(t, err)
			assert.ErrorContains(t, err, "can't parse hosted dsn", err.Error())
			require.Nil(t, res)
		})

		t.Run("when can't create connection to hosted db", func(t *testing.T) {
			tc := suite.Case(t)
			exp := errors.New(uuid.NewString())
			cfg := config{
				hostedDSN:         envHost,
				fs:                afero.OsFs{},
				hostedDBNamespace: uuid.NewString(),
				migrationsPath:    "./sql",
				poolConstructor: func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
					return tc.Deps.DB, exp
				},
			}

			res, err := hostedBootstrapper[Deps](cfg)(t.Context())
			require.Error(t, err)
			assert.ErrorIs(t, err, exp, err.Error())
			require.Nil(t, res)
		})

		t.Run("when can't open migrations dir", func(t *testing.T) {
			cfg := config{
				hostedDSN:         envHost,
				fs:                afero.OsFs{},
				hostedDBNamespace: uuid.NewString(),
				migrationsPath:    uuid.NewString(),
			}

			res, err := hostedBootstrapper[Deps](cfg)(t.Context())
			require.Error(t, err)
			assert.ErrorContains(t, err, "no such file or directory")
			require.Nil(t, res)
		})
	})
}

func TestHostedPostgresql_Injector(t *testing.T) {
	t.Run("should be able to be able", func(t *testing.T) {

		tc := suite.Case(t)
		envHost := "host=localhost"

		hostedClick := hostedPostgres[Deps]{
			root: tc.Deps.DB,
			cfg: config{
				hostedDSN:         envHost,
				hostedDBNamespace: "pgrx_" + strings.Replace(uuid.NewString(), "-", "_", -1),
				poolConstructor: func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
					return tc.Deps.DB, nil
				},
				migrator: func(ctx context.Context, migratorConfig MigratorConfig) error {
					return nil
				},
			},
			forks: &atomic.Int32{},
			ctx:   t.Context(),
		}
		var deps Deps

		deps = hostedClick.Injector(t, deps)
	})

	t.Run("should be able failed", func(t *testing.T) {
		t.Run("when can't drop database", func(t *testing.T) {
			tc := suite.Case(t)
			envHost := "host=localhost"

			hostedClick := hostedPostgres[Deps]{
				root: tc.Deps.DB,
				cfg: config{
					hostedDSN:         envHost,
					hostedDBNamespace: "pgrx_" + strings.Replace(uuid.NewString(), "-", "_", -1),
					poolConstructor: func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
						return tc.Deps.DB, nil
					},
					migrator: func(ctx context.Context, migratorConfig MigratorConfig) error {
						return nil
					},
				},
				forks: &atomic.Int32{},
				ctx:   t.Context(),
			}
			var deps Deps

			deps = hostedClick.Injector(t, deps)
		})
	})
}
