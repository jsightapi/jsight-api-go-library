package catalog

import (
	"encoding/json"
	"github.com/jsightapi/jsight-api-go-library/notation"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
)

type PathVariables struct {
	ASTNodes *jschemaLib.ASTNodes
}

func (p *PathVariables) MarshalJSON() ([]byte, error) {
	var data struct {
		Schema Schema `json:"schema"`
	}

	n := jschemaLib.ASTNode{
		TokenType:  jschemaLib.TokenTypeObject,
		SchemaType: "object",
		Rules:      jschemaLib.MakeRuleASTNodes(0),
		Children:   make([]jschemaLib.ASTNode, 0, p.ASTNodes.Len()),
	}

	err := p.ASTNodes.Each(func(k string, v jschemaLib.ASTNode) error {
		n.Children = append(n.Children, v)
		return nil
	})
	if err != nil {
		return []byte{}, err
	}

	data.Schema = NewSchema(notation.SchemaNotationJSight)
	data.Schema.ContentJSight = astNodeToJsightContent(n, data.Schema.UsedUserTypes, data.Schema.UsedUserEnums)

	return json.Marshal(data)
}
