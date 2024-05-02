package sqlkit

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectForTesting(t *testing.T, dialect Dialect, migrationsFS embed.FS) (*sql.DB, error) {
	t.Helper()

	connector, err := GetConnector(dialect)
	if err != nil {
		return nil, fmt.Errorf("get connector: %w", err)
	}

	conf := connector.Defaults()
	conf.Flags.Set("sslmode", "disable")

	// Create a connection to the default database.
	conn, err := Connect(context.Background(), dialect, conf)
	if err != nil {
		return nil, fmt.Errorf("connect to default database: %w", err)
	}

	seedDB := fmt.Sprintf("test_%d", time.Now().UnixNano())

	if err := createDatabase(t, conn, seedDB); err != nil {
		return nil, fmt.Errorf("create seed database: %w", err)
	}

	migrationsPath, err := findMigrationsPath(t, migrationsFS, ".")
	if err != nil {
		return nil, fmt.Errorf("find migrations path: %w", err)
	}

	conf.Database = seedDB
	source, err := iofs.New(migrationsFS, migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("create migrations source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, connector.DSN(conf))
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("migrate up: %w", err)
	}

	seedConn, err := Connect(context.Background(), dialect, conf)
	if err != nil {
		return nil, fmt.Errorf("connect to seed database: %w", err)
	}

	t.Cleanup(func() {
		// Close the connection to the seed database.
		_ = seedConn.Close()

		_ = migrator.Down()
		_, _ = migrator.Close()

		_ = dropDatabase(t, conn, seedDB)
		_ = conn.Close()
	})

	return seedConn, nil
}

func createDatabase(t *testing.T, conn *sql.DB, name string) error {
	t.Helper()

	_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	if err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	return nil
}

func dropDatabase(t *testing.T, conn *sql.DB, name string) error {
	t.Helper()

	_, err := conn.Exec(fmt.Sprintf("DROP DATABASE %s", name))
	if err != nil {
		return fmt.Errorf("drop database: %w", err)
	}

	return nil
}

func findMigrationsPath(t *testing.T, migrations embed.FS, startDir string) (string, error) {
	t.Helper()

	dirQueue := []string{startDir}

	for len(dirQueue) > 0 {
		dir := dirQueue[0]
		dirQueue = dirQueue[1:]

		entries, err := migrations.ReadDir(dir)
		if err != nil {
			return "", fmt.Errorf("read dir: %w", err)
		}

		for _, entry := range entries {
			// All files are read relative to the root, we need to keep track of the full path of the files here.
			path := filepath.Join(dir, entry.Name())

			if entry.IsDir() {
				// Store the directory for later, and move on.
				dirQueue = append(dirQueue, path)

				continue
			}

			if strings.HasSuffix(entry.Name(), ".sql") {
				// We want SQL files.
				return dir, nil
			}
		}
	}

	return "", fmt.Errorf("no migrations found")
}
