package flux

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nickcorin/toolkit/sqlkit"
)

// NewPostgresCursorStore returns a new instance of PostgresCursorStore.
func NewPostgresCursorStore(conn *sql.DB, tableName string) *PostgresCursorStore {
	return &PostgresCursorStore{conn: conn, tableName: tableName}
}

// Compile-time assertion that PostgresCursorStore implements the CursorStore interface.
var _ CursorStore = (*PostgresCursorStore)(nil)

// PostgresCursorStore is a CursorStore that uses a PostgreSQL database as its storage backend.
type PostgresCursorStore struct {
	conn      *sql.DB
	tableName string
}

func (store *PostgresCursorStore) CreateCursor(ctx context.Context, name string, sequence uint) (Cursor, error) {
	query := `
	INSERT INTO ` + store.tableName + ` (id, name, sequence, created_at, updated_at)
	VALUES (uuid_generate_v4(), $1, $2, $3, $4)
	RETURNING id, name, sequence, created_at, updated_at`

	return scanCursor(store.conn.QueryRowContext(ctx, query, name, sequence, time.Now().UTC(), time.Now().UTC()))
}

func (store *PostgresCursorStore) LookupCursorByID(ctx context.Context, id string) (Cursor, error) {
	query := "SELECT * FROM " + store.tableName + " WHERE id = $1"
	return scanCursor(store.conn.QueryRowContext(ctx, query, id))
}

func (store *PostgresCursorStore) LookupCursorByName(ctx context.Context, name string) (Cursor, error) {
	query := "SELECT * FROM " + store.tableName + " WHERE name = $1"
	return scanCursor(store.conn.QueryRowContext(ctx, query, name))
}

func (store *PostgresCursorStore) UpdateCursor(ctx context.Context, id string, sequence uint) error {
	query := "UPDATE " + store.tableName + " SET sequence = $1, updated_at = $2 WHERE id = $3"

	_, err := store.conn.ExecContext(ctx, query, sequence, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("update cursor: %w", err)
	}

	return nil
}

func scanCursor(s sqlkit.Scannable) (*defaultCursor, error) {
	var cursor defaultCursor

	err := s.Scan(&cursor.id, &cursor.name, &cursor.sequence, &cursor.createdAt, &cursor.updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCursorNotFound
		}

		return nil, fmt.Errorf("scan cursor: %w", err)
	}

	return &cursor, nil
}

// NewPostgresEventStore creates a new instance of a PostgresEventStore.
func NewPostgresEventStore(conn *sql.DB, tableName string) *PostgresEventStore {
	return &PostgresEventStore{conn: conn, tableName: tableName}
}

// PostgresEventStore is an implementation of an EventStore that uses a PostgreSQL database as its storage backend.
type PostgresEventStore struct {
	conn      *sql.DB
	tableName string
}

func (store *PostgresEventStore) CreateEvent(ctx context.Context, topic, key string) (Event, error) {
	query := `
	INSERT INTO ` + store.tableName + ` (id, key, topic, timestamp)
	VALUES (uuid_generate_v4(), $1, $2, $3)
	RETURNING id, topic, sequence, key, timestamp`

	return scanEvent(store.conn.QueryRowContext(ctx, query, key, topic, time.Now().UTC()))
}

func (store *PostgresEventStore) Head(ctx context.Context) (Event, error) {
	query := "SELECT * FROM " + store.tableName + " ORDER BY sequence DESC LIMIT 1"
	return scanEvent(store.conn.QueryRowContext(ctx, query))
}

func (store *PostgresEventStore) NextEvents(
	ctx context.Context,
	from uint,
	batchSize uint,
	streamLag time.Duration,
) ([]Event, error) {
	query := `
	SELECT * FROM ` + store.tableName + ` WHERE sequence > $1 AND timestamp < $2 ORDER BY sequence ASC LIMIT $3
	`

	rows, err := store.conn.QueryContext(ctx, query, from, time.Now().Add(-1*streamLag), batchSize)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	var events []Event
	for rows.Next() {
		event, err := scanEvent(rows)
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

func scanEvent(s sqlkit.Scannable) (*defaultEvent, error) {
	var event defaultEvent

	err := s.Scan(&event.id, &event.topic, &event.sequence, &event.key, &event.timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventNotFound
		}

		return nil, fmt.Errorf("scan event: %w", err)
	}

	return &event, nil
}
