package core

import (
	"j/japi/catalog"
	"j/japi/directive"
)

type rawPathVariable struct {
	pathDirective   directive.Directive // to detect and display an error
	parentDirective directive.Directive
	schema          catalog.Schema
	parameters      []PathParameter
}
