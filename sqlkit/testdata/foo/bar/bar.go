package bar

import (
	"database/sql"
	"time"

	"github.com/nickcorin/toolkit/sqlkit/testdata/foo"
)

type bar struct {
	foo.Foo

	C int
	D time.Time `sqlkit:"d_override"`
	E sql.Null[time.Time]
}
