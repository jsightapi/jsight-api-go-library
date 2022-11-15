package catalog

import (
	schema "github.com/jsightapi/jsight-schema-go-library"
)

// UserSchemas represent available user type's schemas.
// gen:UnsafeOrderedMap
type UserSchemas struct {
	data  map[string]schema.Schema
	order []string
}
