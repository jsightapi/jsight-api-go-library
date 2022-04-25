package directive

import (
	"errors"
	"j/japi/jerr"
	"strings"
)

func (d Directive) Path() (string, error) {
	var path string

	switch {
	case d.Type() == Url:
		path = d.Parameter("Path")

	case d.Type().IsHTTPRequestMethod():
		path = d.Parameter("Path")
		if path == "" {
			if d.Parent == nil {
				return "", errors.New(jerr.PathNotFound)
			}
			return d.Parent.Path() // Parent is the URL directive
		}

	default:
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
