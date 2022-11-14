package catalog

import (
	"github.com/jsightapi/jsight-schema-go-library/kit"
)

type PathVariables struct {
	Schema *ExchangeJSightSchema `json:"schema"`
}

func (p *PathVariables) Validate(key, value string) kit.Error {
	return p.Schema.ValidateObjectProperty(key, value)
}
