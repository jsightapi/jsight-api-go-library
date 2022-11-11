package catalog

import (
	"encoding/json"
	"sync"

	"github.com/jsightapi/jsight-api-go-library/notation"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type ExchangeJSightSchema struct { //nolint:govet
	*jschema.Schema

	onceCompile            sync.Once
	catalogUserTypes       *UserTypes
	disableExchangeExample bool

	// exchangeContent a JSight schema.
	exchangeContent *ExchangeContent
	// UsedUserTypes a list of used user types.
	exchangeUsedUserTypes *StringSet
	// UserUserTypes a list of used user enums.
	exchangeUsedUserEnums *StringSet
}

func newExchangeJSightSchema(s *jschema.Schema) *ExchangeJSightSchema {
	return &ExchangeJSightSchema{
		Schema:                s,
		exchangeContent:       nil,
		exchangeUsedUserTypes: &StringSet{},
		exchangeUsedUserEnums: &StringSet{},
	}
}

func NewExchangeJSightSchema[T bytes.Byter](
	name string,
	b T,
	coreUserTypes *UserSchemas,
	coreRules map[string]jschemaLib.Rule,
	catalogUserTypes *UserTypes,
) (*ExchangeJSightSchema, error) {
	if s, ok := coreUserTypes.Get(name); ok {
		if ss, ok := s.(*jschema.Schema); ok {
			es := newExchangeJSightSchema(ss)
			es.catalogUserTypes = catalogUserTypes
			return es, nil
		}
	}

	es := newExchangeJSightSchema(jschema.New(name, b))
	es.catalogUserTypes = catalogUserTypes

	for n, v := range coreRules {
		if err := es.Schema.AddRule(n, v); err != nil {
			return nil, err
		}
	}

	err := coreUserTypes.Each(func(k string, v jschemaLib.Schema) error {
		return es.Schema.AddType(k, v)
	})
	if err != nil {
		return nil, err
	}

	err = es.Build()
	if err != nil {
		return nil, err
	}

	return es, nil
}

func (e *ExchangeJSightSchema) Compile() (err error) {
	e.onceCompile.Do(func() {
		err = e.buildContent()
		if err != nil {
			return
		}

		err = e.processAllOf(e.exchangeUsedUserTypes)
		if err != nil {
			return
		}
	})
	return err
}

func (e *ExchangeJSightSchema) buildContent() error {
	n, err := e.Schema.GetAST()
	if err != nil {
		return err
	}
	e.exchangeContent = astNodeToJsightContent(n, e.exchangeUsedUserTypes, e.exchangeUsedUserEnums)
	return nil
}

func (e *ExchangeJSightSchema) processAllOf(uut *StringSet) error {
	return e.exchangeContent.processAllOf(uut, e.catalogUserTypes)
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
		Content:  e.exchangeContent,
		Notation: notation.SchemaNotationJSight,
	}

	if !e.disableExchangeExample {
		example, err := e.Example()
		if err != nil {
			return nil, err
		}
		data.Example = string(example)
	}

	if e.exchangeUsedUserTypes != nil && e.exchangeUsedUserTypes.Len() > 0 {
		data.UsedUserTypes = e.exchangeUsedUserTypes.Data()
	}
	if e.exchangeUsedUserEnums != nil && e.exchangeUsedUserEnums.Len() > 0 {
		data.UsedUserEnums = e.exchangeUsedUserEnums.Data()
	}

	return json.Marshal(data)
}
