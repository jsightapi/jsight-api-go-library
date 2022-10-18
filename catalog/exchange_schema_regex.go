package catalog

import (
	"encoding/json"
	"github.com/jsightapi/jsight-api-go-library/notation"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/notations/regex"
)

type ExchangeRegexSchema struct {
	*regex.Schema
}

func (e ExchangeRegexSchema) MarshalJSON() ([]byte, error) {
	data := struct {
		Content  interface{}             `json:"content,omitempty"`
		Example  string                  `json:"example,omitempty"`
		Notation notation.SchemaNotation `json:"notation"`
	}{
		Notation: notation.SchemaNotationRegex,
	}

	data.Content, _ = e.Pattern()

	ex, _ := e.Example()
	data.Example = string(ex)

	return json.Marshal(data)
}

func PrepareRegexSchema(name string, regexStr bytes.Bytes) (*ExchangeRegexSchema, error) {
	var oo []regex.Option
	oo = append(oo, regex.WithGeneratorSeed(0))

	s := regex.New(name, regexStr, oo...)

	// n, err := s.GetAST()
	// if err != nil {
	// 	return Schema{}, err
	// }

	return &ExchangeRegexSchema{Schema: s}, nil
}
