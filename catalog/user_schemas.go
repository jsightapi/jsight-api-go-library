package catalog

import (
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
)

// UserSchemas represent available user type's schemas.
// gen:UnsafeOrderedMap
type UserSchemas struct {
	data  map[string]jschemaLib.Schema
	order []string
}
