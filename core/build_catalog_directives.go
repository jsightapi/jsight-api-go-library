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
	version := d.NamedParameter("Version")
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
	if d.HasNamedParameter() {
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
	title := d.NamedParameter("Title")
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
	version := d.NamedParameter("Version")
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

	bb, err := description(d.BodyCoords.Read())
	if err != nil {
		return d.BodyError(err.Error())
	}
	if len(bb) == 0 {
		return d.KeywordError(jerr.EmptyDescription)
	}

	text := string(bb)

	switch d.Parent.Type() {
	case directive.Info:
		return core.addInfoDescription(d, text)
	case directive.Get, directive.Post, directive.Put, directive.Patch, directive.Delete:
		return core.addHTTPMethodDescription(d, text)
	case directive.Method:
		return core.addJsonRpcMethodDescription(d, text)
	case directive.TAG:
		return core.addTagDescription(d, d.Parent.NamedParameter("TagName"), text)
	default:
		return d.KeywordError("wrong description context")
	}
}

func (core JApiCore) addInfoDescription(d *directive.Directive, description string) *jerr.JApiError {
	if err := core.catalog.AddDescriptionToInfo(description); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addHTTPMethodDescription(d *directive.Directive, description string) *jerr.JApiError {
	if err := core.catalog.AddDescriptionToHTTPMethod(*d, description); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addJsonRpcMethodDescription(d *directive.Directive, description string) *jerr.JApiError {
	if err := core.catalog.AddDescriptionToJsonRpcMethod(*d, description); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addTagDescription(d *directive.Directive, name, description string) *jerr.JApiError {
	if err := core.catalog.AddDescriptionToTag(name, description); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addServer(d *directive.Directive) *jerr.JApiError {
	name := d.NamedParameter("Name")
	if name == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))
	}
	if err := core.catalog.AddServer(name, d.Annotation); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addBaseUrl(d *directive.Directive) *jerr.JApiError {
	path := d.NamedParameter("Path")
	if path == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Path"))
	}
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}

	server := d.Parent
	if err := core.catalog.AddBaseURL(server.NamedParameter("Name"), path); err != nil {
		return d.KeywordError(err.Error())
	}
	return nil
}

