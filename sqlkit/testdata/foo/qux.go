package foo

import "time"

type Qux struct {
	A string
	B int
	C time.Time
}

type gen struct {
	Qux
}
