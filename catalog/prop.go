package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
)

type Prop struct {
	Parameter string
	ASTNode   jschemaLib.ASTNode
	Directive directive.Directive
}
