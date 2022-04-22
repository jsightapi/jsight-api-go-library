package core

import (
	"j/japi/jerr"
	"j/schema/bytes"
)

func (core *JApiCore) japiError(msg string, i bytes.Index) *jerr.JAPIError {
	return jerr.NewJAPIError(msg, core.file, i)
}
