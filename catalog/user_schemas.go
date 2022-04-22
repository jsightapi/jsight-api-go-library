package catalog

import (
	"j/schema"
)

// UserSchemas represent available user type's schemas.
// gen:UnsafeOrderedMap
type UserSchemas struct {
	data  map[string]jschema.Schema
	order []string
}
