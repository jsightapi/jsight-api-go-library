package catalog

import (
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	innerSchema "github.com/jsightapi/jsight-schema-go-library/notations/jschema/schema"
)

type PathVariablesBuilder struct {
	catalogUserTypes *UserTypes
	objectBuilder    jschema.ObjectBuilder
}

func NewPathVariablesBuilder(catalogUserTypes *UserTypes) PathVariablesBuilder {
	return PathVariablesBuilder{
		catalogUserTypes: catalogUserTypes,
		objectBuilder:    jschema.NewObjectBuilder(),
	}
}

func (b PathVariablesBuilder) AddProperty(key string, node innerSchema.Node, types map[string]innerSchema.Type) {
	b.objectBuilder.AddProperty(key, node, types)
}

func (b PathVariablesBuilder) Len() int {
	return b.objectBuilder.Len()
}

func (b PathVariablesBuilder) Build() *PathVariables {
	uutNames := b.objectBuilder.UserTypeNames()
	for _, name := range uutNames {
		if ut, ok := b.catalogUserTypes.Get(name); ok {
			if es, ok := ut.Schema.(*ExchangeJSightSchema); ok {
				b.objectBuilder.AddType(name, es.Schema)
			}
		}
	}

	s := b.objectBuilder.Build()

	es := newExchangeJSightSchema(s)
	es.disableExchangeExample = true
	es.catalogUserTypes = b.catalogUserTypes

	return &PathVariables{
		Schema: es,
	}
}
