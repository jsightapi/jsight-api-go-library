package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type HTTPRequest struct {
	*HTTPRequestHeaders `json:"headers,omitempty"`
	*HTTPRequestBody    `json:"body,omitempty"`
	Directive           directive.Directive `json:"-"`
}

type HTTPRequestHeaders struct {
	Schema    *jschema.Schema     `json:"schema"`
	Directive directive.Directive `json:"-"`
}

type HTTPRequestBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    jschemaLib.Schema   `json:"schema"`
	Directive directive.Directive `json:"-"`
}
