package catalog

import (
	"encoding/json"
	"github.com/jsightapi/jsight-api-go-library/notation"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"sync"
)

type ExchangeJSightSchema struct {
	*jschema.Schema

	onceCompile            sync.Once
	catalogUserTypes       *UserTypes
	DisableExchangeExample bool

	// TODO make the following properties private

	// ExchangeContent a JSight schema.
	ExchangeContent *ExchangeContent
	// UsedUserTypes a list of used user types.
	ExchangeUsedUserTypes *StringSet
	// UserUserTypes a list of used user enums.
	ExchangeUsedUserEnums *StringSet
	// ExchangeExample of schema.
	// ExchangeExample []byte
}

func newExchangeJSightSchema(s *jschema.Schema) *ExchangeJSightSchema {
	return &ExchangeJSightSchema{
		Schema:                s,
		ExchangeContent:       nil,
		ExchangeUsedUserTypes: &StringSet{},
		ExchangeUsedUserEnums: &StringSet{},
	}
}

func NewExchangeJSightSchema(
	name string,
	b []byte,
	userTypes *UserSchemas,
	enumRules map[string]jschemaLib.Rule,
	catalogUserTypes *UserTypes,
) (*ExchangeJSightSchema, error) {
	s := newExchangeJSightSchema(jschema.New(name, b))
	s.catalogUserTypes = catalogUserTypes

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

	err = s.Build() // TODO is there a place for this?
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (e *ExchangeJSightSchema) Compile() (err error) {
	e.onceCompile.Do(func() {
		err = e.buildContent()
		if err != nil {
			return
		}

		err = e.processAllOf(e.ExchangeUsedUserTypes)
		if err != nil {
			return
		}
	})
	return
}

func (e *ExchangeJSightSchema) buildContent() error {
	n, err := e.Schema.GetAST()
	if err != nil {
		return err
	}
	e.ExchangeContent = astNodeToJsightContent(n, e.ExchangeUsedUserTypes, e.ExchangeUsedUserEnums)
	return nil
}

func (e *ExchangeJSightSchema) processAllOf(uut *StringSet) error {
	return e.ExchangeContent.processAllOf(uut, e.catalogUserTypes)
}

func (e *ExchangeJSightSchema) Example() ([]byte, error) {
	// TODO once
	return e.Schema.Example()
}

func (e *ExchangeJSightSchema) MarshalJSON() ([]byte, error) {
	err := e.Compile()
	if err != nil {
		return nil, err
	}

	data := struct {
		Content       interface{}             `json:"content,omitempty"`
		Example       string                  `json:"example,omitempty"`
		Notation      notation.SchemaNotation `json:"notation"`
		UsedUserTypes []string                `json:"usedUserTypes,omitempty"`
		UsedUserEnums []string                `json:"usedUserEnums,omitempty"`
	}{
		Content:  e.ExchangeContent,
		Notation: notation.SchemaNotationJSight,
	}

	if !e.DisableExchangeExample {
		example, err := e.Example()
		if err != nil {
			return nil, err
		}
		data.Example = string(example)
	}

	if e.ExchangeUsedUserTypes != nil && e.ExchangeUsedUserTypes.Len() > 0 {
		data.UsedUserTypes = e.ExchangeUsedUserTypes.Data()
	}
	if e.ExchangeUsedUserEnums != nil && e.ExchangeUsedUserEnums.Len() > 0 {
		data.UsedUserEnums = e.ExchangeUsedUserEnums.Data()
	}

	return json.Marshal(data)
}
