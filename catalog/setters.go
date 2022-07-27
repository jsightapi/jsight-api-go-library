package catalog

import (
	"errors"
	"fmt"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/kit"
	"github.com/jsightapi/jsight-schema-go-library/rules/enum"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

// pathTag returns a Tag from the collection, or creates a new one and adds it to the collection
func (c *Catalog) pathTag(r InteractionId) *Tag {
	t := newPathTag(r)

	if tt, ok := c.Tags.Get(t.Name); ok {
		return tt
	}

	c.Tags.Set(t.Name, t)

	return t
}

func (c *Catalog) AddTag(name, title string) error {
	if c.Tags.Has(TagName(name)) {
		return fmt.Errorf("%s (%q)", jerr.DuplicateNames, name)
	}

	t := NewTag(name, title)

	c.Tags.Set(t.Name, t)

	return nil
}

func (c *Catalog) AddTags(d directive.Directive) error {
	// TODO URL & RPC
	id, err := newHttpInteractionId(*d.Parent)
	if err != nil {
		return err
	}

	for _, name := range d.UnnamedParameter() {
		tn := TagName(name)
		t, ok := c.Tags.Get(tn)
		if !ok {
			return fmt.Errorf("%s %q", jerr.TagNotFound, tn)
		}
		t.appendInteractionId(id)

		c.Interactions.Update(id, func(v Interaction) Interaction {
			v.appendTagName(tn)
			return v
		})
	}

	return nil
}

func (c *Catalog) AddDescriptionToTag(name, description string) error {
	n := TagName(name)

	t, ok := c.Tags.Get(n)
	if !ok {
		return fmt.Errorf("%s (%q)", jerr.TagNotFound, n)
	}

	if t.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Tags.Update(n, func(v *Tag) *Tag {
		v.Description = &description
		return v
	})

	return nil
}

func (c *Catalog) AddTagDescription(name, title string) error {
	t := NewTag(name, title)

	if _, ok := c.Tags.Get(t.Name); ok {
		return errors.New("")
	}

	c.Tags.Set(t.Name, t)

	return nil
}

func (c *Catalog) AddJSight(version string) error {
	if c.JSightVersion != "" {
		return errors.New("directive JSIGHT gotta be only one time")
	}

	c.JSightVersion = version

	return nil
}

func (c *Catalog) AddInfo(d directive.Directive) error {
	if c.Info != nil {
		return errors.New("directive INFO gotta be only one time")
	}
	c.Info = &Info{Directive: d}
	return nil
}

func (c *Catalog) AddTitle(name string) error {
	if c.Info.Title != "" {
		return errors.New(jerr.NotUniqueDirective)
	}
	c.Info.Title = name
	return nil
}

func (c *Catalog) AddVersion(version string) error {
	if c.Info.Version != "" {
		return errors.New(jerr.NotUniqueDirective)
	}
	c.Info.Version = version
	return nil
}

func (c *Catalog) AddDescriptionToInfo(text string) error {
	if c.Info.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Info.Description = &text

	return nil
}

func (c *Catalog) AddHTTPMethod(d directive.Directive) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if c.Interactions.Has(httpId) {
		return fmt.Errorf("method is already defined in resource %s", httpId.String())
	}

	t := c.pathTag(httpId)
	t.appendInteractionId(httpId)

	c.Interactions.Set(httpId, newHttpInteraction(httpId, d.Annotation, t.Name))

	return nil
}

func (c *Catalog) AddDescriptionToHttpMethod(d directive.Directive, text string) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpId) {
		return fmt.Errorf("%s %q", jerr.HttpResourceNotFound, httpId.String())
	}

	v := c.Interactions.GetValue(httpId).(*HttpInteraction) //nolint:errcheck

	if v.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		v.(*HttpInteraction).Description = &text
		return v
	})

	return nil
}

func (c *Catalog) AddDescriptionToJsonRpcMethod(d directive.Directive, text string) error {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(rpcId) {
		return fmt.Errorf("%s %q", jerr.HttpResourceNotFound, rpcId.String())
	}

	v := c.Interactions.GetValue(rpcId).(*JsonRpcInteraction) //nolint:errcheck

	if v.Description != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(rpcId, func(v Interaction) Interaction {
		v.(*JsonRpcInteraction).Description = &text
		return v
	})

	return nil
}

func (c *Catalog) AddQueryToCurrentMethod(d directive.Directive, q Query) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpId) {
		return fmt.Errorf("%s %q", jerr.HttpResourceNotFound, httpId.String())
	}

	v := c.Interactions.GetValue(httpId).(*HttpInteraction) //nolint:errcheck

	if v.Query != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		v.(*HttpInteraction).Query = &q
		return v
	})

	return nil
}

