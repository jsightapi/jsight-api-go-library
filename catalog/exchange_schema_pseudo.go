package catalog

import (
	"encoding/json"
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/notation"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
)

type ExchangePseudoSchema struct {
	jschemaLib.Schema
	Notation notation.SchemaNotation
}

func NewExchangePseudoSchema(n notation.SchemaNotation) *ExchangePseudoSchema {
	return &ExchangePseudoSchema{
		Notation: n,
	}
}

func (e ExchangePseudoSchema) MarshalJSON() ([]byte, error) {
	if e.Notation != notation.SchemaNotationAny && e.Notation != notation.SchemaNotationEmpty {
		return nil, fmt.Errorf(`invalid schema notation "%s"`, e.Notation)
	}

	data := struct {
		Notation notation.SchemaNotation `json:"notation"`
	}{
		Notation: e.Notation,
	}

	return json.Marshal(data)
}
