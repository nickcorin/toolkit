package bar

import (
	"database/sql"
	"time"

	"github.com/nickcorin/toolkit/sqlkit/testdata/foo"
)

type bar struct {
	foo.Foo

	C int
	D sql.NullTime `sqlkit:"d_override"`
	E sql.Null[time.Time]
}
