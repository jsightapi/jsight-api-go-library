package catalog

import (
	schema "github.com/jsightapi/jsight-schema-go-library"
)

type ExchangeSchema interface {
	schema.Schema
	MarshalJSON() ([]byte, error)
}
