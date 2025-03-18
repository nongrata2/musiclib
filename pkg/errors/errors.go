package errors

import (
	"errors"
)

var (
	NotFoundErr   = errors.New("no song found with the given ID")
	OutOfRangeErr = errors.New("page out of range")
)
