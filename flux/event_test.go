package flux_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nickcorin/toolkit/flux"
	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

func TestPostgresEventStore(t *testing.T) {
	conn, err := sqlkit.ConnectForTesting(t, sqlkit.Postgres, pgMigrations)
	require.NoError(t, err)
	require.NotNil(t, conn)

	testEventStore(t, flux.NewPostgresEventStore[string, uint](conn, "events"))
}

func testEventStore(t *testing.T, eventStore flux.EventStore[string, uint]) {
	events, err := eventStore.NextEvents(context.Background(), uuid.NewString(), 0, 5)
	require.Empty(t, events)
	require.ErrorIs(t, err, flux.ErrEventNotFound)

	topic := uuid.NewString()
	foreignKey := uuid.NewString()

	event, err := eventStore.CreateEvent(context.Background(), topic, foreignKey)
	require.NoError(t, err)
	require.NotEmpty(t, event)

	require.Equal(t, topic, event.Topic())
	require.Equal(t, foreignKey, event.Key())
	require.Equal(t, uint(1), event.Sequence())
}
