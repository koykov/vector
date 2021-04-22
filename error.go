package vector

import "errors"

var (
	ErrNotImplement = errors.New("parser not implemented")
	ErrIncompatType = errors.New("incompatible type")
	ErrNotFound     = errors.New("node not found")
	ErrInternal     = errors.New("internal vector error")
)
