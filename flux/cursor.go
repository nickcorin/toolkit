package flux

import (
	"context"
	"errors"
	"time"
)

// ErrCursorNotFound is an error that is returned when a cursor cannot be found.
var ErrCursorNotFound = errors.New("cursor not found")

// Cursor is a "bookmark" that keeps track of the most recent event that a Consumer has processed.
type Cursor interface {
	// ID returns a unique identifier for the cursor.
	ID() string

	// Name is a non-generated identifier for the cursor that can be used to lookup the cursor. This should be unique
	// across all cursors.
	Name() string

	// Sequence returns a value which refers to the most recent event that a Consumer has processed.
	Sequence() uint

	// CreatedAt returns the time at which the cursor was created.
	CreatedAt() time.Time

	// UpdatedAt returns the time at which the cursor was last updated.
	UpdatedAt() time.Time
}

type defaultCursor struct {
	id        string
	name      string
	sequence  uint
	createdAt time.Time
	updatedAt time.Time
}

func (cursor *defaultCursor) ID() string           { return cursor.id }
func (cursor *defaultCursor) Name() string         { return cursor.name }
func (cursor *defaultCursor) Sequence() uint       { return cursor.sequence }
func (cursor *defaultCursor) CreatedAt() time.Time { return cursor.createdAt }
func (cursor *defaultCursor) UpdatedAt() time.Time { return cursor.updatedAt }

// CursorStore is an interface that allows the creation and lookup of cursors.
type CursorStore interface {
	CursorReader
	CursorWriter
}

// CursorReader allows read-only access to a cursor store.
type CursorReader interface {
	LookupCursorByID(ctx context.Context, id string) (Cursor, error)
	LookupCursorByName(ctx context.Context, name string) (Cursor, error)
}

// CursorWriter allows write-only access to a cursor store.
type CursorWriter interface {
	CreateCursor(ctx context.Context, name string, sequence uint) (Cursor, error)
	UpdateCursor(ctx context.Context, id string, sequence uint) error
}
