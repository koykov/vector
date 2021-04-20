package vector

import "errors"

var (
	ErrIncompatType = errors.New("incompatible type")
	ErrNotFound     = errors.New("node not found")
	ErrInternal     = errors.New("internal vector error")
)
