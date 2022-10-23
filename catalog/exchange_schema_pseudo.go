package catalog

import (
	"encoding/json"
	"fmt"
	"github.com/jsightapi/jsight-api-go-library/notation"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
)

type ExchangePseudoSchema struct {
	jschemaLib.Schema
	notation notation.SchemaNotation
}

func NewExchangePseudoSchema(n notation.SchemaNotation) *ExchangePseudoSchema {
	return &ExchangePseudoSchema{
		notation: n,
	}
}

func (e ExchangePseudoSchema) MarshalJSON() ([]byte, error) {
	if e.notation != notation.SchemaNotationAny && e.notation != notation.SchemaNotationEmpty {
		return nil, fmt.Errorf(`invalid schema notation "%s"`, e.notation)
	}

	data := struct {
		Notation notation.SchemaNotation `json:"notation"`
	}{
		Notation: e.notation,
	}

	return json.Marshal(data)
}
