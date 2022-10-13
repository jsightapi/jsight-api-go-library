package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
)

type UserType struct {
	Annotation  string              `json:"annotation,omitempty"`
	Description string              `json:"description,omitempty"`
	Schema      jschemaLib.Schema   `json:"schema"`
	Directive   directive.Directive `json:"-"`
}
