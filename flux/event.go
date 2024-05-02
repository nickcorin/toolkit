package flux

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nickcorin/toolkit/sqlkit"
)

// ErrEventNotFound is an error that is returned when an event cannot be found.
var ErrEventNotFound = errors.New("event not found")

type Event[TOPIC, SEQ comparable] interface {
	// ID returns a unique identifier for the event.
	ID() string

	// Topic is a partitioning key that can be used to filter events.
	Topic() TOPIC

	// Sequence returns an unsigned integer that can be used to order events. Events with a lower sequence number should
	// be considered to have been emitted before those with a higher sequence number.
	Sequence() SEQ

	// Key is a reference to the entity that the event is related to.
	Key() string

	// Timestamp returns the time at which the event was emitted.
	Timestamp() time.Time
}

// Compile-time assertion that defaultEvent[TOPIC, SEQ] implements the Event[TOPIC, SEQ] interface.
var _ Event[string, uint] = (*defaultEvent[string, uint])(nil)

type defaultEvent[TOPIC, SEQ comparable] struct {
	id        string
	topic     TOPIC
	sequence  SEQ
	key       string
	timestamp time.Time
}

func (event *defaultEvent[TOPIC, SEQ]) ID() string           { return event.id }
func (event *defaultEvent[TOPIC, SEQ]) Topic() TOPIC         { return event.topic }
func (event *defaultEvent[TOPIC, SEQ]) Sequence() SEQ        { return event.sequence }
func (event *defaultEvent[TOPIC, SEQ]) Key() string          { return event.key }
func (event *defaultEvent[TOPIC, SEQ]) Timestamp() time.Time { return event.timestamp }

// EventStore is an interface that combines an EventReader and an EventWriter.
type EventStore[TOPIC, SEQ comparable] interface {
	EventReader[TOPIC, SEQ]
	EventWriter[TOPIC, SEQ]
}

// EventReader allows read-only access to an event store.
type EventReader[TOPIC, SEQ comparable] interface {
	NextEvents(ctx context.Context, topic TOPIC, from SEQ, batchSize uint) ([]Event[TOPIC, SEQ], error)
}

// EventWriter allows write-only access to an event store.
type EventWriter[TOPIC, SEQ comparable] interface {
	CreateEvent(ctx context.Context, topic TOPIC, key string) (Event[TOPIC, SEQ], error)
}

// NewPostgresEventStore creates a new instance of a PostgresEventStore.
func NewPostgresEventStore[TOPIC, SEQ comparable](conn *sql.DB, tableName string) *PostgresEventStore[TOPIC, SEQ] {
	return &PostgresEventStore[TOPIC, SEQ]{conn: conn, tableName: tableName}
}

// PostgresEventStore is an implementation of an EventStore that uses a PostgreSQL database as its storage backend.
type PostgresEventStore[TOPIC, SEQ comparable] struct {
	conn      *sql.DB
	tableName string
}

func (store *PostgresEventStore[TOPIC, SEQ]) CreateEvent(ctx context.Context, topic TOPIC, key string) (Event[TOPIC, SEQ], error) {
	query := `
	INSERT INTO ` + store.tableName + ` (id, key, topic, timestamp)
	VALUES (uuid_generate_v4(), $1, $2, $3)
	RETURNING id, topic, sequence, key, timestamp`

	return scanEvent[TOPIC, SEQ](store.conn.QueryRowContext(ctx, query, key, topic, time.Now().UTC()))
}

func (store *PostgresEventStore[TOPIC, SEQ]) NextEvents(ctx context.Context, topic TOPIC, from SEQ, batchSize uint) ([]Event[TOPIC, SEQ], error) {
	query := "SELECT * FROM " + store.tableName + " WHERE sequence > $1 ORDER BY sequence ASC LIMIT $2"

	rows, err := store.conn.QueryContext(ctx, query, from, batchSize)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	var events []Event[TOPIC, SEQ]
	for rows.Next() {
		event, err := scanEvent[TOPIC, SEQ](rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return nil, ErrEventNotFound
	}

	return events, nil
}

func scanEvent[TOPIC, SEQ comparable](s sqlkit.Scannable) (*defaultEvent[TOPIC, SEQ], error) {
	var event defaultEvent[TOPIC, SEQ]

	err := s.Scan(&event.id, &event.topic, &event.sequence, &event.key, &event.timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventNotFound
		}

		return nil, fmt.Errorf("scan event: %w", err)
	}

	return &event, nil
}
