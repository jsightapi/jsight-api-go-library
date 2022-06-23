package core

import (
	"errors"
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/kit"
)

func (core *JApiCore) addDirectives() *jerr.JApiError {
	for _, d := range core.directivesWithPastes {
		if je := core.addDirectiveBranch(d); je != nil {
			return je
		}
	}
	return nil
}

func (core *JApiCore) addDirectiveBranch(d *directive.Directive) *jerr.JApiError {
	if je := core.addDirective(d); je != nil {
		return je
	}

	for _, c := range d.Children {
		if je := core.addDirectiveBranch(c); je != nil {
			return je
		}
	}

	return nil
}

func (core *JApiCore) addDirective(d *directive.Directive) *jerr.JApiError {
	if _, ok := core.bannedDirectives[d.Type()]; ok {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.DirectiveNotAllowed, d.Type().String()))
	}

	f, ok := core.directiveFunctions[d.Type()]
	if !ok { // Path
		return nil
	}

	return f(d)
}

func (core JApiCore) addJSight(d *directive.Directive) *jerr.JApiError {
	version := d.Parameter("Version")
	if version == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Version"))
	}

	if version != lastJSightVersion {
		return d.KeywordError("unsupported version of JSIGHT")
	}
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if err := core.catalog.AddJSight(version); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addInfo(d *directive.Directive) *jerr.JApiError {
	if d.HasAnyParameters() {
		return d.KeywordError(jerr.ParametersAreForbiddenForTheDirective)
	}
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if err := core.catalog.AddInfo(*d); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addTitle(d *directive.Directive) *jerr.JApiError {
	title := d.Parameter("Title")
	if title == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Title"))
	}
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if err := core.catalog.AddTitle(title); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addVersion(d *directive.Directive) *jerr.JApiError {
	version := d.Parameter("Version")
	if version == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Version"))
	}
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if err := core.catalog.AddVersion(version); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addDescription(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if !d.BodyCoords.IsSet() {
		return d.KeywordError(jerr.EmptyDescription)
	}

	text := string(description(d.BodyCoords.Read()))
	if text == "" {
		return d.KeywordError(jerr.EmptyDescription)
	}

	switch d.Parent.Type() {
	case directive.Info:
		return core.addInfoDescription(d, text)
	case directive.Get, directive.Post, directive.Put, directive.Patch, directive.Delete:
		return core.addHTTPMethodDescription(d, text)
	case directive.Method:
		return core.addMethodDescription(d, text)
	default:
		return d.KeywordError("wrong description context")
	}
}

func (core JApiCore) addInfoDescription(d *directive.Directive, text string) *jerr.JApiError {
	if err := core.catalog.AddDescriptionToInfo(text); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addHTTPMethodDescription(d *directive.Directive, text string) *jerr.JApiError {
	if err := core.catalog.AddDescriptionToHttpMethod(*d, text); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (_ JApiCore) addMethodDescription(_ *directive.Directive, _ string) *jerr.JApiError {
	// TODO
	return nil
}

func (core JApiCore) addServer(d *directive.Directive) *jerr.JApiError {
	name := d.Parameter("Name")
	if name == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))
	}
	if err := core.catalog.AddServer(name, d.Annotation); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addBaseUrl(d *directive.Directive) *jerr.JApiError {
	path := d.Parameter("Path")
	if path == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Path"))
	}
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	server := d.Parent
	if err := core.catalog.AddBaseUrl(server.Parameter("Name"), path); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addType(d *directive.Directive) *jerr.JApiError {
	if d.Parameter("Name") == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))
	}
	return core.catalog.AddType(*d, core.userTypes, core.rules)
}

func (core *JApiCore) addURL(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	path, err := d.Path()
	if err != nil {
		return d.KeywordError(err.Error())
	}

	pp, err := PathParameters(path)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	err = core.checkSimilarPaths(pp)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	p := catalog.Path(path)

	if _, ok := core.uniqURLPath[p]; ok {
		return d.KeywordError(fmt.Sprintf("non-unique path %q in the URL directive", p))
	}

	core.uniqURLPath[p] = struct{}{}

	return nil
}

func (core JApiCore) addHTTPMethod(d *directive.Directive) *jerr.JApiError {
	path, err := d.Path()
	if err != nil {
		return d.KeywordError(err.Error())
	}

	pp, err := PathParameters(path)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	err = core.checkSimilarPaths(pp)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if err = core.catalog.AddMethod(*d); err != nil {
		return d.KeywordError(err.Error())
	}

	return nil
}

func (core JApiCore) addQuery(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if !d.BodyCoords.IsSet() {
		return d.KeywordError(jerr.EmptyBody)
	}

	q := catalog.NewQuery(*d)

	q.Format = d.Parameter("Format")
	if q.Format == "" {
		q.Format = "htmlFormEncoded"
	}

	example := d.Parameter("QueryExample")
	if example != "" {
		q.Example = example
	}

	s, err := catalog.UnmarshalSchema("", d.BodyCoords.Read(), core.userTypes, core.rules)
	if err != nil {
		var e kit.Error
		if errors.As(err, &e) {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}
		return d.BodyError(err.Error())
	}

	q.Schema = &s

	if err = core.catalog.AddQueryToCurrentMethod(*d, q); err != nil {
		return d.KeywordError(err.Error())
	}

	return nil
}

