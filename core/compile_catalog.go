package core

import (
	"errors"
	"fmt"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"strings"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

func (core *JApiCore) compileCatalog() *jerr.JApiError {
	if je := core.ProcessAllOf(); je != nil {
		return je
	}

	if je := core.ExpandRawPathVariableShortcuts(); je != nil {
		return je
	}

	if je := core.CheckRawPathVariableSchemas(); je != nil {
		return je
	}

	return core.BuildResourceMethodsPathVariables()
}

func (core *JApiCore) ExpandRawPathVariableShortcuts() *jerr.JApiError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		r := &core.rawPathVariables[i]

		for r.schema.AstNode.TokenType == jschemaLib.TokenTypeShortcut {
			typeName := r.schema.AstNode.SchemaType
			if typeName == "mixed" {
				return r.pathDirective.KeywordError("The root schema object cannot have an OR rule")
			}

			ut, ok := core.catalog.UserTypes.Get(typeName)
			if !ok {
				return r.pathDirective.KeywordError(fmt.Sprintf(`User type "%s" not found`, typeName))
			}

			r.schema = ut.Schema.JSchema.(*jschema.Schema) // copy schema TODO ???
		}

		// TODO
		// if err := checkPathSchema(r.schema); err != nil {
		// 	return r.pathDirective.KeywordError(err.Error())
		// }
	}

	return nil
}

func (core *JApiCore) CheckRawPathVariableSchemas() *jerr.JApiError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		if err := checkPathSchema(core.rawPathVariables[i].schema); err != nil {
			return core.rawPathVariables[i].pathDirective.KeywordError(err.Error())
		}
	}
	return nil
}

func checkPathSchema(s *jschema.Schema) error {
	if s.AstNode.TokenType != jschemaLib.TokenTypeObject {
		return errors.New("the body of the Path DIRECTIVE must be an object")
	}

	if s.AstNode.Rules.Has("additionalProperties") {
		return errors.New(`the "additionalProperties" rule is invalid in the Path directive`)
	}

	if s.AstNode.Rules.Has("nullable") {
		return errors.New(`the "nullable" rule is invalid in the Path directive`)
	}

	if s.AstNode.Rules.Has("or") {
		return errors.New(`the "or" rule is invalid in the Path directive`)
	}

	if s.AstNode.Children == nil || len(s.AstNode.Children) == 0 {
		return errors.New("an empty object in the Path directive")
	}

	for _, v := range s.AstNode.Children {
		if v.TokenType == jschemaLib.TokenTypeObject || v.TokenType == jschemaLib.TokenTypeArray {
			return fmt.Errorf("the multi-level property %q is not allowed in the Path directive", v.Key)
		}
	}

	return nil
}

func (core *JApiCore) BuildResourceMethodsPathVariables() *jerr.JApiError {
	allProjectProperties := make(map[catalog.Path]prop)
	for _, v := range core.rawPathVariables {
		n, err := v.schema.GetAST()
		if err != nil {
			return v.pathDirective.KeywordError(jerr.InternalServerError)
		}

		pp := core.propertiesToMap(n)

		for _, p := range v.parameters {
			if nn, ok := pp[p.parameter]; ok {
				if _, ok := allProjectProperties[p.path]; ok {
					return v.pathDirective.KeywordError(fmt.Sprintf(
						"The parameter %q has already been defined earlier",
						p.parameter,
					))
				}

				allProjectProperties[p.path] = prop{
					astNode:   nn,
					directive: v.pathDirective,
				}

				delete(pp, p.parameter)
			}
		}

		// Check that all path properties in schema is exists in the path.
		if len(pp) > 0 {
			ss := core.getPropertiesNames(pp)
			return v.pathDirective.KeywordError(fmt.Sprintf("Has unused parameters %q in schema", ss))
		}
	}

	// set PathVariables
	err := core.catalog.Interactions.Map(
		func(_ catalog.InteractionID, v catalog.Interaction) (catalog.Interaction, error) {
			if hi, ok := v.(*catalog.HTTPInteraction); ok {
				pp := pathParameters(v.Path().String())
				properties := make([]prop, 0, len(pp))

				for _, p := range pp {
					if pr, ok := allProjectProperties[p.path]; ok {
						pr.parameter = p.parameter
						properties = append(properties, pr)
					}
				}

				if len(properties) != 0 {
					hi.SetPathVariables(core.newPathVariables(properties))
				}
			}
			return v, nil
		},
	)
	if err != nil {
		return err.(*jerr.JApiError) //nolint:errorlint
	}

	return nil
}

