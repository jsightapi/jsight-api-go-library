package core

import (
	innerSchema "github.com/jsightapi/jsight-schema-go-library/notations/jschema/schema"
)

type PieceOfPathVariable struct {
	node  innerSchema.Node
	types map[string]innerSchema.Type
}
