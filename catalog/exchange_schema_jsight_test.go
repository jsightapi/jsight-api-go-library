package catalog

import (
	"testing"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"

	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExchangeJSightSchema(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := []struct {
			project               string
			userTypes             map[string]string
			expectedUsedUserTypes *StringSet
			expectedContent       *ExchangeContent
			expectedExample       string
		}{
			{
				project: "42",
				expectedContent: &ExchangeContent{
					TokenType:   "number",
					Type:        "integer",
					ScalarValue: "42",
				},
				expectedExample: "42",
			},
			{
				project: "42 # a comment",
				expectedContent: &ExchangeContent{
					TokenType:   "number",
					Type:        "integer",
					ScalarValue: "42",
				},
				expectedExample: "42",
			},
			{
				project: `42 // {type: "integer"}`,
				expectedContent: &ExchangeContent{
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
				expectedExample: "42",
			},

			{
				project: `"foo"`,
				expectedContent: &ExchangeContent{
					TokenType:   "string",
					Type:        "string",
					ScalarValue: "foo",
				},
				expectedExample: `"foo"`,
			},

			{
				project: `"foo" // {type: "string"}`,
				expectedContent: &ExchangeContent{
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

				expectedExample: `"foo"`,
			},

			{
				project: "3.14",
				expectedContent: &ExchangeContent{
					TokenType:   "number",
					Type:        "float",
					ScalarValue: "3.14",
				},

				expectedExample: "3.14",
			},

			{
				project: `3.14 // {type: "float"}`,
				expectedContent: &ExchangeContent{
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

				expectedExample: "3.14",
			},

			{
				project: "true",
				expectedContent: &ExchangeContent{
					TokenType:   "boolean",
					Type:        "boolean",
					ScalarValue: "true",
				},

				expectedExample: "true",
			},

			{
				project: `true // {type: "boolean"}`,
				expectedContent: &ExchangeContent{
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

				expectedExample: "true",
			},

			{
				project: "false",
				expectedContent: &ExchangeContent{
					TokenType:   "boolean",
					Type:        "boolean",
					ScalarValue: "false",
				},

				expectedExample: "false",
			},

			{
				project: `false // {type: "boolean"}`,
				expectedContent: &ExchangeContent{
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

				expectedExample: "false",
			},

			{
				project: "null",
				expectedContent: &ExchangeContent{
					TokenType:   "null",
					Type:        "null",
					ScalarValue: "null",
				},

				expectedExample: "null",
			},

			{
				project: "42 //   some note   ",
				expectedContent: &ExchangeContent{
					TokenType:   "number",
					Type:        "integer",
					ScalarValue: "42",
					Note:        "some note",
				},

				expectedExample: "42",
			},

			{
				project: "42 // {min: 1, max: 100}",
				expectedContent: &ExchangeContent{
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

				expectedExample: "42",
			},

			{
				project: "42 // {min: 1, max: 100} -\tsome note\t",
				expectedContent: &ExchangeContent{
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

				expectedExample: "42",
			},

			{
				project: `"fizz" // {enum: ["fizz", "buzz"]}`,
				expectedContent: &ExchangeContent{
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

				expectedExample: `"fizz"`,
			},

			{
				project: "[1, 2]",
				expectedContent: &ExchangeContent{
					TokenType: "array",
					Type:      "array",
					Children: []*ExchangeContent{
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

				expectedExample: "[1,2]",
			},

			{
				project: "[@foo]",
				userTypes: map[string]string{
					"@foo": "42",
				},
				expectedContent: &ExchangeContent{
					TokenType: "array",
					Type:      "array",
					Children: []*ExchangeContent{
						{
							TokenType:   "reference",
							Type:        "@foo",
							ScalarValue: "@foo",
							Optional:    true,
						},
					},
				},
				expectedUsedUserTypes: NewStringSet("@foo"),
				expectedExample:       "[42]",
			},

			{
				project: `{"foo": "bar"}`,
				expectedContent: &ExchangeContent{
					TokenType: "object",
					Type:      "object",
					Children: []*ExchangeContent{
						{
							Key:         StrPtr("foo"),
							TokenType:   "string",
							Type:        "string",
							ScalarValue: "bar",
						},
					},
				},

				expectedExample: `{"foo":"bar"}`,
			},

			{
				project: `{ // {additionalProperties: true}
			"foo": "bar" // {optional: true}
		}`,
				expectedContent: &ExchangeContent{
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
					Children: []*ExchangeContent{
						{
							Key:         StrPtr("foo"),
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

				expectedExample: `{"foo":"bar"}`,
			},

			{
				project: `{"@foo": "bar", @foo: "buzz"}`,
				userTypes: map[string]string{
					"@foo": `"abc" // {minLength: 3}`,
				},
				expectedContent: &ExchangeContent{
					TokenType: "object",
					Type:      "object",
					Children: []*ExchangeContent{
						{
							Key:         StrPtr("@foo"),
							TokenType:   "string",
							Type:        "string",
							ScalarValue: "bar",
						},
						{
							Key:              StrPtr("@foo"),
							IsKeyUserTypeRef: true,
							TokenType:        "string",
							Type:             "string",
							ScalarValue:      "buzz",
						},
					},
				},
				expectedUsedUserTypes: NewStringSet("@foo"),
				expectedExample:       `{"@foo":"bar","abc":"buzz"}`,
			},

			{
				project: `42 // {type: "mixed", or: [{type: "@foo"}, {type: "@bar"}]}`,
				userTypes: map[string]string{
					"@foo": "42",
					"@bar": `"bar"`,
				},
				expectedContent: &ExchangeContent{
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
				expectedUsedUserTypes: NewStringSet("@foo", "@bar"),
				expectedExample:       "42",
			},

			{
				project: "@foo",
				userTypes: map[string]string{
					"@foo": "42",
				},
				expectedContent: &ExchangeContent{
					TokenType:   "reference",
					Type:        "@foo",
					ScalarValue: "@foo",
				},
				expectedUsedUserTypes: NewStringSet("@foo"),
				expectedExample:       "42",
			},

			{
				project: "@foo | @bar",
				userTypes: map[string]string{
					"@foo": "42",
					"@bar": `"bar"`,
				},
				expectedContent: &ExchangeContent{
					TokenType:   "reference",
					Type:        "mixed",
					ScalarValue: "@foo | @bar",
				},
				expectedUsedUserTypes: NewStringSet("@foo", "@bar"),
				expectedExample:       "42",
			},

			{
				project: `{} // {allOf: "@foo"}`,
				userTypes: map[string]string{
					"@foo": `{"id": 42}`,
				},
				expectedContent: &ExchangeContent{
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
					Children: []*ExchangeContent{
						{
							Key:           StrPtr("id"),
							TokenType:     "number",
							Type:          "integer",
							ScalarValue:   "42",
							Optional:      false,
							InheritedFrom: "@foo",
						},
					},
					InheritedFrom: "", // Handled by catalog compilation logic.
				},
				expectedUsedUserTypes: NewStringSet("@foo"),
				expectedExample:       `{"id":42}`,
			},

			{
				project: `{} // {allOf: ["@foo", "@bar"]}`,
				userTypes: map[string]string{
					"@foo": `{"foo": 42}`,
					"@bar": `{"bar": 42}`,
				},
				expectedContent: &ExchangeContent{
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
					Children: []*ExchangeContent{
						{
							Key:           StrPtr("foo"),
							TokenType:     "number",
							Type:          "integer",
							ScalarValue:   "42",
							Optional:      false,
							InheritedFrom: "@foo",
						},
						{
							Key:           StrPtr("bar"),
							TokenType:     "number",
							Type:          "integer",
							ScalarValue:   "42",
							Optional:      false,
							InheritedFrom: "@bar",
						},
					},
					InheritedFrom: "", // Handled by catalog compilation logic.
				},
				expectedUsedUserTypes: NewStringSet("@foo", "@bar"),
				expectedExample:       `{"foo":42,"bar":42}`,
			},
		}

		for _, c := range cc {
			t.Run(c.project, func(t *testing.T) {
				tt := &UserSchemas{}
				cut := &UserTypes{}
				for n, p := range c.userTypes {
					s := jschema.New(n, p)
					tt.Set(n, s)
					cut.Set(n, &UserType{
						Schema: &ExchangeJSightSchema{
							Schema: s,
						},
					})
				}
				rr := make(map[string]jschemaLib.Rule)

				es, err := NewExchangeJSightSchema("", []byte(c.project), tt, rr, cut)
				require.NoError(t, err)

				err = es.Compile()
				require.NoError(t, err)

				assert.EqualValues(t, c.expectedContent, es.exchangeContent)

				if c.expectedUsedUserTypes == nil {
					c.expectedUsedUserTypes = &StringSet{}
				}
				assert.EqualValues(t, c.expectedUsedUserTypes, es.exchangeUsedUserTypes)

				example, err := es.Example()
				require.NoError(t, err)

				assert.EqualValues(t, c.expectedExample, example)
			})
		}
	})
}
