package catalog

import (
	"encoding/json"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"strconv"
)

type ExchangeSchemaContentJSight struct {
	// Key is key of object element.
	Key *string

	// TokenType a JSON type.
	TokenType string

	// Type a JSight type.
	Type string

	// ScalarValue contains scalar value from the example.
	// Make sense only for scalar types like string, integer, and etc.
	ScalarValue string

	// InheritedFrom a user defined type from which this property is inherited.
	InheritedFrom string

	// Note a user note.
	Note string

	// Rules a list of attached rules.
	Rules *Rules

	// Children represent available object properties or array items.
	Children []*ExchangeSchemaContentJSight

	// IsKeyUserTypeRef indicates that this is an object property which is described
	// by user defined type.
	IsKeyUserTypeRef bool

	// Optional indicates that this schema item is option or not.
	Optional bool
}

var (
	_ json.Marshaler = ExchangeSchemaContentJSight{}
	_ json.Marshaler = &ExchangeSchemaContentJSight{}
)

func (c ExchangeSchemaContentJSight) MarshalJSON() (b []byte, err error) {
	switch c.TokenType {
	case jschemaLib.TokenTypeObject, jschemaLib.TokenTypeArray:
		b, err = c.marshalJSONObjectOrArray()

	default:
		b, err = c.marshalJSONLiteral()
	}
	return b, err
}

func (c ExchangeSchemaContentJSight) marshalJSONObjectOrArray() ([]byte, error) {
	var data struct {
		Rules            []Rule                         `json:"rules,omitempty"`
		Key              *string                        `json:"key,omitempty"`
		TokenType        string                         `json:"tokenType,omitempty"`
		Type             string                         `json:"type,omitempty"`
		InheritedFrom    string                         `json:"inheritedFrom,omitempty"`
		Note             string                         `json:"note,omitempty"`
		Children         []*ExchangeSchemaContentJSight `json:"children"`
		IsKeyUserTypeRef bool                           `json:"isKeyUserTypeRef,omitempty"`
		Optional         bool                           `json:"optional"`
	}

	data.Key = c.Key
	data.IsKeyUserTypeRef = c.IsKeyUserTypeRef
	data.TokenType = c.TokenType
	data.Type = c.Type
	data.Optional = c.Optional
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() != 0 {
		data.Rules = c.Rules.data
	}
	if len(c.Children) == 0 {
		data.Children = make([]*ExchangeSchemaContentJSight, 0)
	} else {
		data.Children = c.Children
	}

	return json.Marshal(data)
}

func (c ExchangeSchemaContentJSight) marshalJSONLiteral() ([]byte, error) {
	var data struct {
		Note             string  `json:"note,omitempty"`
		Key              *string `json:"key,omitempty"`
		TokenType        string  `json:"tokenType,omitempty"`
		Type             string  `json:"type,omitempty"`
		ScalarValue      string  `json:"scalarValue"`
		InheritedFrom    string  `json:"inheritedFrom,omitempty"`
		Rules            []Rule  `json:"rules,omitempty"`
		IsKeyUserTypeRef bool    `json:"isKeyUserTypeRef,omitempty"`
		Optional         bool    `json:"optional"`
	}

	data.Key = c.Key
	data.IsKeyUserTypeRef = c.IsKeyUserTypeRef
	data.TokenType = c.TokenType
	data.Type = c.Type
	data.Optional = c.Optional
	data.ScalarValue = c.ScalarValue
	data.InheritedFrom = c.InheritedFrom
	data.Note = c.Note
	if c.Rules != nil && c.Rules.Len() != 0 {
		data.Rules = c.Rules.data
	}

	return json.Marshal(data)
}

