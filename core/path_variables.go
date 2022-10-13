package core

import (
	"fmt"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/jsightapi/jsight-schema-go-library/notations/regex"
	"strings"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type prop struct {
	parameter string
	astNode   jschemaLib.ASTNode
	directive directive.Directive
}

func (core *JApiCore) newPathVariables(properties []prop) *catalog.PathVariables {
	n := &jschemaLib.ASTNodes{}

	for _, p := range properties {
		n.Set(p.parameter, p.astNode)
	}

	return &catalog.PathVariables{ASTNodes: n}
}

func (core *JApiCore) collectUsedUserTypes(sc *catalog.SchemaContentJSight, usedUserTypes *catalog.StringSet) error {
	if sc.TokenType == jschemaLib.TokenTypeShortcut {
		// We have two different cases under "reference" type:
		// 1. Single type like "@foo"
		// 2. A list of types like "@foo | @bar"
		//
		// For the first case we have valid user type in the `v.Type` property.
		// But for the second case we got "mixed" there. So we should use `v.ScalarValue`
		// instead. This property should always be string.
		for _, t := range strings.Split(sc.ScalarValue, " | ") {
			if err := core.appendUsedUserType(usedUserTypes, t); err != nil {
				return err
			}
		}
		return nil
	}
	err := sc.Rules.Each(func(k string, v catalog.Rule) error {
		switch k {
		case "type":
			if v.ScalarValue[0] == '@' {
				if err := core.appendUsedUserType(usedUserTypes, v.ScalarValue); err != nil {
					return err
				}
			}

		case "or":
			for _, i := range v.Children {
				var userType string
				if i.ScalarValue != "" {
					userType = i.ScalarValue
				} else {
					for _, c := range i.Children {
						if c.Key == "type" {
							userType = c.ScalarValue
							break
						}
					}
				}

				// Schema types shouldn't be added.
				if jschemaLib.IsValidType(userType) {
					continue
				}

				if err := core.appendUsedUserType(usedUserTypes, userType); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (core *JApiCore) appendUsedUserType(usedUserTypes *catalog.StringSet, s string) error {
	if t, ok := core.catalog.UserTypes.Get(s); ok {
		switch sc := t.Schema.(type) {
		case *jschema.Schema:
			switch sc.AstNode.TokenType {
			case "string", "number", "boolean", "null":
				usedUserTypes.Add(s)
				return nil
			default:
				return fmt.Errorf(
					"unavailable JSON type %q of the UserType %q in Path directive",
					sc.AstNode.TokenType,
					s,
				)
			}
		case *regex.Schema:
			usedUserTypes.Add(s)
			return nil
		default:
			// case notation.SchemaNotationAny, notation.SchemaNotationEmpty:
			// return err (see below)
		}
	}
	return fmt.Errorf(`UserType not found "%s" for Path directive`, s)
}
