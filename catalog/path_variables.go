package catalog

import (
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type PathVariables struct {
	*ExchangeJSightSchema
}

func NewPathVariables(properties []Prop) *PathVariables {
	s := jschema.New("", "")
	_ = s.Build()

	s.ASTNode = jschemaLib.ASTNode{
		TokenType:  jschemaLib.TokenTypeObject,
		SchemaType: "object",
		Rules:      jschemaLib.MakeRuleASTNodes(0),
		Children:   make([]jschemaLib.ASTNode, 0, len(properties)),
	}

	for _, p := range properties {
		s.ASTNode.Children = append(s.ASTNode.Children, p.ASTNode)
	}

	return &PathVariables{
		ExchangeJSightSchema: &ExchangeJSightSchema{
			Schema: s,
			// TODO Exchange: NewExchange(),
		},
	}
}

func (p *PathVariables) Validate(key, value []byte) error {
	return p.Schema.ValidateObjectProperty(key, value)
}

func (p *PathVariables) MarshalJSON() ([]byte, error) {
	// TODO
	return []byte(`{}`), nil
}
