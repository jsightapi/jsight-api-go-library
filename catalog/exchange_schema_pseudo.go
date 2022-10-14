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

func (s ExchangePseudoSchema) MarshalJSON() ([]byte, error) {
	data := struct {
		Notation notation.SchemaNotation `json:"notation"`
	}{
		Notation: s.notation,
	}

	if s.notation != notation.SchemaNotationAny && s.notation != notation.SchemaNotationEmpty {
		return nil, fmt.Errorf(`invalid schema notation "%s"`, s.notation)
	}

	return json.Marshal(data)
}