func (core JApiCore) addType(d *directive.Directive) *jerr.JApiError {
	if d.NamedParameter("Name") == "" {
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

	return checkJsonRpcUrlChildCompatible(d)
}

// checkJsonRpcUrlChildCompatible checks the compatibility of child directives for HTTP and JSON-RPC protocols.
func checkJsonRpcUrlChildCompatible(d *directive.Directive) *jerr.JApiError {
	if len(d.Children) == 0 {
		return nil
	}

	var base *directive.Directive
	var isBaseJsonRpc bool

	for _, dd := range d.Children {
		if base == nil {
			base = dd
			isBaseJsonRpc = isJsonRpcUrlChildDirective(base)
		} else if isBaseJsonRpc != isJsonRpcUrlChildDirective(dd) {
			return dd.KeywordError(
				fmt.Sprintf("directives %q and %q cannot be within the same URL directive",
					base.Type().String(),
					dd.Type().String(),
				),
			)
		}
	}

	return nil
}

func isJsonRpcUrlChildDirective(d *directive.Directive) bool {
	switch d.Type() {
	case directive.Protocol, directive.Method:
		return true
	default:
		return false
	}
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

	return core.catalog.AddHTTPMethod(*d)
}

func (core JApiCore) addQuery(d *directive.Directive) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if !d.BodyCoords.IsSet() {
		return d.KeywordError(jerr.EmptyBody)
	}

	q := catalog.NewQuery(*d)

	q.Format = d.NamedParameter("Format")
	if q.Format == "" {
		q.Format = "htmlFormEncoded"
	}

	example := d.NamedParameter("QueryExample")
	if example != "" {
		q.Example = example
	}

	s, err := catalog.UnmarshalJSightSchema("", d.BodyCoords.Read(), core.userTypes, core.rules)
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

	schemaNotation := d.NamedParameter("SchemaNotation")
	typ := d.NamedParameter("Type")

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
		if s, err = catalog.UnmarshalJSightSchema("", bytes.Bytes(typ), core.userTypes, core.rules); err == nil {
			err = core.catalog.AddRequestBody(s, bodyFormat, *d)
		}

	case sn == notation.SchemaNotationJSight && typ == "" && d.BodyCoords.IsSet():
		if s, err = catalog.UnmarshalJSightSchema("", d.BodyCoords.Read(), core.userTypes, core.rules); err == nil {
			err = core.catalog.AddRequestBody(s, bodyFormat, *d)
		}
		var e kit.Error
		if errors.As(err, &e) {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}

	case sn == notation.SchemaNotationRegex && typ == "" && d.BodyCoords.IsSet():
		if s, err = catalog.UnmarshalRegexSchema("", d.BodyCoords.Read()); err == nil {
			err = core.catalog.AddRequestBody(s, bodyFormat, *d)
		}
		var e kit.Error
		if errors.As(err, &e) {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}

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
	schemaNotationParam := d.NamedParameter("SchemaNotation")
	typeParam := d.NamedParameter("Type")

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
		if d1.Type() == directive.HTTPResponseCode && typeParam != "" && d1.NamedParameter("Type") != "" {
			return d.KeywordError(
				"You cannot specify User Type in the response directive if it has a child Body directive.",
			)
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
		je = core.catalog.AddResponseBody(
			bytes.Bytes(typeParam),
			bodyFormat,
			schemaNotation,
			*d,
			core.userTypes,
			core.rules,
		)

	case d.BodyCoords.IsSet():
		je = core.catalog.AddResponseBody(
			d.BodyCoords.Read(),
			bodyFormat,
			schemaNotation,
			*d,
			core.userTypes,
			core.rules,
		)

	case schemaNotation.IsAnyOrEmpty():
		je = core.catalog.AddResponseBody(
			nil,
			bodyFormat,
			schemaNotation,
			*d,
			core.userTypes,
			core.rules,
		)

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

	s, err = catalog.UnmarshalJSightSchema("", d.BodyCoords.Read(), core.userTypes, core.rules)
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
	if d.Parent.HasNamedParameter() && d.Parent.Type() != directive.Macro {
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

	if d.NamedParameter("ProtocolName") == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "ProtocolName"))
	}

	if d.NamedParameter("ProtocolName") != "json-rpc-2.0" {
		return d.KeywordError("the parameter value have to be \"json-rpc-2.0\"")
	}

	if _, ok := core.onlyOneProtocolIntoURL[d.Parent]; ok {
		return d.KeywordError("the directive Protocol must be unique within the directive URL")
	}
	core.onlyOneProtocolIntoURL[d.Parent] = struct{}{}

	return nil
}

func (core JApiCore) addJsonRpcMethod(d *directive.Directive) *jerr.JApiError {
	if d.NamedParameter("MethodName") == "" {
		return d.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "MethodName"))
	}

	if !isProtocolExists(d) {
		return d.KeywordError("the directive \"Protocol\" was not found")
	}

	return core.catalog.AddJsonRpcMethod(*d)
}

func isProtocolExists(d *directive.Directive) bool {
	for _, dd := range d.Parent.Children {
		if dd.Type() == directive.Protocol {
			return true
		}
	}
	return false
}

func (core JApiCore) addJsonRpcSchema(
	d *directive.Directive,
	f func(catalog.Schema, directive.Directive) error,
) *jerr.JApiError {
	if d.Annotation != "" {
		return d.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if !d.BodyCoords.IsSet() {
		return d.KeywordError(jerr.EmptyBody)
	}

	var s catalog.Schema
	var err error

	s, err = catalog.UnmarshalJSightSchema("", d.BodyCoords.Read(), core.userTypes, core.rules)
	if err != nil {
		var e kit.Error
		if errors.As(err, &e) {
			return d.BodyErrorIndex(e.Message(), e.Position())
		}
		return d.BodyError(err.Error())
	}

	err = f(s, *d)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	return nil
}

func (core JApiCore) addJsonRpcParams(d *directive.Directive) *jerr.JApiError {
	return core.addJsonRpcSchema(d, core.catalog.AddJsonRpcParams)
}

func (core JApiCore) addJsonRpcResult(d *directive.Directive) *jerr.JApiError {
	return core.addJsonRpcSchema(d, core.catalog.AddJsonRpcResult)
}
