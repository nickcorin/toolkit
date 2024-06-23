package flux_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nickcorin/toolkit/flux"
	"github.com/nickcorin/toolkit/sqlkit"
	"github.com/stretchr/testify/require"
)

func TestPostgresCursorStore(t *testing.T) {
	conn, err := sqlkit.ConnectForTesting(t, sqlkit.Postgres, pgMigrations)
	require.NoError(t, err)
	require.NotNil(t, conn)

	testCursorStore(t, flux.NewPostgresCursorStore(conn, "cursors"))
}

func testCursorStore(t *testing.T, cursorStore flux.CursorStore) {
	cursorName := "test-cursor"

	cursor, err := cursorStore.CreateCursor(context.Background(), cursorName, 0)
	require.NoError(t, err)
	require.NotNil(t, cursor)

	cursor2, err := cursorStore.LookupCursorByID(context.Background(), cursor.ID())
	require.NoError(t, err)
	require.NotNil(t, cursor2)

	cursor3, err := cursorStore.LookupCursorByName(context.Background(), cursor.Name())
	require.NoError(t, err)
	require.NotNil(t, cursor3)

	require.Equal(t, cursor.ID(), cursor2.ID())
	require.Equal(t, cursor2.ID(), cursor3.ID())

	require.Equal(t, cursor.Name(), cursor2.Name())
	require.Equal(t, cursor2.Name(), cursor3.Name())

	require.Equal(t, cursor.Sequence(), cursor2.Sequence())
	require.Equal(t, cursor2.Sequence(), cursor3.Sequence())

	require.Equal(t, cursor.CreatedAt(), cursor2.CreatedAt())
	require.Equal(t, cursor2.CreatedAt(), cursor3.CreatedAt())

	require.Equal(t, cursor.UpdatedAt(), cursor2.UpdatedAt())
	require.Equal(t, cursor2.UpdatedAt(), cursor3.UpdatedAt())

	require.Equal(t, cursorName, cursor.Name())
	require.Equal(t, uint(0), cursor.Sequence())

	err = cursorStore.UpdateCursor(context.Background(), cursor.ID(), 1)
	require.NoError(t, err)

	cursor4, err := cursorStore.LookupCursorByID(context.Background(), cursor.ID())
	require.NoError(t, err)
	require.NotNil(t, cursor4)

	require.Equal(t, uint(1), cursor4.Sequence())
}

func TestPostgresEventStore(t *testing.T) {
	conn, err := sqlkit.ConnectForTesting(t, sqlkit.Postgres, pgMigrations)
	require.NoError(t, err)
	require.NotNil(t, conn)

	testEventStore(t, flux.NewPostgresEventStore(conn, "events"))
}

func testEventStore(t *testing.T, eventStore flux.EventStore) {
	var (
		topicA = flux.EventTopic(uuid.NewString())
		topicB = flux.EventTopic(uuid.NewString())
		topicC = flux.EventTopic(uuid.NewString())

		keyA = uuid.NewString()
		keyB = uuid.NewString()
		keyC = uuid.NewString()
	)

	testEvents := []struct {
		topic flux.EventTopic
		key   string
	}{
		{topicA, keyA},
		{topicB, keyB},
		{topicC, keyC},
	}

	storedEvents := make([]flux.Event, 0)

	t.Run("query an empty outbox", func(t *testing.T) {
		event, err := eventStore.Head(context.Background())
		require.ErrorIs(t, err, flux.ErrEventNotFound)
		require.Empty(t, event)

		events, err := eventStore.NextEvents(context.Background(), 0, 5, 0)
		require.Empty(t, events)
		require.ErrorIs(t, err, flux.ErrEventNotFound)
	})

	t.Run("create events", func(t *testing.T) {
		for _, testEvent := range testEvents {
			event, err := eventStore.CreateEvent(context.Background(), testEvent.topic.String(), testEvent.key)
			require.NoError(t, err)
			require.NotNil(t, event)

			storedEvents = append(storedEvents, event)
		}
	})

	t.Run("fetch a full batch of events", func(t *testing.T) {
		const batchSize = 2

		events, err := eventStore.NextEvents(context.Background(), 0, batchSize, 0)
		require.NoError(t, err)
		require.Len(t, events, batchSize)
	})

	t.Run("fetch a partially filled batch of events", func(t *testing.T) {
		const batchSize = 6

		events, err := eventStore.NextEvents(context.Background(), 0, batchSize, 0)
		require.NoError(t, err)
		require.Len(t, events, len(storedEvents))
	})

	t.Run("fetch the head of the queue", func(t *testing.T) {
		event, err := eventStore.Head(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, event)

		require.Equal(t, event, storedEvents[len(storedEvents)-1])
	})
}
