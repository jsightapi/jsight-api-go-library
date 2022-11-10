package catalog

import (
	"errors"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/kit"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

type HTTPResponseBody struct {
	Format    SerializeFormat     `json:"format"`
	Schema    ExchangeSchema      `json:"schema"`
	Directive directive.Directive `json:"-"`
}

func NewHTTPResponseBody(
	b bytes.Bytes,
	f SerializeFormat,
	sn notation.SchemaNotation,
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]jschemaLib.Rule,
	catalogUserTypes *UserTypes,
) (HTTPResponseBody, *jerr.JApiError) {
	body := HTTPResponseBody{
		Format:    f,
		Schema:    nil,
		Directive: d,
	}

	var s ExchangeSchema
	var err error

	switch f {
	case SerializeFormatJSON:
		s, err = NewExchangeJSightSchema("", b.Data(), tt, rr, catalogUserTypes)
		if err != nil {
			return HTTPResponseBody{}, adoptErrorForResponseBody(d, err)
		}
	case SerializeFormatPlainString:
		s, err = PrepareRegexSchema("", b)
		if err != nil {
			return HTTPResponseBody{}, adoptErrorForResponseBody(d, err)
		}
	default:
		s = NewExchangePseudoSchema(sn)
	}

	body.Schema = s

	return body, nil
}

func adoptErrorForResponseBody(d directive.Directive, err error) *jerr.JApiError {
	var e kit.Error
	if errors.As(err, &e) {
		if d.BodyCoords.IsSet() {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}
		return d.ParameterError(e.Message())
	}
	return d.KeywordError(err.Error())
}