func (*JApiCore) propertiesToMap(n jschemaLib.ASTNode) map[string]jschemaLib.ASTNode {
	if len(n.Children) == 0 {
		return nil
	}

	res := make(map[string]jschemaLib.ASTNode, len(n.Children))
	for _, v := range n.Children {
		res[v.Key] = v
	}
	return res
}

func (*JApiCore) getPropertiesNames(m map[string]jschemaLib.ASTNode) string {
	if len(m) == 0 {
		return ""
	}

	buf := strings.Builder{}
	for k := range m {
		buf.WriteString(k)
		buf.WriteString(", ")
	}
	return strings.TrimSuffix(buf.String(), ", ")
}

func (core *JApiCore) ProcessAllOf() *jerr.JApiError {
	if je := core.processUserTypes(); je != nil {
		return je
	}

	if je := core.processBaseUrlAllOf(); je != nil {
		return je
	}

	if je := core.processRawPathVariablesAllOf(); je != nil {
		return je
	}

	if je := core.processQueryAllOf(); je != nil {
		return je
	}

	if je := core.processRequestHeaderAllOf(); je != nil {
		return je
	}

	if je := core.processRequestAllOf(); je != nil {
		return je
	}

	if je := core.processResponseHeaderAllOf(); je != nil {
		return je
	}

	return core.processResponseAllOf()
}

