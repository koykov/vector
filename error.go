package vector

import "errors"

var (
	ErrEmptySrc     = errors.New("can't parse empty source")
	ErrShortSrc     = errors.New("source is too short to parse")
	ErrNotImplement = errors.New("method not implemented")
	ErrIncompatType = errors.New("incompatible type")
	ErrNotFound     = errors.New("node not found")
	ErrInternal     = errors.New("internal vector error")
	ErrUnparsedTail = errors.New("unparsed tail")
	ErrUnexpId      = errors.New("unexpected identifier")
	ErrUnexpEOF     = errors.New("unexpected end of file")
	ErrUnexpEOS     = errors.New("unexpected end of string")

	_, _, _, _, _ = ErrShortSrc, ErrUnparsedTail, ErrUnexpId, ErrUnexpEOF, ErrUnexpEOS
)
