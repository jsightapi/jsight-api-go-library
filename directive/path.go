package directive

import (
	"errors"
	"j/japi/jerr"
	"strings"
)

func (d Directive) Path() (string, error) {
	var path string

	if d.Type() == Url { //nolint:gocritic
		path = d.Parameter("Path")
	} else if d.Type().IsHTTPRequestMethod() { //nolint:revive
		path = d.Parameter("Path")
		if path == "" {
			if d.Parent == nil {
				return "", errors.New(jerr.PathNotFound)
			}
			return d.Parent.Path() // Parent is the URL directive
		}
	} else {
		if d.Parent == nil {
			return "", errors.New(jerr.PathNotFound)
		}
		return d.Parent.Path()
	}

	if !strings.HasPrefix(path, "/") {
		return "", errors.New(jerr.IncorrectPath)
	}

	return path, nil
}
