package sqlkit_test

import (
	"embed"
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/migrations/postgres/*.sql
var pgMigrations embed.FS

func TestConnectForTesting_Postgres(t *testing.T) {
	conn, err := sqlkit.ConnectForTesting(t, sqlkit.Postgres, pgMigrations)
	require.NoError(t, err)
	require.NotNil(t, conn)
}
