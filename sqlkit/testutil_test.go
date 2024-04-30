package sqlkit_test

import (
	"embed"
	"testing"

	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/migrations/postgres/*.sql
var migrations embed.FS

func TestConnectForTesting(t *testing.T) {
	conn, err := sqlkit.ConnectForTesting(t, sqlkit.Postgres, migrations)
	require.NoError(t, err)
	require.NotNil(t, conn)
}
