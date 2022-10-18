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
	// if je := core.processUserTypesAllOf(); je != nil {
	// 	return je
	// }

	// TODO remove if je := core.ProcessAllOf(); je != nil {
	// 	return je
	// }

	if je := core.ExpandRawPathVariableShortcuts(); je != nil {
		return je
	}

	// TODO It's necessary? We did it early, see ExpandRawPathVariableShortcuts.
	// if je := core.CheckRawPathVariableSchemas(); je != nil {
	// 	return je
	// }

	return core.BuildResourceMethodsPathVariables()
}

func (core *JApiCore) ExpandRawPathVariableShortcuts() *jerr.JApiError {
	for i := 0; i < len(core.rawPathVariables); i++ {
		r := &core.rawPathVariables[i]

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

		if err := checkPathSchema(r.schema); err != nil {
			return r.pathDirective.KeywordError(err.Error())
		}
	}

	return nil
}

// func (core *JApiCore) CheckRawPathVariableSchemas() *jerr.JApiError {
// 	for i := 0; i < len(core.rawPathVariables); i++ {
// 		if err := checkPathSchema(core.rawPathVariables[i].schema); err != nil {
// 			return core.rawPathVariables[i].pathDirective.KeywordError(err.Error())
// 		}
// 	}
// 	return nil
// }

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

	if s.ASTNode.Children == nil || len(s.ASTNode.Children) == 0 {
		return errors.New("an empty object in the Path directive")
	}

	for _, v := range s.ASTNode.Children {
		if v.TokenType == jschemaLib.TokenTypeObject || v.TokenType == jschemaLib.TokenTypeArray {
			return fmt.Errorf("the multi-level property %q is not allowed in the Path directive", v.Key)
		}
	}

	return nil
}

