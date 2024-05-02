package flux

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nickcorin/toolkit/sqlkit"
)

// ErrCursorNotFound is an error that is returned when a cursor cannot be found.
var ErrCursorNotFound = errors.New("cursor not found")

// Cursor is a "bookmark" that keeps track of the most recent event that a Consumer has processed.
type Cursor[SEQ comparable] interface {
	// ID returns a unique identifier for the cursor.
	ID() string

	// Name is a non-generated identifier for the cursor that can be used to lookup the cursor. This should be unique
	// across cursors.
	Name() string

	// Sequence returns a value which refers to the most recent event that a Consumer has processed.
	Sequence() SEQ

	// CreatedAt returns the time at which the cursor was created.
	CreatedAt() time.Time

	// UpdatedAt returns the time at which the cursor was last updated.
	UpdatedAt() time.Time
}

type defaultCursor[SEQ comparable] struct {
	id        string
	name      string
	sequence  SEQ
	createdAt time.Time
	updatedAt time.Time
}

func (cursor *defaultCursor[SEQ]) ID() string           { return cursor.id }
func (cursor *defaultCursor[SEQ]) Name() string         { return cursor.name }
func (cursor *defaultCursor[SEQ]) Sequence() SEQ        { return cursor.sequence }
func (cursor *defaultCursor[SEQ]) CreatedAt() time.Time { return cursor.createdAt }
func (cursor *defaultCursor[SEQ]) UpdatedAt() time.Time { return cursor.updatedAt }

// CursorStore is an interface that allows the creation and lookup of cursors.
type CursorStore[SEQ comparable] interface {
	CreateCursor(ctx context.Context, name string, sequence SEQ) (Cursor[SEQ], error)
	LookupCursorByID(ctx context.Context, id string) (Cursor[SEQ], error)
	LookupCursorByName(ctx context.Context, name string) (Cursor[SEQ], error)
}

// NewPostgresCursorStore returns a new instance of PostgresCursorStore.
func NewPostgresCursorStore[SEQ comparable](conn *sql.DB, tableName string) *PostgresCursorStore[SEQ] {
	return &PostgresCursorStore[SEQ]{conn: conn, tableName: tableName}
}

// Compile-time assertion that PostgresCursorStore[int] implements the CursorStore[int] interface.
var _ CursorStore[int] = (*PostgresCursorStore[int])(nil)

// PostgresCursorStore is a CursorStore that uses a PostgreSQL database as its storage backend.
type PostgresCursorStore[SEQ comparable] struct {
	conn      *sql.DB
	tableName string
}

func (store *PostgresCursorStore[SEQ]) CreateCursor(ctx context.Context, name string, sequence SEQ) (Cursor[SEQ], error) {
	query := `
	INSERT INTO ` + store.tableName + ` (id, name, sequence, created_at, updated_at)
	VALUES (uuid_generate_v4(), $1, $2, $3, $4)
	RETURNING id, name, sequence, created_at, updated_at`

	return scanCursor[SEQ](store.conn.QueryRowContext(ctx, query, name, sequence, time.Now().UTC(), time.Now().UTC()))
}

func (store *PostgresCursorStore[SEQ]) LookupCursorByID(ctx context.Context, id string) (Cursor[SEQ], error) {
	query := "SELECT * FROM " + store.tableName + " WHERE id = $1"
	return scanCursor[SEQ](store.conn.QueryRowContext(ctx, query, id))
}

func (store *PostgresCursorStore[SEQ]) LookupCursorByName(ctx context.Context, name string) (Cursor[SEQ], error) {
	query := "SELECT * FROM " + store.tableName + " WHERE name = $1"
	return scanCursor[SEQ](store.conn.QueryRowContext(ctx, query, name))
}

func scanCursor[SEQ comparable](s sqlkit.Scannable) (*defaultCursor[SEQ], error) {
	var cursor defaultCursor[SEQ]

	err := s.Scan(&cursor.id, &cursor.name, &cursor.sequence, &cursor.createdAt, &cursor.updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCursorNotFound
		}

		return nil, fmt.Errorf("scan cursor: %w", err)
	}

	return &cursor, nil
}
