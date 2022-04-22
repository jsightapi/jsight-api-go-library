package catalog

import (
	"errors"
	"j/japi/directive"
	"j/japi/jerr"
	"j/japi/notation"
	"j/schema/bytes"
	"j/schema/kit"
)

type HTTPResponseBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

func NewHTTPResponseBody(
	name string,
	b bytes.Bytes,
	f SerializeFormat,
	sn notation.SchemaNotation,
	d directive.Directive,
	tt *UserSchemas,
) (HTTPResponseBody, *jerr.JAPIError) {
	body := HTTPResponseBody{
		Format:    f,
		Schema:    nil,
		Directive: d,
	}

	var s Schema
	switch f {
	case SerializeFormatJSON:
		var err error
		s, err = UnmarshalSchema(name, b, tt)
		if err != nil {
			var e kit.Error
			if errors.As(err, &e) {
				if d.BodyCoords.IsSet() {
					return body, d.BodyErrorIndex(e.Message(), e.Position())
				}
				return body, d.ParameterError(e.Message())
			}
			return body, d.KeywordError(err.Error())
		}
	case SerializeFormatPlainString:
		s = NewRegexSchema(b)
	default:
		s = NewSchema(sn)
	}

	body.Schema = &s

	return body, nil
}
