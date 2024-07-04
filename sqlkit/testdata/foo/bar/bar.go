package bar

import (
	"database/sql"
	"time"

	"github.com/nickcorin/toolkit/sqlkit/testdata/foo"
)

type bar struct {
	foo.Foo

	B int                 `sqlkit:"-"` // omit this field.
	C int                 // override the type with the underlying primitive.
	D sql.NullTime        `sqlkit:"d_override"` // override the column name.
	E sql.Null[time.Time] // override the type with a generic.
}
