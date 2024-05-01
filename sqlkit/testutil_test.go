package sqlkit_test

import (
	"context"
	"embed"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/migrations/postgres/*.sql
var pgMigrations embed.FS

func TestConnectForTesting_Postgres(t *testing.T) {
	testSchemas := make([]string, 0)

	t.Run("create seed database", func(t *testing.T) {
		conn, err := sqlkit.ConnectForTesting(t, sqlkit.Postgres, pgMigrations)
		require.NoError(t, err)
		require.NotNil(t, conn)

		var version string
		err = conn.QueryRowContext(context.Background(), "SELECT version()").Scan(&version)
		require.NoError(t, err)
		require.NotEmpty(t, version)

		var dbName string
		err = conn.QueryRowContext(context.Background(), "SELECT current_database()").Scan(&dbName)
		require.NoError(t, err)
		require.NotEmpty(t, dbName)

		testSchemas = append(testSchemas, dbName)
	})

	for _, schema := range testSchemas {
		s := schema // capture range variable.

		t.Run("ensure test db was dropped: "+schema, func(t *testing.T) {
			connector, err := sqlkit.GetConnector(sqlkit.Postgres)
			require.NoError(t, err)
			require.NotNil(t, connector)

			conf := connector.Defaults()
			conf.Database = s

			conn, err := sqlkit.Connect(context.Background(), sqlkit.Postgres, conf)
			require.Empty(t, conn)
			require.Error(t, err)
			require.True(t, sqlkit.PgErrorIs(err, pgerrcode.InvalidCatalogName))
		})
	}
}
