package core

import (
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
)

func (core *JApiCore) UserTypesData() map[string]*jschema.Schema {
	ut := make(map[string]*jschema.Schema, core.userTypes.Len())
	_ = core.userTypes.Each(func(k string, s jschemaLib.Schema) error { //nolint:errcheck // It's ok.
		if ss, ok := s.(*jschema.Schema); ok {
			ut[k] = ss
		}
		return nil
	})
	return ut
}