func (core *JApiCore) BuildResourceMethodsPathVariables() *jerr.JApiError {
	allProjectProperties := make(map[catalog.Path]catalog.Prop)
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

				allProjectProperties[p.path] = catalog.Prop{
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
	}

	// set PathVariables
	err := core.catalog.Interactions.Map(
		func(_ catalog.InteractionID, v catalog.Interaction) (catalog.Interaction, error) {
			if hi, ok := v.(*catalog.HTTPInteraction); ok {
				pp := pathParameters(v.Path().String())
				properties := make([]catalog.Prop, 0, len(pp))

				for _, p := range pp {
					if pr, ok := allProjectProperties[p.path]; ok {
						pr.Parameter = p.parameter
						properties = append(properties, pr)
					}
				}

				if len(properties) != 0 {
					hi.SetPathVariables(catalog.NewPathVariables(properties))
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
	// if je := core.processUserTypesAllOf(); je != nil {
	// 	return je
	// }
	// TODO
	// if je := core.processBaseUrlAllOf(); je != nil {
	// 	return je
	// }
	//
	// if je := core.processRawPathVariablesAllOf(); je != nil {
	// 	return je
	// }
	//
	// if je := core.processQueryAllOf(); je != nil {
	// 	return je
	// }
	//
	// if je := core.processRequestHeaderAllOf(); je != nil {
	// 	return je
	// }
	//
	// if je := core.processRequestAllOf(); je != nil {
	// 	return je
	// }
	//
	// if je := core.processResponseHeaderAllOf(); je != nil {
	// 	return je
	// }
	//
	// return core.processResponseAllOf()
	return nil
}

// func (core *JApiCore) processUserTypesAllOf() *jerr.JApiError {
// 	return adoptError(core.catalog.UserTypes.Each(func(k string, v *catalog.UserType) error {
// 		if s, ok := v.Schema.(*catalog.ExchangeJSightSchema); ok {
// 			content, err := s.Compile()
// 			if err != nil {
// 				return v.Directive.BodyError(err.Error())
// 			}
// 			if err = core.processExchangeContentAllOf(content, s.ExchangeUsedUserTypes); err != nil {
// 				return v.Directive.BodyError(err.Error())
// 			}
// 		}
// 		return nil
// 	}))
// }

// func (core *JApiCore) processBaseUrlAllOf() *jerr.JApiError {
// 	return adoptError(core.catalog.Servers.Each(func(k string, v *catalog.Server) error {
// 		s := v.BaseUrlVariables
// 		if s != nil && s.Schema != nil {
// 			if err := core.processExchangeContentAllOf(&s.Schema.ASTNode); err != nil {
// 				return s.Directive.BodyError(err.Error())
// 			}
// 		}
// 		return nil
// 	}))
// }

// func (core *JApiCore) processRawPathVariablesAllOf() *jerr.JApiError {
// 	for _, r := range core.rawPathVariables {
// 		if err := core.processExchangeContentAllOf(&r.schema.ASTNode); err != nil {
// 			return r.pathDirective.BodyError(err.Error())
// 		}
// 	}
// 	return nil
// }

// func (core *JApiCore) processQueryAllOf() *jerr.JApiError {
// 	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
// 		if hi, ok := v.(*catalog.HTTPInteraction); ok {
// 			q := hi.Query
// 			if q != nil && q.Schema != nil {
// 				err := core.processExchangeContentAllOf(&q.Schema.ASTNode)
// 				if err != nil {
// 					return q.Directive.BodyError(err.Error())
// 				}
// 			}
// 		}
// 		return nil
// 	}))
// }

// func (core *JApiCore) processRequestHeaderAllOf() *jerr.JApiError {
// 	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
// 		if hi, ok := v.(*catalog.HTTPInteraction); ok {
// 			r := hi.Request
// 			if r != nil && r.HTTPRequestHeaders != nil && r.HTTPRequestHeaders.Schema != nil {
// 				h := r.HTTPRequestHeaders
// 				err := core.processExchangeContentAllOf(&h.Schema.ASTNode)
// 				if err != nil {
// 					return r.HTTPRequestHeaders.Directive.BodyError(err.Error())
// 				}
// 			}
// 		}
// 		return nil
// 	}))
// }

// func (core *JApiCore) processRequestAllOf() *jerr.JApiError {
// 	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
// 		if hi, ok := v.(*catalog.HTTPInteraction); ok {
// 			r := hi.Request
// 			if r != nil && r.HTTPRequestBody != nil && r.HTTPRequestBody.Schema != nil {
// 				if s, ok := r.HTTPRequestBody.Schema.(*catalog.ExchangeJSightSchema); ok {
// 					err := core.processExchangeContentAllOf(&s.ASTNode)
// 					if err != nil {
// 						return r.HTTPRequestBody.Directive.BodyError(err.Error())
// 					}
// 				}
// 			}
// 		}
// 		return nil
// 	}))
// }

// func (core *JApiCore) processResponseHeaderAllOf() *jerr.JApiError {
// 	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
// 		if hi, ok := v.(*catalog.HTTPInteraction); ok {
// 			for _, resp := range hi.Responses {
// 				h := resp.Headers
// 				if h != nil && h.Schema != nil {
// 					err := core.processExchangeContentAllOf(&h.Schema.ASTNode)
// 					if err != nil {
// 						return resp.Headers.Directive.BodyError(err.Error())
// 					}
// 				}
// 			}
// 		}
// 		return nil
// 	}))
// }

// func (core *JApiCore) processResponseAllOf() *jerr.JApiError {
// 	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
// 		if hi, ok := v.(*catalog.HTTPInteraction); ok {
// 			for _, resp := range hi.Responses {
// 				b := resp.Body
// 				if b != nil && b.Schema != nil {
// 					if s, ok := b.Schema.(*catalog.ExchangeJSightSchema); ok {
// 						err := core.processExchangeContentAllOf(&s.ASTNode)
// 						if err != nil {
// 							return resp.Body.Directive.BodyError(err.Error())
// 						}
// 					}
// 				}
// 			}
// 		}
// 		return nil
// 	}))
// }
