package core

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type rawPathVariable struct {
	schema          *jschema.JSchema
	parameters      []PathParameter
	pathDirective   directive.Directive // to detect and display an error
	parentDirective directive.Directive
}
