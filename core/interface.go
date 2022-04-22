package core

import (
	"j/japi/catalog"
	"j/japi/jerr"
)

// ValidateJAPI should be used to check if .jst file is valid according to specification
func (core *JApiCore) ValidateJAPI() *jerr.JAPIError {
	return core.processJapiProject()
}

func (core *JApiCore) Catalog() *catalog.Catalog {
	return core.catalog
}
