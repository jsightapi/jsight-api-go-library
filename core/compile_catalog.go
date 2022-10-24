package core

import (
	"errors"
	"fmt"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"strings"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (core *JApiCore) compileCatalog() *jerr.JApiError {
	return core.buildPathVariables()
}

func (core *JApiCore) expandRawPathVariableShortcuts(r *rawPathVariable) error {
	for r.schema.ASTNode.TokenType == jschemaLib.TokenTypeShortcut {
		typeName := r.schema.ASTNode.SchemaType
		if typeName == "mixed" {
			return errors.New("The root schema object cannot have an OR rule")
		}

		ut, ok := core.catalog.UserTypes.Get(typeName)
		if !ok {
			return fmt.Errorf(`%s (%s)`, jerr.UserTypeNotFound, typeName)
		}

		r.schema = ut.Schema.(*catalog.ExchangeJSightSchema).Schema // copy schema
	}
	return nil
}

func (core *JApiCore) expandRawPathVariableAllOf(r *rawPathVariable) error {
	return core.expandASTNodeAllOf(&r.schema.ASTNode, r.schema)
}

func (core *JApiCore) expandASTNodeAllOf(an *jschemaLib.ASTNode, s *jschema.Schema) error {
	for i := range an.Children {
		if err := core.expandASTNodeAllOf(&(an.Children[i]), s); err != nil {
			return err
		}
	}

	allOf, ok := an.Rules.Get("allOf")
	if !ok {
		return nil
	}

	switch allOf.TokenType { //nolint:exhaustive // We expects only this types.
	case jschemaLib.TokenTypeArray:
		for i := len(allOf.Items) - 1; i >= 0; i-- {
			r := allOf.Items[i]
			if err := core.inheritPropertiesFromASTNode(an, r.Value, s); err != nil {
				return err
			}
		}
	case jschemaLib.TokenTypeShortcut:
		if err := core.inheritPropertiesFromASTNode(an, allOf.Value, s); err != nil {
			return err
		}
	}

	return nil
}

func (core *JApiCore) inheritPropertiesFromASTNode(
	an *jschemaLib.ASTNode,
	userTypeName string,
	s *jschema.Schema,
) error {
	ut, ok := core.catalog.UserTypes.Get(userTypeName)
	if !ok {
		return fmt.Errorf(`the user type %q not found`, userTypeName)
	}

	utn, err := ut.Schema.GetAST()
	if err != nil {
		return err
	}

	err = core.expandASTNodeAllOf(&utn, s)
	if err != nil {
		return err
	}

	if utn.TokenType != jschemaLib.TokenTypeObject {
		return fmt.Errorf(`%s (%s)`, jerr.UserTypeIsNotAnObject, userTypeName)
	}

	if an.Children == nil {
		an.Children = make([]jschemaLib.ASTNode, 0, 10)
	}

	for i := len(utn.Children) - 1; i >= 0; i-- {
		child := utn.Children[i]

		if child.Key == "" {
			return errors.New(jerr.InternalServerError)
		}

		p := an.ObjectProperty(child.Key)
		if p != nil && p.InheritedFrom == "" {
			// Don't allow to override original properties.
			return fmt.Errorf(jerr.NotAllowedToOverrideTheProperty,
				child.Key,
				userTypeName,
			)
		}

		if p != nil && p.InheritedFrom != "" {
			// This property already defined, skip.
			continue
		}

		dup := child
		if dup.InheritedFrom == "" {
			s.AddUserTypeName(userTypeName)
		}
		dup.InheritedFrom = userTypeName
		an.Children = append(an.Children, dup)
	}

	return nil
}

func (core *JApiCore) checkPathSchema(s *jschema.Schema) error {
	if err := core.checkPathSchemaRoot(s); err != nil {
		return err
	}

	for _, an := range s.ASTNode.Children {
		if err := core.checkPathSchemaProperty(an); err != nil {
			return err
		}
	}

	return nil
}

func (core *JApiCore) checkPathSchemaRoot(s *jschema.Schema) error {
	if s.ASTNode.TokenType != jschemaLib.TokenTypeObject {
		return errors.New("the body of the Path DIRECTIVE must be an object")
	}

	if s.ASTNode.Rules.Has("additionalProperties") {
		return errors.New(`the "additionalProperties" rule is invalid in the Path directive`)
	}

	if s.ASTNode.Rules.Has("nullable") {
		return errors.New(`the "nullable" rule is invalid in the Path directive`)
	}

	if s.ASTNode.Rules.Has("or") {
		return errors.New(`the "or" rule is invalid in the Path directive`)
	}

	if s.ASTNode.Children == nil || len(s.ASTNode.Children) == 0 {
		return errors.New("an empty object in the Path directive")
	}

	return nil
}

func (core *JApiCore) checkPathSchemaProperty(an jschemaLib.ASTNode) error {
	if an.TokenType == jschemaLib.TokenTypeObject || an.TokenType == jschemaLib.TokenTypeArray {
		return fmt.Errorf("%s (%s)", jerr.MultiLevelPropertyIsNotAllowed, an.Key)
	}

	rule, ok := an.Rules.Get("or")
	if ok {
		for _, v := range rule.Items {
			switch v.TokenType {
			case jschemaLib.TokenTypeShortcut:
				if err := core.checkPathSchemaPropertyUserType(v.Value); err != nil {
					return err
				}
			case jschemaLib.TokenTypeObject:
				if t, ok := v.Properties.Get("type"); ok {
					if t.TokenType == "string" && (t.Value == jschemaLib.TokenTypeObject ||
						t.Value == jschemaLib.TokenTypeArray) {
						return fmt.Errorf("%s (%s)", jerr.MultiLevelPropertyIsNotAllowed, an.Key)
					}
				}
			}
		}
	} else if an.TokenType == jschemaLib.TokenTypeShortcut {
		if err := core.checkPathSchemaPropertyUserType(an.Value); err != nil {
			return err
		}
	}

	return nil
}

func (core *JApiCore) checkPathSchemaPropertyUserType(typeName string) error {
	ut, ok := core.catalog.UserTypes.Get(typeName)
	if !ok {
		return fmt.Errorf(`%s (%s)`, jerr.UserTypeNotFound, typeName)
	}

	rootNode, err := ut.Schema.GetAST()
	if err != nil {
		return errors.New(jerr.InternalServerError)
	}

	if err := core.checkPathSchemaProperty(rootNode); err != nil {
		return err
	}

	return nil
}

func (core *JApiCore) collectAllProjectProperties(v *rawPathVariable) error {
	n, err := v.schema.GetAST()
	if err != nil {
		return errors.New(jerr.InternalServerError)
	}

	pp := core.propertiesToMap(n)

	for _, p := range v.parameters {
		if nn, ok := pp[p.parameter]; ok {
			if _, ok := core.allProjectProperties[p.path]; ok {
				return fmt.Errorf(
					"The parameter %q has already been defined earlier",
					p.parameter,
				)
			}

			core.allProjectProperties[p.path] = catalog.Prop{
				ASTNode:   nn,
				Directive: v.pathDirective,
			}

			delete(pp, p.parameter)
		}
	}

	// Check that all path properties in schema is exists in the path.
	if len(pp) > 0 {
		ss := core.getPropertiesNames(pp)
		return fmt.Errorf("Has unused parameters %q in schema", ss)
	}

	return nil
}

func (core *JApiCore) buildPathVariables() *jerr.JApiError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		v := &core.rawPathVariables[i]

		if err := core.expandRawPathVariableShortcuts(v); err != nil {
			return v.pathDirective.KeywordError(err.Error())
		}

		if err := core.expandRawPathVariableAllOf(v); err != nil {
			return v.pathDirective.KeywordError(err.Error())
		}

		if err := core.checkPathSchema(v.schema); err != nil {
			return v.pathDirective.KeywordError(err.Error())
		}

		if err := core.collectAllProjectProperties(v); err != nil {
			return v.pathDirective.KeywordError(err.Error())
		}
	}

	// set PathVariables
	err := core.catalog.Interactions.Map(
		func(_ catalog.InteractionID, v catalog.Interaction) (catalog.Interaction, error) {
			if hi, ok := v.(*catalog.HTTPInteraction); ok {
				pp := pathParameters(v.Path().String())
				properties := make([]catalog.Prop, 0, len(pp))

				for _, p := range pp {
					if pr, ok := core.allProjectProperties[p.path]; ok {
						pr.Parameter = p.parameter
						properties = append(properties, pr)
					}
				}

				if len(properties) != 0 {
					hi.SetPathVariables(catalog.NewPathVariables(properties, core.catalog.UserTypes))
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
