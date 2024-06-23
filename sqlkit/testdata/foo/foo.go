package foo

import (
	"database/sql"
	"time"
)

type Foo struct {
	A string
	B int
	C Baz
	D time.Time
	E time.Time
	F sql.NullTime
}

type Baz int
type Qux struct{}
