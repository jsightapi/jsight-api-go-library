package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type HTTPResponse struct {
	Code       string               `json:"code"`
	Annotation string               `json:"annotation,omitempty"`
	Headers    *HTTPResponseHeaders `json:"headers,omitempty"`
	Body       *HTTPResponseBody    `json:"body"`
	Directive  directive.Directive  `json:"-"`
}

type HTTPResponseHeaders struct {
	Schema    *jschema.Schema     `json:"schema"`
	Directive directive.Directive `json:"-"`
}
