package catalog

import (
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

type PathVariables struct {
	ExchangeJSightSchema
}

func NewPathVariables(properties []Prop) *PathVariables {
	n := jschemaLib.ASTNode{
		TokenType:  jschemaLib.TokenTypeObject,
		SchemaType: "object",
		Rules:      jschemaLib.MakeRuleASTNodes(0),
		Children:   make([]jschemaLib.ASTNode, 0, len(properties)),
	}

	for _, p := range properties {
		n.Children = append(n.Children, p.ASTNode)
	}

	return &PathVariables{
		ExchangeJSightSchema: ExchangeJSightSchema{
			Schema: &jschema.Schema{
				AstNode: n,
			},
		},
	}
}
