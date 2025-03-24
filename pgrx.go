//go:generate mockery
package pgrx

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/godepo/groat/integration"
	"github.com/godepo/groat/pkg/ctxgroup"
	"github.com/godepo/groat/pkg/generics"
	"github.com/godepo/pgrx/internal/pkg/containersync"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultTestsDeadline = time.Second * 5
	defaultPoolMaxCons   = 8
	defaultPoolMinCons   = 2
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type config struct {
	userName            string
	dbName              string
	password            string
	containerImage      string
	imageEnvValue       string
	deadline            time.Duration
	poolMaxConns        int32
	poolMinConns        int32
	poolMaxConnIdleTime time.Duration
	migrator            Migrator
	hasSetMigrator      bool
	migrationsPath      string
	fs                  afero.Fs
	runner              containerRunner
	poolConstructor     func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error)
}

type containerRunner func(
	ctx context.Context,
	img string,
	opts ...testcontainers.ContainerCustomizer,
) (PostgresContainer, error)

type PostgresContainer interface {
	ConnectionString(ctx context.Context, args ...string) (string, error)
	Terminate(ctx context.Context, opts ...testcontainers.TerminateOption) error
}

type Container[T any] struct {
	forks          *atomic.Int32
	pc             PostgresContainer
	connString     string
	pgCfg          *pgxpool.Config
	ctx            context.Context
	migrator       Migrator
	migrationsPath string
	root           *pgxpool.Pool
}

type MigratorConfig struct {
	DBName   string
	Pool     DB
	Path     string
	UserName string
}

type Migrator func(ctx context.Context, migratorConfig MigratorConfig) error

func (c *Container[T]) Injector(t *testing.T, to T) T {
	t.Helper()
	cfg, err := pgxpool.ParseConfig(c.connString)
	require.NoError(t, err)

	cfg.MinConns = c.pgCfg.MinConns
	cfg.MaxConns = c.pgCfg.MaxConns
	cfg.MaxConnIdleTime = c.pgCfg.MaxConnIdleTime
	cfg.ConnConfig.Database = fmt.Sprintf("%s_%d", cfg.ConnConfig.Database, c.forks.Add(1))

	_, err = c.root.Exec(c.ctx, fmt.Sprintf("CREATE DATABASE %s WITH OWNER = '%s'",
		cfg.ConnConfig.Database, cfg.ConnConfig.User))
	require.NoError(t, err,
		"can't created database=%s for user %s", cfg.ConnConfig.Database, cfg.ConnConfig.User,
	)

	p, err := pgxpool.NewWithConfig(c.ctx, cfg)
	require.NoError(t, err)

	require.NoError(t, p.Ping(c.ctx))
	err = c.migrator(c.ctx, MigratorConfig{
		DBName:   cfg.ConnConfig.Database,
		Pool:     p,
		Path:     c.migrationsPath,
		UserName: cfg.ConnConfig.User,
	})
	require.NoError(t, err)
	generics.Injector[*pgxpool.Config, T](t, cfg, to, "pgxconfig")
	return generics.Injector[*pgxpool.Pool, T](t, p, to, "pgxpool")
}

func New[T any](all ...Option) integration.Bootstrap[T] {
	cfg := config{
		userName:            "test",
		dbName:              "test",
		password:            "test",
		containerImage:      "postgres:16",
		imageEnvValue:       "GROAT_I9N_PG_IMAGE",
		deadline:            defaultTestsDeadline,
		poolMaxConns:        defaultPoolMaxCons,
		poolMinConns:        defaultPoolMinCons,
		poolMaxConnIdleTime: time.Minute,
		migrationsPath:      "../sql/migrations",
		fs:                  afero.NewOsFs(),
		runner: func(
			ctx context.Context,
			img string,
			opts ...testcontainers.ContainerCustomizer,
		) (PostgresContainer, error) {
			return postgres.Run(ctx, img, opts...)
		},
		poolConstructor: pgxpool.NewWithConfig,
	}
	for _, opt := range all {
		opt(&cfg)
	}

	if env := os.Getenv(cfg.imageEnvValue); env != "" {
		cfg.containerImage = env
	}

	return bootstrapper[T](cfg)
}

func bootstrapper[T any](cfg config) integration.Bootstrap[T] {
	return func(ctx context.Context) (integration.Injector[T], error) {
		if !cfg.hasSetMigrator {
			mig, err := PlainMigrator(cfg.fs, cfg.migrationsPath)
			if err != nil {
				return nil, err
			}
			cfg.migrator = mig
		}

		pc, err := cfg.runner(
			ctx,
			cfg.containerImage,
			postgres.WithPassword(cfg.password),
			postgres.WithDatabase(cfg.dbName),
			postgres.WithUsername(cfg.userName),
			testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
				ContainerRequest: testcontainers.ContainerRequest{
					Tmpfs: map[string]string{"/tmpfs": "rw"},
					Env:   map[string]string{"PGDATA": "/tmpfs"},
				},
			}),
			testcontainers.WithWaitStrategyAndDeadline(
				cfg.deadline,
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).WithStartupTimeout(4*time.Second), //nolint:mnd
			),
		)
		if err != nil {
			return nil, fmt.Errorf("postgres container failed to run: %w", err)
		}
		ctxgroup.IncAt(ctx)

		go containersync.Terminator(ctx, pc.Terminate)()
		container, err := newContainer[T](ctx, pc, cfg)
		if err != nil {
			return nil, err
		}

		return container.Injector, nil
	}
}

func newContainer[T any](ctx context.Context,
	pc PostgresContainer,
	cfg config,
) (*Container[T], error) {
	container := &Container[T]{
		forks:          &atomic.Int32{},
		pc:             pc,
		ctx:            ctx,
		migrator:       cfg.migrator,
		migrationsPath: cfg.migrationsPath,
	}

	connString, err := pc.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get connection string for postgres container: %w", err)
	}
	container.connString = connString
	container.pgCfg, err = pgxpool.ParseConfig(container.connString)
	if err != nil {
		return nil, fmt.Errorf("can't parse connection string: %w", err)
	}

	container.pgCfg.MinConns = cfg.poolMinConns
	container.pgCfg.MaxConnIdleTime = cfg.poolMaxConnIdleTime
	container.pgCfg.MaxConns = cfg.poolMaxConns

	root, err := cfg.poolConstructor(ctx, container.pgCfg)
	if err != nil {
		return nil, fmt.Errorf("can't connect to root db: %w", err)
	}
	container.root = root
	return container, nil
}
