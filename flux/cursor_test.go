package flux_test

import (
	"context"
	"testing"

	"github.com/nickcorin/toolkit/flux"
	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

func TestPostgresCursorStore(t *testing.T) {
	conn, err := sqlkit.ConnectForTesting(t, sqlkit.Postgres, pgMigrations)
	require.NoError(t, err)
	require.NotNil(t, conn)

	testCursorStore(t, flux.NewPostgresCursorStore[uint](conn, "cursors"))
}

func testCursorStore(t *testing.T, cursorStore flux.CursorStore[uint]) {
	cursor, err := cursorStore.CreateCursor(context.Background(), "cursors", 0)
	require.NoError(t, err)
	require.NotNil(t, cursor)

	cursor2, err := cursorStore.LookupCursorByID(context.Background(), cursor.ID())
	require.NoError(t, err)
	require.NotNil(t, cursor2)

	cursor3, err := cursorStore.LookupCursorByName(context.Background(), cursor.Name())
	require.NoError(t, err)
	require.NotNil(t, cursor3)

	require.EqualValues(t, cursor.ID(), cursor2.ID())
	require.EqualValues(t, cursor2.ID(), cursor3.ID())

	require.EqualValues(t, cursor.Name(), cursor2.Name())
	require.EqualValues(t, cursor2.Name(), cursor3.Name())

	require.EqualValues(t, cursor.Sequence(), cursor2.Sequence())
	require.EqualValues(t, cursor2.Sequence(), cursor3.Sequence())

	require.EqualValues(t, cursor.CreatedAt(), cursor2.CreatedAt())
	require.EqualValues(t, cursor2.CreatedAt(), cursor3.CreatedAt())

	require.EqualValues(t, cursor.UpdatedAt(), cursor2.UpdatedAt())
	require.EqualValues(t, cursor2.UpdatedAt(), cursor3.UpdatedAt())
}
