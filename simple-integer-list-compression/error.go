package simple

import "errors"

var (
	ErrCardinalityHeaderSizeOutOfBound = errors.New("simple: CardinalityHeaderSize out of bound")
	ErrInputTooLong                    = errors.New("simple: input too long")
)