func astNodeToJsightContent(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) *ExchangeSchemaContentJSight {
	rules := collectJSightContentRules(node, usedUserTypes)

	var isOptional bool
	if r, ok := rules.Get("optional"); ok {
		var err error
		isOptional, err = strconv.ParseBool(r.ScalarValue)
		if err != nil {
			// Normally this shouldn't happen.
			panic(err)
		}
	}

	if rules.Len() == 0 {
		rules = nil
	}

	c := &ExchangeSchemaContentJSight{
		IsKeyUserTypeRef: node.IsKeyShortcut,
		TokenType:        node.TokenType,
		Type:             node.SchemaType,
		Optional:         isOptional,
		ScalarValue:      node.Value,
		InheritedFrom:    node.InheritedFrom,
		Note:             Annotation(node.Comment),
		Rules:            rules,
	}

	switch node.TokenType {
	case jschemaLib.TokenTypeObject:
		c.collectJSightContentObjectProperties(node, usedUserTypes, usedUserEnums)
	case jschemaLib.TokenTypeArray:
		c.collectJSightContentArrayItems(node, usedUserTypes, usedUserEnums)
	}

	return c
}

func collectJSightContentRules(node jschemaLib.ASTNode, usedUserTypes *StringSet) *Rules {
	if node.Rules.Len() == 0 {
		return &Rules{}
	}

	rr := newRulesBuilder(node.Rules.Len())

	node.Rules.EachSafe(func(k string, v jschemaLib.RuleASTNode) {
		switch k {
		case "type":
			if v.Value[0] == '@' {
				usedUserTypes.Add(v.Value)
			}
			if v.Source == jschemaLib.RuleASTNodeSourceGenerated {
				return
			}

		case "allOf":
			if v.Value != "" {
				usedUserTypes.Add(v.Value)
			}

			for _, i := range v.Items {
				usedUserTypes.Add(i.Value)
			}

		case "additionalProperties":
			if v.Value[0] == '@' {
				usedUserTypes.Add(v.Value)
			}

		case "or":
			for _, i := range v.Items {
				var userType string
				if i.Value != "" {
					userType = i.Value
				} else {
					v, ok := i.Properties.Get("type")
					if ok {
						userType = v.Value
					} else {
						userType = node.SchemaType
					}
				}

				if userType[0] != '@' {
					continue
				}

				usedUserTypes.Add(userType)
			}

			if v.Source == jschemaLib.RuleASTNodeSourceGenerated {
				return
			}
		}
		rr.Set(k, astNodeToSchemaRule(v))
	})

	return rr.Rules()
}

func (c *ExchangeSchemaContentJSight) collectJSightContentObjectProperties(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) {
	if len(node.Children) > 0 {
		if c.Children == nil {
			c.Children = make([]*ExchangeSchemaContentJSight, 0, len(node.Children))
		}
		for _, v := range node.Children {
			an := astNodeToJsightContent(v, usedUserTypes, usedUserEnums)
			an.Key = SrtPtr(v.Key)

			c.Children = append(c.Children, an)

			if v.IsKeyShortcut {
				usedUserTypes.Add(v.Key)
			}
		}
	}
}

func (c *ExchangeSchemaContentJSight) collectJSightContentArrayItems(
	node jschemaLib.ASTNode,
	usedUserTypes, usedUserEnums *StringSet,
) {
	if len(node.Children) > 0 {
		if c.Children == nil {
			c.Children = make([]*ExchangeSchemaContentJSight, 0, len(node.Children))
		}
		for _, n := range node.Children {
			an := astNodeToJsightContent(n, usedUserTypes, usedUserEnums)
			an.Optional = true
			c.Children = append(c.Children, an)
		}
	}
}

func astNodeToSchemaRule(node jschemaLib.RuleASTNode) Rule {
	rr := newRulesBuilder(node.Properties.Len() + len(node.Items))

	if node.Properties.Len() != 0 {
		node.Properties.EachSafe(func(k string, v jschemaLib.RuleASTNode) {
			rr.Set(k, astNodeToSchemaRule(v))
		})
	}

	if len(node.Items) != 0 {
		for _, n := range node.Items {
			rr.Append(astNodeToSchemaRule(n))
		}
	}

	var children []Rule
	if rr.Rules().Len() != 0 {
		children = rr.Rules().data
	}

	return Rule{
		TokenType:   RuleTokenType(node.TokenType),
		ScalarValue: node.Value,
		Note:        node.Comment,
		Children:    children,
	}
}
