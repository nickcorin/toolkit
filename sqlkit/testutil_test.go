package sqlkit_test

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
)

//go:embed testdata/migrations/postgres/*.sql
var pgMigrations embed.FS

func TestConnectForTesting_PreloadedMigrations(t *testing.T) {
	tests := []struct {
		name    string
		dialect sqlkit.Dialect
		fs      fs.FS
		err     error
	}{
		{
			name:    "postgres embedded migrations",
			dialect: sqlkit.Postgres,
			fs:      pgMigrations,
			err:     nil,
		},
		{
			name:    "postgres loaded migrations",
			dialect: sqlkit.Postgres,
			fs:      os.DirFS("testdata/migrations/postgres"),
			err:     nil,
		},
		{
			name:    "postgres loaded migrations, invalid path",
			dialect: sqlkit.Postgres,
			fs:      os.DirFS("testdata/migrations/invalid"),
			err:     fs.ErrNotExist,
		},
	}

	for _, tt := range tests {
		var schema string

		t.Run(tt.name, func(t *testing.T) {
			conn, err := sqlkit.ConnectForTesting(t, tt.dialect, tt.fs)
			if err != nil {
				if errors.Is(err, tt.err) {
					return
				}

				t.Fatalf("sqlkit.ConnectForTesting() = %v, want %v", err, tt.err)
			}

			var version string
			err = conn.QueryRowContext(context.Background(), "SELECT version()").Scan(&version)
			if err != nil {
				if errors.Is(err, tt.err) {
					return
				}

				t.Errorf("SELECT version() = %v, want %v", err, tt.err)
			}

			if version == "" {
				t.Errorf("SELECT version() = %v, want non-empty", version)
			}

			var dbName string
			err = conn.QueryRowContext(context.Background(), "SELECT current_database()").Scan(&dbName)
			if err != nil {
				if errors.Is(err, tt.err) {
					return
				}

				t.Errorf("SELECT current_database() = %v, want %v", err, tt.err)
			}

			if dbName == "" {
				t.Errorf("SELECT current_database() = %v, want non-empty", dbName)
			}

			// Capture the schema name for the cleanup test below.
			schema = dbName
		})

		t.Run(tt.name+"_cleanup", func(t *testing.T) {
			if schema == "" {
				t.Skip("no schema to cleanup")
			}

			connector, err := sqlkit.GetConnector(tt.dialect)
			if err != nil {
				if !errors.Is(err, tt.err) {
					return
				}

				t.Errorf("sqlkit.GetConnector() = %v, want %v", err, tt.err)
			}

			if connector == nil {
				t.Errorf("sqlkit.GetConnector() = %v, want non-nil", connector)
			}

			conf, err := connector.Defaults()
			if err != nil {
				if !errors.Is(err, tt.err) {
					return
				}

				t.Errorf("connector.Defaults() = %v, want %v", err, tt.err)
			}

			if conf == nil {
				t.Errorf("connector.Defaults() = %v, want non-nil", conf)
			}

			conf.Database = schema

			conn, err := sqlkit.Connect(context.Background(), conf)
			if err == nil {
				t.Errorf("sqlkit.Connect() = %v, want not-nil", err)
			}

			if conn != nil {
				t.Errorf("sqlkit.Connect() = %v, want nil", conn)
			}
		})
	}
}
