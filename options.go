package pgrx

import (
	"time"

	"github.com/spf13/afero"
)

type Option func(*config)

// WithContainerImageEnv override ENV contained container image name, by default using GROAT_I9N_PG_IMAGE.
func WithContainerImageEnv(env string) Option {
	return func(c *config) {
		c.imageEnvValue = env
	}
}

// WithContainerImage override initial container image for run postgresql instance,
// by default using value 'postgres:16'.
func WithContainerImage(image string) Option {
	return func(c *config) {
		c.containerImage = image
	}
}

// WithUserName override initial username of the user to be created when the container starts. By default,
// using value 'test'.
func WithUserName(userName string) Option {
	return func(c *config) {
		c.userName = userName
	}
}

// WithPassword override the initial password of the user to be created when the container starts as superuser by
// default using value 'test'.
func WithPassword(password string) Option {
	return func(c *config) {
		c.password = password
	}
}

// WithDBName override the initial database for check correct migrationsPath process by default using value 'test'.
func WithDBName(dbName string) Option {
	return func(c *config) {
		c.dbName = dbName
	}
}

// WithDeadline override the initial timeout for bootstrap container, by default used 5 seconds deadline.
func WithDeadline(deadline time.Duration) Option {
	return func(c *config) {
		c.deadline = deadline
	}
}

// WithPoolMaxConnections override default number of pgx pool max connections by default used  8.
func WithPoolMaxConnections(maxCons int32) Option {
	return func(c *config) {
		c.poolMaxConns = maxCons
	}
}

// WithPoolMinConnections override default number of pgx pool min connections by default used  2.
func WithPoolMinConnections(minCons int32) Option {
	return func(c *config) {
		c.poolMinConns = minCons
	}
}

// WithPoolMaxIdleTime  override default number of pgx pool min connections by default used one minute.
func WithPoolMaxIdleTime(maxIdleTime time.Duration) Option {
	return func(c *config) {
		c.poolMaxConnIdleTime = maxIdleTime
	}
}

// WithMigrator override default migration func for pgx pool by default used plain sql files in directory.
func WithMigrator(migrator Migrator) Option {
	return func(c *config) {
		c.migrator = migrator
		c.hasSetMigrator = true
	}
}

func WithFileSystem(fs afero.Fs) Option {
	return func(c *config) {
		c.fs = fs
	}
}

// WithMigrationsPath set default path to migrations files.
func WithMigrationsPath(path string) Option {
	return func(c *config) {
		c.migrationsPath = path
	}
}

func WithPoolInjectLabel(label string) Option {
	return func(c *config) {
		c.injectPoolLabel = label
	}
}

func WithPoolConfigInjectLabel(label string) Option {
	return func(c *config) {
		c.injectConfigLabel = label
	}
}
