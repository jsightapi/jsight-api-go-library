package core

import (
	"j/japi/catalog"
	"j/japi/directive"
)

type rawPathVariable struct {
	schema          catalog.Schema
	parameters      []PathParameter
	pathDirective   directive.Directive // to detect and display an error
	parentDirective directive.Directive
}
