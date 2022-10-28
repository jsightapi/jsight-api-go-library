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

func (b PathVariablesBuilder) AddProperty(key string, value innerSchema.Node) {
	b.objectBuilder.AddProperty(key, value)
}

func (b PathVariablesBuilder) Len() int {
	return b.objectBuilder.Len()
}

func (b PathVariablesBuilder) Build() *PathVariables {
	s := b.objectBuilder.Build()

	es := newExchangeJSightSchema(s)
	es.disableExchangeExample = true
	es.catalogUserTypes = b.catalogUserTypes

	return &PathVariables{
		Schema: es,
	}
}
