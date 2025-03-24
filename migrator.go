package pgrx

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

const defaultExpMigrations = 8

func PlainMigrator(fs afero.Fs, path string) (Migrator, error) {
	migrations := make([]string, 0, defaultExpMigrations)

	dir, err := fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open migration dir: %w", err)
	}
	defer dir.Close()

	list, err := dir.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("can't read migrations file list: %w", err)
	}

	for _, info := range list {
		if info.IsDir() {
			continue
		}
		fh, err := fs.Open(filepath.Join(path, info.Name()))
		if err != nil {
			return nil, fmt.Errorf("can't open migration file %s: %w", info.Name(), err)
		}
		defer fh.Close()

		data, err := io.ReadAll(fh)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, string(data))
	}

	return func(ctx context.Context, cfg MigratorConfig) error {
		for _, migration := range migrations {
			if _, err := cfg.Pool.Exec(ctx, migration); err != nil {
				return err
			}
		}
		return nil
	}, nil
}
