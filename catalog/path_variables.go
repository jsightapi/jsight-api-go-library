package catalog

import (
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type PathVariables struct {
	Schema *ExchangeJSightSchema `json:"schema"`
}

func NewPathVariables(properties []Prop, catalogUserTypes *UserTypes) *PathVariables {
	s := jschema.New("", "")
	_ = s.Build() // TODO It's necessary?

	s.ASTNode = jschemaLib.ASTNode{
		TokenType:  jschemaLib.TokenTypeObject,
		SchemaType: "object",
		Rules:      jschemaLib.MakeRuleASTNodes(0),
		Children:   make([]jschemaLib.ASTNode, 0, len(properties)),
	}

	for _, p := range properties {
		s.ASTNode.Children = append(s.ASTNode.Children, p.ASTNode)
	}

	es := newExchangeJSightSchema(s)
	es.DisableExchangeExample = true
	es.catalogUserTypes = catalogUserTypes

	return &PathVariables{
		Schema: es,
	}
}

func (p *PathVariables) Validate(key, value []byte) error {
	return p.Schema.ValidateObjectProperty(key, value)
}