func (c *Catalog) AddResponse(code string, annotation string, d directive.Directive) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	r := HTTPResponse{Code: code, Annotation: annotation, Directive: d}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		v.(*HttpInteraction).Responses = append(v.(*HttpInteraction).Responses, r)
		return v
	})

	return nil
}

func (c *Catalog) AddResponseBody(
	schemaBytes bytes.Bytes,
	bodyFormat SerializeFormat,
	sn notation.SchemaNotation,
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]jschemaLib.Rule,
) *jerr.JApiError {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	if !c.Interactions.Has(httpId) {
		return d.KeywordError(fmt.Sprintf("%s %q", jerr.HttpResourceNotFound, httpId.String()))
	}

	v := c.Interactions.GetValue(httpId).(*HttpInteraction) //nolint:errcheck

	i := len(v.Responses) - 1
	if i == -1 {
		return d.KeywordError(fmt.Sprintf("%s for %q", jerr.ResponsesIsEmpty, httpId.String()))
	}

	httpResponseBody, je := NewHTTPResponseBody(schemaBytes, bodyFormat, sn, d, tt, rr)
	if je != nil {
		return je
	}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		v.(*HttpInteraction).Responses[i].Body = &httpResponseBody
		return v
	})

	return nil
}

func (c *Catalog) AddResponseHeaders(s Schema, d directive.Directive) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpId) {
		return fmt.Errorf("%s %q", jerr.HttpResourceNotFound, httpId.String())
	}

	v := c.Interactions.GetValue(httpId).(*HttpInteraction) //nolint:errcheck

	i := len(v.Responses) - 1
	if i == -1 {
		return fmt.Errorf("%s for %q", jerr.ResponsesIsEmpty, httpId.String())
	}

	if v.Responses[i].Headers != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		v.(*HttpInteraction).Responses[i].Headers = &HTTPResponseHeaders{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddServer(name string, annotation string) error {
	if c.Servers.Has(name) {
		return fmt.Errorf("%s (%q)", jerr.DuplicateNames, name)
	}

	server := new(Server)
	server.Annotation = annotation

	c.Servers.Set(name, server)

	return nil
}

func (c *Catalog) AddBaseUrl(serverName string, path string) error {
	if !c.Servers.Has(serverName) {
		return fmt.Errorf("server not found for %s", serverName)
	}

	v := c.Servers.GetValue(serverName)

	if v.BaseUrl != "" {
		return errors.New("BaseUrl already defined")
	}

	v.BaseUrl = path

	// if d.BodyCoords.IsSet() {
	// 	// baseurl has jschema body
	// 	s, err := UnmarshalSchema("", d.BodyCoords.Read())
	// 	if err != nil {
	// 		if e, ok := err.(kit.Error); ok {
	// 			return c.japiError(e.Message(), d.BodyCoords.B()+bytes.Index(e.Position()))
	// 		}
	// 		return c.japiError(err.Error(), d.BodyBegin())
	// 	}
	//
	// 	if s.ContentJSight.TokenType != objectStr && s.ContentJSight.TokenType != shortcutStr {
	// 		return c.japiError("the body of the BaseUrl directive can contain an object schema", d.BodyBegin())
	// 	}
	//
	// 	c.Servers.Update(serverName, func(v *Server) *Server {
	// 		v.BaseUrlVariables = &baseUrlVariables{
	// 			Schema:    &s,
	// 			directive: d,
	// 		}
	// 		return v
	// 	})
	// }

	return nil
}

func (c *Catalog) AddType(
	d directive.Directive,
	tt *UserSchemas,
	rr map[string]jschemaLib.Rule,
) *jerr.JApiError {
	name := d.NamedParameter("Name")

	if c.UserTypes.Has(name) {
		return d.KeywordError(fmt.Sprintf("%s (%q)", jerr.DuplicateNames, name))
	}

	userType := &UserType{
		Annotation: d.Annotation,
		Directive:  d,
	}
	typeNotation, err := notation.NewSchemaNotation(d.NamedParameter("SchemaNotation"))
	if err != nil {
		return d.KeywordError(err.Error())
	}

	switch typeNotation {
	case notation.SchemaNotationJSight:
		if !d.BodyCoords.IsSet() {
			return d.KeywordError(jerr.EmptyBody)
		}
		b := d.BodyCoords.Read()
		schema, err := UnmarshalSchema(name, b, tt, rr)
		if err != nil {
			var e kit.Error
			if errors.As(err, &e) {
				return d.BodyErrorIndex(e.Message(), e.Position())
			}
			return d.KeywordError(err.Error())
		}
		userType.Schema = schema
	case notation.SchemaNotationRegex:
		if !d.BodyCoords.IsSet() {
			return d.KeywordError(jerr.EmptyBody)
		}
		userType.Schema = NewRegexSchema(d.BodyCoords.Read())
	case notation.SchemaNotationAny, notation.SchemaNotationEmpty:
		userType.Schema = NewSchema(typeNotation)
	}

	c.UserTypes.Set(name, userType)

	return nil
}

func (c *Catalog) AddRequest(d directive.Directive) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		if v.(*HttpInteraction).Request == nil {
			v.(*HttpInteraction).Request = &HTTPRequest{
				Directive: d,
			}
		}
		return v
	})

	return nil
}

