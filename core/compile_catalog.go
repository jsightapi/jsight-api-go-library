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

func (core *JApiCore) expandRawPathVariableShortcuts(r *rawPathVariable) *jerr.JApiError {
	for r.schema.ASTNode.TokenType == jschemaLib.TokenTypeShortcut {
		typeName := r.schema.ASTNode.SchemaType
		if typeName == "mixed" {
			return r.pathDirective.KeywordError("The root schema object cannot have an OR rule")
		}

		ut, ok := core.catalog.UserTypes.Get(typeName)
		if !ok {
			return r.pathDirective.KeywordError(fmt.Sprintf(`User type "%s" not found`, typeName))
		}

		r.schema = ut.Schema.(*catalog.ExchangeJSightSchema).Schema // copy schema
	}
	return nil
}

func checkPathSchema(s *jschema.Schema) error {
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

	// TODO remove condition allOf
	if (s.ASTNode.Children == nil || len(s.ASTNode.Children) == 0) && !s.ASTNode.Rules.Has("allOf") {
		return errors.New("an empty object in the Path directive")
	}

	for _, v := range s.ASTNode.Children {
		if v.TokenType == jschemaLib.TokenTypeObject || v.TokenType == jschemaLib.TokenTypeArray {
			return fmt.Errorf("the multi-level property %q is not allowed in the Path directive", v.Key)
		}
	}

	return nil
}

func (core *JApiCore) collectAllProjectProperties(v rawPathVariable) *jerr.JApiError {
	n, err := v.schema.GetAST()
	if err != nil {
		return v.pathDirective.KeywordError(jerr.InternalServerError)
	}

	pp := core.propertiesToMap(n)

	for _, p := range v.parameters {
		if nn, ok := pp[p.parameter]; ok {
			if _, ok := core.allProjectProperties[p.path]; ok {
				return v.pathDirective.KeywordError(fmt.Sprintf(
					"The parameter %q has already been defined earlier",
					p.parameter,
				))
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
		return v.pathDirective.KeywordError(fmt.Sprintf("Has unused parameters %q in schema", ss))
	}

	return nil
}

func (core *JApiCore) buildPathVariables() *jerr.JApiError {
	for _, v := range core.rawPathVariables {
		if je := core.expandRawPathVariableShortcuts(&v); je != nil {
			return je
		}

		// TODO allOf

		if err := checkPathSchema(v.schema); err != nil {
			return v.pathDirective.KeywordError(err.Error())
		}

		if je := core.collectAllProjectProperties(v); je != nil {
			return je
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
