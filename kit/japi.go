package kit

import (
	"fmt"
	"j/japi/core"
	"j/japi/jerr"
	"j/schema/bytes"
	"j/schema/fs"
	"j/schema/reader"
)

// JApi is an interface-level wrapper for JApiCore
type JApi struct {
	core *core.JApiCore
}

// NewJapi returns interface-level wrapper for JApiCore
// Does not include .jst file validation. File validation should be called explicitly.
func NewJapi(filepath string) (JApi, error) {
	f, err := readPanicFree(filepath)
	if err != nil {
		return JApi{}, err
	}
	return JApi{core.NewJApiCore(f)}, nil
}

func readPanicFree(filename string) (f *fs.File, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	f = reader.Read(filename)
	return f, err
}

func NewJapiFromBytes(b bytes.Bytes) JApi {
	return JApi{
		core.NewJApiCore(fs.NewFile("root", b)),
	}
}

// ValidateJAPI validates .jst file
func (j *JApi) ValidateJAPI() *jerr.JAPIError {
	return j.core.ValidateJAPI()
}

func (j JApi) ToJson() ([]byte, error) {
	c := j.core.Catalog()
	return c.ToJson()
}

func (j JApi) Title() string {
	if j.core != nil && j.core.Catalog() != nil && j.core.Catalog().Info != nil {
		return j.core.Catalog().Info.Title
	}
	return ""
}

func (j JApi) ToJsonIndent() ([]byte, error) {
	c := j.core.Catalog()
	return c.ToJsonIndent()
}
