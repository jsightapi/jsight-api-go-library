package directive

import (
	"errors"
	"j/japi/jerr"
)

func (d Directive) Method() (Enumeration, error) {
	if d.Type().IsHTTPRequestMethod() { //nolint:gocritic
		return d.Type(), nil
	} else if d.Parent != nil {
		return d.Parent.Method()
	} else {
		return Get, errors.New(jerr.MethodNotFound)
	}
}
