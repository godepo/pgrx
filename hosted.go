package pgrx

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/godepo/groat/integration"
	"github.com/godepo/groat/pkg/generics"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var ErrRequireNamespacePrefixForHostedDB = errors.New("hosted db requires namespace prefix")

type hostedPostgres[T any] struct {
	root  DB
	cfg   config
	forks *atomic.Int32
	ctx   context.Context
}

func hostedBootstrapper[T any](cfg config) integration.Bootstrap[T] {
	return func(ctx context.Context) (integration.Injector[T], error) {
		if cfg.hostedDBNamespace == "" {
			return nil, ErrRequireNamespacePrefixForHostedDB
		}

		if cfg.migrator == nil {
			mig, err := PlainMigrator(cfg.fs, cfg.migrationsPath)
			if err != nil {
				return nil, err
			}
			cfg.migrator = mig
		}

		local := &hostedPostgres[T]{
			cfg:   cfg,
			forks: &atomic.Int32{},
			ctx:   ctx,
		}

		opts, err := pgxpool.ParseConfig(cfg.hostedDSN)
		if err != nil {
			return nil, fmt.Errorf("can't parse hosted dsn: %w", err)
		}

		conn, err := cfg.poolConstructor(ctx, opts)

		if err != nil {
			return nil, fmt.Errorf("can't create connection to hosted db: %w", err)
		}

		local.root = conn

		return local.Injector, nil
	}
}

func (c *hostedPostgres[T]) Injector(t *testing.T, to T) T {
	t.Helper()

	cfg, err := pgxpool.ParseConfig(c.cfg.hostedDSN)
	require.NoError(t, err)

	dbName := fmt.Sprintf("%s_%d", c.cfg.hostedDBNamespace, c.forks.Add(1))
	cfg.ConnConfig.Database = dbName

	_, err = c.root.Exec(
		c.ctx,
		"CREATE DATABASE "+dbName,
	)
	require.NoError(t, err,
		"can't created database=%s for user %s",
		dbName, cfg.ConnConfig.User,
	)

	var con *pgxpool.Pool
	t.Cleanup(func() {
		if con != nil {
			con.Close()
		}
		_, err := c.root.Exec(c.ctx, "DROP DATABASE "+dbName)
		if err != nil {
			t.Logf("can't cleanup database %s: %v", dbName, err)
		}
	})

	con, err = c.cfg.poolConstructor(c.ctx, cfg)
	require.NoError(t, err)

	require.NoError(t, con.Ping(c.ctx))

	err = c.cfg.migrator(c.ctx, MigratorConfig{
		Pool:     con,
		DBName:   dbName,
		Path:     c.cfg.migrationsPath,
		UserName: cfg.ConnConfig.User,
	})
	require.NoError(t, err)
	res := generics.Injector(t, con, to, c.cfg.injectPoolLabel)
	res = generics.Injector(t, cfg, res, c.cfg.injectConfigLabel)
	return res
}
