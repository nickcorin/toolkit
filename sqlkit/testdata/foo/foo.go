package foo

import (
	"database/sql"
	"time"
)

type FooType struct {
	A string
	B int
	C BazType
	D time.Time
	E time.Time
	F sql.NullTime
	G QuxType
}

type BazType int
type QuxType struct{}
