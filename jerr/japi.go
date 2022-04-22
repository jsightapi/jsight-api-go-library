package jerr

import (
	"j/schema/bytes"
	"j/schema/fs"
)

type JAPIError struct {
	Location
	Msg string
}

func NewJAPIError(msg string, f *fs.File, i bytes.Index) *JAPIError {
	loc := NewLocation(f, i)
	return &JAPIError{Location: loc, Msg: msg}
}

func (e JAPIError) Error() string {
	return e.Msg
}
