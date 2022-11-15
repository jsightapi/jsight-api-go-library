package core

import (
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/ischema"
)

type PieceOfPathVariable struct {
	node  ischema.Node
	types map[string]ischema.Type
}
