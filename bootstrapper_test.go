package pgrx

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestBootstrapper(t *testing.T) {
	tcs := newMigratorTestCase(t).Given(ArrangeExpectError)
	tcs.When(ActDirOpen(tcs.State.ExpectError)).Then(ExpectError)
	_, tcs.State.ResultError = bootstrapper[any](config{
		fs:             tcs.Deps.FS,
		migrationsPath: "./sql",
	})(context.Background())

	t.Run("should be able to be able faile", func(t *testing.T) {
		t.Run("when can't run postgres container", func(t *testing.T) {
			exp := errors.New(uuid.NewString())
			_, err := bootstrapper[Deps](config{
				hasSetMigrator: true,
				runner: func(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (PostgresContainer, error) {
					return nil, exp
				},
			})(context.Background())
			require.ErrorIs(t, err, exp)
		})

		t.Run("when can't get connection string from postgres container", func(t *testing.T) {
			pc := NewMockPostgresContainer(t)
			exp := errors.New(uuid.NewString())

			pc.EXPECT().ConnectionString(mock.Anything).Return("", exp)

			_, err := bootstrapper[Deps](config{
				hasSetMigrator: true,
				runner: func(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (PostgresContainer, error) {
					return pc, nil
				},
			})(context.Background())
			require.ErrorIs(t, err, exp)
		})

		t.Run("when can't parse connection string", func(t *testing.T) {
			pc := NewMockPostgresContainer(t)

			pc.EXPECT().ConnectionString(mock.Anything).Return(uuid.NewString(), nil)

			_, err := bootstrapper[Deps](config{
				hasSetMigrator: true,
				runner: func(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (PostgresContainer, error) {
					return pc, nil
				},
			})(context.Background())
			require.Error(t, err)
			assert.ErrorContains(t, err, "cannot parse")
		})

		t.Run("when can't create connection poo;", func(t *testing.T) {
			pc := NewMockPostgresContainer(t)

			exp := errors.New(uuid.NewString())

			pc.EXPECT().ConnectionString(mock.Anything).Return("", nil)

			_, err := bootstrapper[Deps](config{
				hasSetMigrator: true,
				runner: func(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (PostgresContainer, error) {
					return pc, nil
				},
				poolConstructor: func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
					return nil, exp
				},
			})(context.Background())
			assert.ErrorIs(t, err, exp)
		})
	})

}
