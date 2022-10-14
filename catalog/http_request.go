package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type HTTPRequest struct {
	*HTTPRequestHeaders `json:"headers,omitempty"`
	*HTTPRequestBody    `json:"body,omitempty"`
	Directive           directive.Directive `json:"-"`
}

type HTTPRequestHeaders struct {
	Schema    *ExchangeJSightSchema `json:"schema"`
	Directive directive.Directive   `json:"-"`
}

type HTTPRequestBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    ExchangeSchema      `json:"schema"`
	Directive directive.Directive `json:"-"`
}
