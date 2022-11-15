package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	schema "github.com/jsightapi/jsight-schema-go-library"
)

type Prop struct {
	Parameter string
	ASTNode   schema.ASTNode
	Directive directive.Directive
}