func (core JApiCore) addRequest(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	schemaNotation := d.Parameter("SchemaNotation")
	typ := d.Parameter("Type")

	if schemaNotation != "" && typ != "" {
		return d.KeywordError(jerr.CannotUseTheTypeAndSchemaNotationParametersTogether)
	}

	sn, err := notation.NewSchemaNotation(schemaNotation)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	bodyFormat, err := catalog.SchemaSerializeFormat(sn)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if d.Type() == directive.Request {
		if err = core.catalog.AddRequest(*d); err != nil {
			return d.KeywordError(err.Error())
		}
	}

	var s catalog.Schema

	switch {
	case sn == notation.SchemaNotationJSight && typ != "" && !d.BodyCoords.IsSet():
		if s, err = catalog.UnmarshalSchema(typ, bytes.Bytes(typ), core.userTypes, core.rules); err == nil {
			err = core.catalog.AddRequestBody(s, bodyFormat, *d)
		}

	case sn == notation.SchemaNotationJSight && typ == "" && d.BodyCoords.IsSet():
		if s, err = catalog.UnmarshalSchema("", d.BodyCoords.Read(), core.userTypes, core.rules); err == nil {
			err = core.catalog.AddRequestBody(s, bodyFormat, *d)
		}
		var e kit.Error
		if errors.As(err, &e) {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}

	case sn == notation.SchemaNotationRegex && typ == "" && d.BodyCoords.IsSet():
		s = catalog.NewRegexSchema(d.BodyCoords.Read())
		err = core.catalog.AddRequestBody(s, bodyFormat, *d)

	case (sn == notation.SchemaNotationAny || sn == notation.SchemaNotationEmpty) && !d.BodyCoords.IsSet():
		s = catalog.NewSchema(sn)
		err = core.catalog.AddRequestBody(s, bodyFormat, *d)

	case d.Type() == directive.Body:
		err = errors.New("incorrect request")
	}

	if err != nil {
		return d.KeywordError(err.Error())
	}

	return nil
}

func (core JApiCore) addResponse(d *directive.Directive) *jerr.JApiError {
	schemaNotationParam := d.Parameter("SchemaNotation")
	typeParam := d.Parameter("Type")

	if schemaNotationParam != "" && typeParam != "" {
		return d.KeywordError(jerr.CannotUseTheTypeAndSchemaNotationParametersTogether)
	}

	schemaNotation, err := notation.NewSchemaNotation(schemaNotationParam)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	bodyFormat, err := catalog.SchemaSerializeFormat(schemaNotation)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if d.Type() == directive.Body {
		d1 := d.Parent
		if d1.Type() == directive.HTTPResponseCode && typeParam != "" && d1.Parameter("Type") != "" {
			return d.KeywordError("You cannot specify User Type in the response directive if it has a child Body directive.")
		}
	}

	if d.Type() == directive.HTTPResponseCode {
		if err = core.catalog.AddResponse(d.Keyword, d.Annotation, *d); err != nil {
			return d.KeywordError(err.Error())
		}
	}

	var je *jerr.JApiError

	switch {
	case typeParam != "":
		je = core.catalog.AddResponseBody(typeParam, bytes.Bytes(typeParam), bodyFormat, schemaNotation, *d, core.userTypes, core.rules)

	case d.BodyCoords.IsSet():
		je = core.catalog.AddResponseBody("", d.BodyCoords.Read(), bodyFormat, schemaNotation, *d, core.userTypes, core.rules)

	case schemaNotation.IsAnyOrEmpty():
		je = core.catalog.AddResponseBody("", bytes.Bytes{}, bodyFormat, schemaNotation, *d, core.userTypes, core.rules)

	case d.Type() == directive.Body:
		je = d.KeywordError("body is empty")
	}

	return je
}

func (core JApiCore) addHeaders(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if !d.BodyCoords.IsSet() {
		return d.KeywordError(jerr.EmptyBody)
	}

	var s catalog.Schema
	var err error

	s, err = catalog.UnmarshalSchema("", d.BodyCoords.Read(), core.userTypes, core.rules)
	if err != nil {
		var e kit.Error
		if errors.As(err, &e) {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}
		return d.BodyError(err.Error())
	}

	switch d.Parent.Type() {
	case directive.Request:
		err = core.catalog.AddRequestHeaders(s, *d)
	case directive.HTTPResponseCode:
		err = core.catalog.AddResponseHeaders(s, *d)
	default:
		err = errors.New("incorrect directive context")
	}

	if err != nil {
		return d.KeywordError(err.Error())
	}

	return nil
}

func (core JApiCore) addBody(d *directive.Directive) *jerr.JApiError {
	if d.Parent.HasAnyParameters() && d.Parent.Type() != directive.Macro {
		return d.Parent.KeywordError("parameters are unacceptable, according to the Body directive")
	}

	switch d.Parent.Type() {
	case directive.Request:
		return core.addRequest(d)
	case directive.HTTPResponseCode:
		return core.addResponse(d)
	default:
		return nil
	}
}

func (core JApiCore) addProtocol(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	if d.Parameter("Protocol") != "json-rpc-2.0" {
		return d.KeywordError("the parameter value should be \"json-rpc-2.0\"")
	}

	if err := core.catalog.AddProtocol(*d); err != nil {
		return d.KeywordError(err.Error())
	}

	return nil
}
