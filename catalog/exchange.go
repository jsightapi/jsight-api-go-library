package catalog

import (
	"sync"
)

type Exchange struct {
	once sync.Once
	// ContentJSight a JSight schema.
	ContentJSight *ExchangeSchemaContentJSight
	// UsedUserTypes a list of used user types.
	UsedUserTypes *StringSet
	// UserUserTypes a list of used user enums.
	UsedUserEnums *StringSet
	// Example of schema.
	Example string
}

func SrtPtr(s string) *string {
	return &s
}