func (c *Catalog) AddRequestBody(s Schema, f SerializeFormat, d directive.Directive) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpId) {
		return fmt.Errorf("%s %q", jerr.HttpResourceNotFound, httpId.String())
	}

	v := c.Interactions.GetValue(httpId).(*HttpInteraction) //nolint:errcheck

	if v.Request == nil {
		return fmt.Errorf("%s for %q", jerr.RequestIsEmpty, httpId.String())
	}

	if v.Request.HTTPRequestBody != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		v.(*HttpInteraction).Request.HTTPRequestBody = &HTTPRequestBody{Format: f, Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddRequestHeaders(s Schema, d directive.Directive) error {
	httpId, err := newHttpInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(httpId) {
		return fmt.Errorf("%s %q", jerr.HttpResourceNotFound, httpId.String())
	}

	v := c.Interactions.GetValue(httpId).(*HttpInteraction) //nolint:errcheck

	if v.Request == nil {
		return fmt.Errorf("%s for %q", jerr.RequestIsEmpty, httpId.String())
	}

	if v.Request.HTTPRequestHeaders != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(httpId, func(v Interaction) Interaction {
		v.(*HttpInteraction).Request.HTTPRequestHeaders = &HTTPRequestHeaders{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddJsonRpcMethod(d directive.Directive) error {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return err
	}

	if c.Interactions.Has(rpcId) {
		return fmt.Errorf("method is already defined in resource %s", rpcId.String())
	}

	t := c.pathTag(rpcId)
	t.appendInteractionId(rpcId)

	c.Interactions.Set(rpcId, newJsonRpcInteraction(rpcId, d.NamedParameter("MethodName"), d.Annotation, t.Name))

	return nil
}

func (c *Catalog) AddJsonRpcParams(s Schema, d directive.Directive) error {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(rpcId) {
		return fmt.Errorf("%s %q", jerr.JsonRpcResourceNotFound, rpcId.String())
	}

	v := c.Interactions.GetValue(rpcId).(*JsonRpcInteraction) //nolint:errcheck

	if v.Params != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(rpcId, func(v Interaction) Interaction {
		v.(*JsonRpcInteraction).Params = &jsonRpcParams{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddJsonRpcResult(s Schema, d directive.Directive) error {
	rpcId, err := newJsonRpcInteractionId(d)
	if err != nil {
		return err
	}

	if !c.Interactions.Has(rpcId) {
		return fmt.Errorf("%s %q", jerr.JsonRpcResourceNotFound, rpcId.String())
	}

	v := c.Interactions.GetValue(rpcId).(*JsonRpcInteraction) //nolint:errcheck

	if v.Result != nil {
		return errors.New(jerr.NotUniqueDirective)
	}

	c.Interactions.Update(rpcId, func(v Interaction) Interaction {
		v.(*JsonRpcInteraction).Result = &jsonRpcResult{Schema: &s, Directive: d}
		return v
	})

	return nil
}

func (c *Catalog) AddEnum(d *directive.Directive, e *enum.Enum) *jerr.JApiError {
	name := d.NamedParameter("Name")
	if c.UserEnums.Has(name) {
		return d.KeywordError(fmt.Sprintf("%s (%q)", jerr.DuplicateNames, name))
	}

	r, err := c.enumDirectiveToUserRule(d, e)
	if err != nil {
		return d.KeywordError(err.Error())
	}

	c.UserEnums.Set(name, r)
	return nil
}

func (*Catalog) enumDirectiveToUserRule(d *directive.Directive, e *enum.Enum) (*UserRule, error) {
	vv, err := e.Values()
	if err != nil {
		return nil, err
	}

	r := Rule{
		Key:       d.NamedParameter("Name"),
		TokenType: RuleTokenTypeObject,
	}

	for _, v := range vv {
		r.Children = append(r.Children, Rule{
			Key:         v.Value.Unquote().String(),
			TokenType:   RuleTokenType(v.Type.ToTokenType()),
			ScalarValue: v.Value.Unquote().String(),
		})
	}

	return &UserRule{
		Annotation: d.Annotation,
		Value:      r,
		Directive:  d,
	}, nil
}
