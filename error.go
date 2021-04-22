package vector

import "errors"

var (
	ErrEmptySrc     = errors.New("can't parse empty source")
	ErrNotImplement = errors.New("parser not implemented")
	ErrIncompatType = errors.New("incompatible type")
	ErrNotFound     = errors.New("node not found")
	ErrInternal     = errors.New("internal vector error")
)
