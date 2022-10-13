package catalog

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/notation"
)

func TestUnmarshalJSightSchema(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			expected  Schema
			userTypes map[string]string
		}{
			"42": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "integer",
						ScalarValue: "42",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
			},
			"42 # a comment": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "integer",
						ScalarValue: "42",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
			},
			`42 // {type: "integer"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "type",
									TokenType:   RuleTokenTypeString,
									ScalarValue: "integer",
								},
							},
							index: map[string]int{"type": 0},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
			},

			`"foo"`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "string",
						Type:        "string",
						ScalarValue: "foo",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       `"foo"`,
				},
			},

			`"foo" // {type: "string"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "string",
						Type:        "string",
						ScalarValue: "foo",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "type",
									TokenType:   RuleTokenTypeString,
									ScalarValue: "string",
								},
							},
							index: map[string]int{"type": 0},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       `"foo"`,
				},
			},

			"3.14": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "float",
						ScalarValue: "3.14",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "3.14",
				},
			},

			`3.14 // {type: "float"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "float",
						ScalarValue: "3.14",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "type",
									TokenType:   RuleTokenTypeString,
									ScalarValue: "float",
								},
							},
							index: map[string]int{"type": 0},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "3.14",
				},
			},

			"true": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "boolean",
						Type:        "boolean",
						ScalarValue: "true",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "true",
				},
			},

			`true // {type: "boolean"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "boolean",
						Type:        "boolean",
						ScalarValue: "true",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "type",
									TokenType:   RuleTokenTypeString,
									ScalarValue: "boolean",
								},
							},
							index: map[string]int{"type": 0},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "true",
				},
			},

			"false": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "boolean",
						Type:        "boolean",
						ScalarValue: "false",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "false",
				},
			},

			`false // {type: "boolean"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "boolean",
						Type:        "boolean",
						ScalarValue: "false",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "type",
									TokenType:   RuleTokenTypeString,
									ScalarValue: "boolean",
								},
							},
							index: map[string]int{"type": 0},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "false",
				},
			},

			"null": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "null",
						Type:        "null",
						ScalarValue: "null",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "null",
				},
			},

			"42 //   some note   ": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "integer",
						ScalarValue: "42",
						Note:        "some note",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
			},

			"42 // {min: 1, max: 100}": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "min",
									TokenType:   RuleTokenTypeNumber,
									ScalarValue: "1",
								},
								{
									Key:         "max",
									TokenType:   RuleTokenTypeNumber,
									ScalarValue: "100",
								},
							},
							index: map[string]int{"min": 0, "max": 1},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
			},

			"42 // {min: 1, max: 100} -\tsome note\t": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "min",
									TokenType:   RuleTokenTypeNumber,
									ScalarValue: "1",
								},
								{
									Key:         "max",
									TokenType:   RuleTokenTypeNumber,
									ScalarValue: "100",
								},
							},
							index: map[string]int{"min": 0, "max": 1},
						},
						Note: "some note",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
			},

			`"fizz" // {enum: ["fizz", "buzz"]}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "string",
						Type:        "enum",
						ScalarValue: "fizz",
						Rules: &Rules{
							data: []Rule{
								{
									Key:       "enum",
									TokenType: RuleTokenTypeArray,
									Children: []Rule{
										{
											TokenType:   RuleTokenTypeString,
											ScalarValue: "fizz",
										},
										{
											TokenType:   RuleTokenTypeString,
											ScalarValue: "buzz",
										},
									},
								},
							},
							index: map[string]int{"enum": 0},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       `"fizz"`,
				},
			},

			"[1, 2]": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "array",
						Type:      "array",
						Children: []*SchemaContentJSight{
							{
								TokenType:   "number",
								Type:        "integer",
								ScalarValue: "1",
								Optional:    true,
							},
							{
								TokenType:   "number",
								Type:        "integer",
								ScalarValue: "2",
								Optional:    true,
							},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       "[1,2]",
				},
			},

			"[@foo]": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "array",
						Type:      "array",
						Children: []*SchemaContentJSight{
							{
								TokenType:   "reference",
								Type:        "@foo",
								ScalarValue: "@foo",
								Optional:    true,
							},
						},
					},
					UsedUserTypes: NewStringSet("@foo"),
					UsedUserEnums: &StringSet{},
					Example:       "[42]",
				},
				userTypes: map[string]string{
					"@foo": "42",
				},
			},

			`{"foo": "bar"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "object",
						Type:      "object",
						Children: []*SchemaContentJSight{
							{
								Key:         SrtPtr("foo"),
								TokenType:   "string",
								Type:        "string",
								ScalarValue: "bar",
							},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       `{"foo":"bar"}`,
				},
			},

			`{ // {additionalProperties: true}
	"foo": "bar" // {optional: true}
}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "object",
						Type:      "object",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "additionalProperties",
									TokenType:   RuleTokenTypeBoolean,
									ScalarValue: "true",
								},
							},
							index: map[string]int{"additionalProperties": 0},
						},
						Children: []*SchemaContentJSight{
							{
								Key:         SrtPtr("foo"),
								TokenType:   "string",
								Type:        "string",
								ScalarValue: "bar",
								Optional:    true,
								Rules: &Rules{
									data: []Rule{
										{
											Key:         "optional",
											TokenType:   RuleTokenTypeBoolean,
											ScalarValue: "true",
										},
									},
									index: map[string]int{"optional": 0},
								},
							},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
					Example:       `{"foo":"bar"}`,
				},
			},

			`{"@foo": "bar", @foo: "buzz"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "object",
						Type:      "object",
						Children: []*SchemaContentJSight{
							{
								Key:         SrtPtr("@foo"),
								TokenType:   "string",
								Type:        "string",
								ScalarValue: "bar",
							},
							{
								Key:              SrtPtr("@foo"),
								IsKeyUserTypeRef: true,
								TokenType:        "string",
								Type:             "string",
								ScalarValue:      "buzz",
							},
						},
					},
					UsedUserTypes: NewStringSet("@foo"),
					UsedUserEnums: &StringSet{},
					Example:       `{"@foo":"bar","abc":"buzz"}`,
				},
				userTypes: map[string]string{
					"@foo": `"abc" // {minLength: 3}`,
				},
			},

			`42 // {type: "mixed", or: [{type: "@foo"}, {type: "@bar"}]}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "mixed",
						ScalarValue: "42",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "type",
									TokenType:   RuleTokenTypeString,
									ScalarValue: "mixed",
								},
								{
									Key:       "or",
									TokenType: RuleTokenTypeArray,
									Children: []Rule{
										{
											TokenType: RuleTokenTypeObject,
											Children: []Rule{
												{
													Key:         "type",
													TokenType:   RuleTokenTypeReference,
													ScalarValue: "@foo",
												},
											},
										},
										{
											TokenType: RuleTokenTypeObject,
											Children: []Rule{
												{
													Key:         "type",
													TokenType:   RuleTokenTypeReference,
													ScalarValue: "@bar",
												},
											},
										},
									},
								},
							},
							index: map[string]int{"type": 0, "or": 1},
						},
					},
					UsedUserTypes: NewStringSet("@foo", "@bar"),
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
				userTypes: map[string]string{
					"@foo": "42",
					"@bar": `"bar"`,
				},
			},

			"@foo": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "reference",
						Type:        "@foo",
						ScalarValue: "@foo",
					},
					UsedUserTypes: NewStringSet("@foo"),
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
				userTypes: map[string]string{
					"@foo": "42",
				},
			},

			"@foo | @bar": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "reference",
						Type:        "mixed",
						ScalarValue: "@foo | @bar",
					},
					UsedUserTypes: NewStringSet("@foo", "@bar"),
					UsedUserEnums: &StringSet{},
					Example:       "42",
				},
				userTypes: map[string]string{
					"@foo": "42",
					"@bar": `"bar"`,
				},
			},

			`{} // {allOf: "@foo"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "object",
						Type:      "object",
						Rules: &Rules{
							data: []Rule{
								{
									Key:         "allOf",
									TokenType:   RuleTokenTypeReference,
									ScalarValue: "@foo",
								},
							},
							index: map[string]int{"allOf": 0},
						},
						InheritedFrom: "", // Handled by catalog compilation logic.
					},
					UsedUserTypes: NewStringSet("@foo"),
					UsedUserEnums: &StringSet{},
					Example:       `{"id":42}`,
				},
				userTypes: map[string]string{
					"@foo": `{"id": 42}`,
				},
			},

			`{} // {allOf: ["@foo", "@bar"]}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "object",
						Type:      "object",
						Rules: &Rules{
							data: []Rule{
								{
									Key:       "allOf",
									TokenType: RuleTokenTypeArray,
									Children: []Rule{
										{
											TokenType:   RuleTokenTypeReference,
											ScalarValue: "@foo",
										},
										{
											TokenType:   RuleTokenTypeReference,
											ScalarValue: "@bar",
										},
									},
								},
							},
							index: map[string]int{"allOf": 0},
						},
						InheritedFrom: "", // Handled by catalog compilation logic.
					},
					UsedUserTypes: NewStringSet("@foo", "@bar"),
					UsedUserEnums: &StringSet{},
					Example:       `{"foo":42,"bar":42}`,
				},
				userTypes: map[string]string{
					"@foo": `{"foo": 42}`,
					"@bar": `{"bar": 42}`,
				},
			},
		}

		for b, c := range cc {
			t.Run(b, func(t *testing.T) {
				tt := &UserSchemas{}
				for n, p := range c.userTypes {
					tt.Set(n, jschema.New(n, p))
				}
				actual, err := UnmarshalJSightSchema("foo", []byte(b), tt, nil)
				actual.JSchema = nil
				require.NoError(t, err)
				assert.EqualValues(t, c.expected, actual)
			})
		}
	})
}
