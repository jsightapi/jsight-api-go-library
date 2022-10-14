package catalog

import (
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
)

type ExchangeSchema interface {
	jschemaLib.Schema
	MarshalJSON() ([]byte, error)
}
