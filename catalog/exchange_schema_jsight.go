package catalog

import (
	"encoding/json"
	"github.com/jsightapi/jsight-api-go-library/notation"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type ExchangeJSightSchema struct {
	*jschema.Schema
	Exchange Exchange
}

func NewExchangeJSightSchema() ExchangeJSightSchema {
	return ExchangeJSightSchema{
		Schema: nil,
		Exchange: Exchange{
			ContentJSight: nil,
			UsedUserTypes: &StringSet{},
			UsedUserEnums: &StringSet{},
		},
	}
}

func PrepareJSightSchema(
	name string,
	b []byte,
	userTypes *UserSchemas,
	enumRules map[string]jschemaLib.Rule,
) (*ExchangeJSightSchema, error) {
	s := NewExchangeJSightSchema()
	s.Schema = jschema.New(name, b)

	for n, v := range enumRules {
		if err := s.Schema.AddRule(n, v); err != nil {
			return nil, err
		}
	}

	err := userTypes.Each(func(k string, v jschemaLib.Schema) error {
		return s.Schema.AddType(k, v)
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *ExchangeJSightSchema) buildExchangeContent() error {
	var err error
	s.Exchange.once.Do(func() {
		n, err := s.Schema.GetAST()
		if err != nil {
			return
		}

		example, err := s.Schema.Example()
		if err != nil {
			return
		}

		s.Exchange.ContentJSight = astNodeToJsightContent(n, s.Exchange.UsedUserTypes, s.Exchange.UsedUserEnums)
		s.Exchange.Example = string(example)
	})
	return err
}

func (s *ExchangeJSightSchema) MarshalJSON() ([]byte, error) {
	data := struct {
		Content       interface{}             `json:"content,omitempty"`
		Example       string                  `json:"example,omitempty"`
		Notation      notation.SchemaNotation `json:"notation"`
		UsedUserTypes []string                `json:"usedUserTypes,omitempty"`
		UsedUserEnums []string                `json:"usedUserEnums,omitempty"`
	}{
		Notation: notation.SchemaNotationJSight,
	}

	err := s.buildExchangeContent()
	if err != nil {
		return nil, err
	}

	data.Content = s.Exchange.ContentJSight
	if s.Exchange.UsedUserTypes != nil && s.Exchange.UsedUserTypes.Len() > 0 {
		data.UsedUserTypes = s.Exchange.UsedUserTypes.Data()
	}
	if s.Exchange.UsedUserEnums != nil && s.Exchange.UsedUserEnums.Len() > 0 {
		data.UsedUserEnums = s.Exchange.UsedUserEnums.Data()
	}

	data.Example = s.Exchange.Example

	return json.Marshal(data)
}