func (core *JApiCore) processUserTypes() *jerr.JApiError {
	return adoptError(core.catalog.UserTypes.Each(func(k string, v *catalog.UserType) error {
		if v.Schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(v.Schema.ContentJSight, v.Schema.UsedUserTypes); err != nil {
				return v.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processBaseUrlAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Servers.Each(func(k string, v *catalog.Server) error {
		s := v.BaseUrlVariables
		if s != nil && s.Schema != nil && s.Schema.Notation == notation.SchemaNotationJSight {
			if err := core.processSchemaContentJSightAllOf(s.Schema.ContentJSight, s.Schema.UsedUserTypes); err != nil {
				return s.Directive.BodyError(err.Error())
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRawPathVariablesAllOf() *jerr.JApiError {
	for _, r := range core.rawPathVariables {
		if err := core.processSchemaContentJSightAllOf(r.schema.AstNode, r.UsedUserTypes); err != nil {
			return r.pathDirective.BodyError(err.Error())
		}
	}
	return nil
}

func (core *JApiCore) processQueryAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			q := hi.Query
			if q != nil && q.Schema != nil && q.Schema.Notation == notation.SchemaNotationJSight {
				err := core.processSchemaContentJSightAllOf(q.Schema.ContentJSight, q.Schema.UsedUserTypes)
				if err != nil {
					return q.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRequestHeaderAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			r := hi.Request
			isJSight := r != nil &&
				r.HTTPRequestHeaders != nil &&
				r.HTTPRequestHeaders.Schema != nil &&
				r.HTTPRequestHeaders.Schema.Notation == notation.SchemaNotationJSight
			if isJSight {
				h := r.HTTPRequestHeaders
				err := core.processSchemaContentJSightAllOf(h.Schema.ContentJSight, h.Schema.UsedUserTypes)
				if err != nil {
					return r.HTTPRequestHeaders.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processRequestAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			r := hi.Request
			isJSight := r != nil &&
				r.HTTPRequestBody != nil &&
				r.HTTPRequestBody.Schema != nil &&
				r.HTTPRequestBody.Schema.Notation == notation.SchemaNotationJSight
			if isJSight {
				b := r.HTTPRequestBody
				err := core.processSchemaContentJSightAllOf(b.Schema.ContentJSight, b.Schema.UsedUserTypes)
				if err != nil {
					return r.HTTPRequestBody.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processResponseHeaderAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			for _, resp := range hi.Responses {
				h := resp.Headers
				if h != nil && h.Schema != nil && h.Schema.Notation == notation.SchemaNotationJSight {
					err := core.processSchemaContentJSightAllOf(h.Schema.ContentJSight, h.Schema.UsedUserTypes)
					if err != nil {
						return resp.Headers.Directive.BodyError(err.Error())
					}
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processResponseAllOf() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			for _, resp := range hi.Responses {
				b := resp.Body
				if b != nil && b.Schema != nil && b.Schema.Notation == notation.SchemaNotationJSight {
					err := core.processSchemaContentJSightAllOf(b.Schema.ContentJSight, b.Schema.UsedUserTypes)
					if err != nil {
						return resp.Body.Directive.BodyError(err.Error())
					}
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) processSchemaContentJSightAllOf(node jschemaLib.ASTNode, uut *catalog.StringSet) error {
	if node.TokenType != jschemaLib.TokenTypeObject {
		return nil
	}

	for _, v := range node.Children {
		if err := core.processSchemaContentJSightAllOf(v, uut); err != nil {
			return err
		}
	}

	rule, ok := node.Rules.Get("allOf")
	if !ok {
		return nil
	}

	switch rule.TokenType { //nolint:exhaustive // We expects only this types.
	case jschemaLib.TokenTypeArray:
		for i := len(rule.Items) - 1; i >= 0; i-- {
			r := rule.Items[i]
			if err := core.inheritPropertiesFromUserType(node, uut, r.Value); err != nil {
				return err
			}
		}
	case jschemaLib.TokenTypeShortcut:
		if err := core.inheritPropertiesFromUserType(node, uut, rule.Value); err != nil {
			return err
		}
	}
	return nil
}

func (core *JApiCore) inheritPropertiesFromUserType(
	node jschemaLib.ASTNode,
	uut *catalog.StringSet,
	userTypeName string,
) error {
	ut, ok := core.catalog.UserTypes.Get(userTypeName)
	if !ok {
		return fmt.Errorf(`the user type %q not found`, userTypeName)
	}

	if ut.Schema.ContentJSight.TokenType != jschemaLib.TokenTypeObject {
		return fmt.Errorf(`the user type %q is not an object`, userTypeName)
	}

	if _, ok := core.processedByAllOf[userTypeName]; !ok {
		core.processedByAllOf[userTypeName] = struct{}{}
		n, err := ut.Schema.JSchema.GetAST()
		if err != nil {
			return err
		}
		if err := core.processSchemaContentJSightAllOf(n, uut); err != nil {
			return err
		}
	}

	if node.Children == nil {
		node.Children = make([]jschemaLib.ASTNode, 0, 10)
	}

	for i := len(ut.Schema.ContentJSight.Children) - 1; i >= 0; i-- {
		v := ut.Schema.ContentJSight.Children[i]

		if v.Key == nil {
			return fmt.Errorf(jerr.InternalServerError)
		}

		p := node.ObjectProperty(*(v.Key))
		if p != nil && p.InheritedFrom == "" {
			// Don't allow to override original properties.
			return fmt.Errorf(
				"it is not allowed to override the %q property from the user type %q",
				*(v.Key),
				userTypeName,
			)
		}

		if p != nil && p.InheritedFrom != "" {
			// This property already defined, skip.
			continue
		}

		vv := *v
		if vv.InheritedFrom == "" {
			uut.Add(userTypeName)
		}
		vv.InheritedFrom = userTypeName
		node.Unshift(&vv)
	}

	return nil
}
