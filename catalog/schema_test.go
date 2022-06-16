package catalog

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/notation"
)

func TestUnmarshalSchema(t *testing.T) {
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
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},
			"42 # a comment": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"type": {
									TokenType:   "string",
									ScalarValue: "integer",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`"foo"`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "string",
						Type:        "string",
						ScalarValue: "foo",
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"type": {
									TokenType:   "string",
									ScalarValue: "string",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"3.14": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "number",
						Type:        "float",
						ScalarValue: "3.14",
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"type": {
									TokenType:   "string",
									ScalarValue: "float",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"true": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "boolean",
						Type:        "boolean",
						ScalarValue: "true",
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"type": {
									TokenType:   "string",
									ScalarValue: "boolean",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"false": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "boolean",
						Type:        "boolean",
						ScalarValue: "false",
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"type": {
									TokenType:   "string",
									ScalarValue: "boolean",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"null": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "null",
						Type:        "null",
						ScalarValue: "null",
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
						Rules:       &Rules{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"min": {
									TokenType:   "number",
									ScalarValue: "1",
									Properties:  &Rules{},
								},
								"max": {
									TokenType:   "number",
									ScalarValue: "100",
									Properties:  &Rules{},
								},
							},
							order: []string{"min", "max"},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"min": {
									TokenType:   "number",
									ScalarValue: "1",
									Properties:  &Rules{},
								},
								"max": {
									TokenType:   "number",
									ScalarValue: "100",
									Properties:  &Rules{},
								},
							},
							order: []string{"min", "max"},
						},
						Note: "some note",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"enum": {
									TokenType:  "array",
									Properties: &Rules{},
									Items: []Rule{
										{
											TokenType:   "string",
											ScalarValue: "fizz",
											Properties:  &Rules{},
										},
										{
											TokenType:   "string",
											ScalarValue: "buzz",
											Properties:  &Rules{},
										},
									},
								},
							},
							order: []string{"enum"},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			// todo ERROR (code 806): An array was expected as a value for the "enum"
			// `"foo" // {enum: "@AnEnum"}`: {
			//	Notation: notation.SchemaNotationJSight,
			//	ContentJSight: &SchemaContentJSight{
			//		TokenType: "string",
			//		Type:     "enum",
			//		ScalarValue: "foo",
			//		Rules: map[string]Rule{
			//			"enum": {
			//				TokenType: "string",
			//				ScalarValue: "@AnEnum",
			//			},
			//		},
			//	},
			//	UsedUserEnums: []string{"@AnEnum"},
			// },

			"[1, 2]": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "array",
						Type:      "array",
						Rules:     &Rules{},
						Children: []*SchemaContentJSight{
							{
								TokenType:   "number",
								Type:        "integer",
								ScalarValue: "1",
								Optional:    true,
								Rules:       &Rules{},
							},
							{
								TokenType:   "number",
								Type:        "integer",
								ScalarValue: "2",
								Optional:    true,
								Rules:       &Rules{},
							},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"[@foo]": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "array",
						Type:      "array",
						Rules:     &Rules{},
						Children: []*SchemaContentJSight{
							{
								TokenType:   "shortcut",
								Type:        "@foo",
								ScalarValue: "@foo",
								Optional:    true,
								Rules:       &Rules{},
							},
						},
					},
					UsedUserTypes: NewStringSet("@foo"),
					UsedUserEnums: &StringSet{},
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
						Rules:     &Rules{},
						Children: []*SchemaContentJSight{
							{
								Key:         "foo",
								TokenType:   "string",
								Type:        "string",
								ScalarValue: "bar",
								Rules:       &Rules{},
							},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"additionalProperties": {
									TokenType:   "boolean",
									ScalarValue: "true",
									Properties:  &Rules{},
								},
							},
							order: []string{"additionalProperties"},
						},
						Children: []*SchemaContentJSight{
							{
								Key:         "foo",
								TokenType:   "string",
								Type:        "string",
								ScalarValue: "bar",
								Optional:    true,
								Rules: &Rules{
									data: map[string]Rule{
										"optional": {
											TokenType:   "boolean",
											ScalarValue: "true",
											Properties:  &Rules{},
										},
									},
									order: []string{"optional"},
								},
							},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`{"@foo": "bar", @foo: "baz"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType: "object",
						Type:      "object",
						Rules:     &Rules{},
						Children: []*SchemaContentJSight{
							{
								Key:         "@foo",
								TokenType:   "string",
								Type:        "string",
								ScalarValue: "bar",
								Rules:       &Rules{},
							},
							{
								Key:         "@foo",
								TokenType:   "string",
								Type:        "string",
								ScalarValue: "baz",
								Rules:       &Rules{},
							},
						},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"type": {
									TokenType:   "string",
									ScalarValue: "mixed",
									Properties:  &Rules{},
								},
								"or": {
									TokenType:  "array",
									Properties: &Rules{},
									Items: []Rule{
										{
											TokenType: "object",
											Properties: &Rules{
												data: map[string]Rule{
													"type": {
														TokenType:   "string",
														ScalarValue: "@foo",
														Properties:  &Rules{},
													},
												},
												order: []string{"type"},
											},
										},
										{
											TokenType: "object",
											Properties: &Rules{
												data: map[string]Rule{
													"type": {
														TokenType:   "string",
														ScalarValue: "@bar",
														Properties:  &Rules{},
													},
												},
												order: []string{"type"},
											},
										},
									},
								},
							},
							order: []string{"type", "or"},
						},
					},
					UsedUserTypes: NewStringSet("@foo", "@bar"),
					UsedUserEnums: &StringSet{},
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
						TokenType:   "shortcut",
						Type:        "@foo",
						ScalarValue: "@foo",
						Rules:       &Rules{},
					},
					UsedUserTypes: NewStringSet("@foo"),
					UsedUserEnums: &StringSet{},
				},
				userTypes: map[string]string{
					"@foo": "42",
				},
			},

			"@foo | @bar": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:   "shortcut",
						Type:        "mixed",
						ScalarValue: "@foo | @bar",
						Rules:       &Rules{},
					},
					UsedUserTypes: NewStringSet("@foo", "@bar"),
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"allOf": {
									TokenType:   "string",
									ScalarValue: "@foo",
									Properties:  &Rules{},
								},
							},
							order: []string{"allOf"},
						},
						InheritedFrom: "", // Handled by catalog compilation logic.
					},
					UsedUserTypes: NewStringSet("@foo"),
					UsedUserEnums: &StringSet{},
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
							data: map[string]Rule{
								"allOf": {
									TokenType:  "array",
									Properties: &Rules{},
									Items: []Rule{
										{
											TokenType:   "string",
											ScalarValue: "@foo",
											Properties:  &Rules{},
										},
										{
											TokenType:   "string",
											ScalarValue: "@bar",
											Properties:  &Rules{},
										},
									},
								},
							},
							order: []string{"allOf"},
						},
						InheritedFrom: "", // Handled by catalog compilation logic.
					},
					UsedUserTypes: NewStringSet("@foo", "@bar"),
					UsedUserEnums: &StringSet{},
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
					tt.Set(n, jschema.New(n, []byte(p)))
				}
				actual, err := UnmarshalSchema("foo", []byte(b), tt)
				require.NoError(t, err)
				assert.EqualValues(t, c.expected, actual)
			})
		}

		// todo ERROR (code 806): An array was expected as a value for the "enum"
		//		t.Run("full", func(t *testing.T) {
		//			s, err := UnmarshalSchema("bar", []byte(`
		// { # An example.
		//	"page": 1,
		//	"per_page": 50,  // {optional: true, max: 100}
		//	"filter": {     // {optional: true}
		//		"size": "XXL" // {optional: true, enum: "@catSizeEnum"} - Filter by cat size.
		//		"age" : 12    // {optional: true                      } - Filter by age of the cat.
		//	},
		//	"parameters": {} // {type: "@catParameters"},
		//	"foo": @foo,
		//	"bar": @fizz | @buzz,
		//	@keyName: @keyValue
		// `))
		//			require.NoError(t, err)
		//
		//			assert.Equal(t, Schema{
		//				Notation: notation.SchemaNotationJSight,
		//				ContentJSight: &SchemaContentJSight{
		//					TokenType: "object",
		//					Type:     "object",
		//					Properties: map[string]SchemaContentJSight{
		//						"page": {
		//							TokenType:    "number",
		//							Type:        "integer",
		//							ScalarValue: 1,
		//						},
		//						"per_page": {
		//							TokenType:    "number",
		//							Type:        "integer",
		//							Optional:    true,
		//							ScalarValue: 50,
		//							Rules: map[string]Rule{
		//								"optional": {
		//									TokenType:    "boolean",
		//									ScalarValue: true,
		//								},
		//								"max": {
		//									TokenType:    "number",
		//									ScalarValue: 100,
		//								},
		//							},
		//						},
		//						"filter": {
		//							TokenType: "object",
		//							Type:     "object",
		//							Optional: true,
		//							Rules: map[string]Rule{
		//								"optional": {
		//									TokenType:    "boolean",
		//									ScalarValue: true,
		//								},
		//							},
		//							Properties: map[string]SchemaContentJSight{
		//								"size": {
		//									TokenType:    "string",
		//									Type:        "enum",
		//									Optional:    true,
		//									ScalarValue: "XXL",
		//									Note:        "Filter by cat size",
		//									Rules: map[string]Rule{
		//										"optional": {
		//											TokenType:    "boolean",
		//											ScalarValue: true,
		//										},
		//										"enum": {
		//											TokenType:    "string",
		//											ScalarValue: "@catSizeEnum",
		//										},
		//									},
		//								},
		//								"age": {
		//									TokenType:    "number",
		//									Type:        "integer",
		//									Optional:    true,
		//									ScalarValue: 12,
		//									Note:        "Filter by cat's age",
		//									Rules: map[string]Rule{
		//										"optional": {
		//											TokenType:    "boolean",
		//											ScalarValue: true,
		//										},
		//									},
		//								},
		//							},
		//						},
		//						"parameters": {
		//							TokenType: "object",
		//							Type:     "object",
		//							Optional: true,
		//							Rules: map[string]Rule{
		//								"type": {
		//									TokenType:    "string",
		//									ScalarValue: "@catParameters",
		//								},
		//							},
		//						},
		//						"foo": {
		//							TokenType: "shortcut",
		//							Type:     "@foo",
		//						},
		//						"bar": {
		//							TokenType: "shortcut",
		//							Type:     "@foo | @buzz",
		//						},
		//						"@keyName": {
		//							IsKeyUserTypeRef: true,
		//							TokenType:      "shortcut",
		//							Type:          "@keyValue",
		//						},
		//					},
		//				},
		//				UsedUserTypes: []string{
		//					"@catParameters",
		//					"@foo",
		//					"@fizz",
		//					"@buzz",
		//					"@keyName",
		//					"@keyValue",
		//				},
		//				UsedUserEnums: []string{"@catSizeEnum"},
		//			}, s)
		//		})
	})
}

func TestSchema_MarshalJSON(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := []struct {
			schema   Schema
			expected string
		}{
			{
				Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						TokenType:     "string",
						Type:          "bar",
						ScalarValue:   "100",
						InheritedFrom: "inherited_from",
						Note:          "note",
					},
					UsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
					UsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
					ContentRegexp: "/[A-z]*/",
				},
				`{
	"notation": "jsight",
	"content": {
		"tokenType": "string",
		"type": "bar",
		"optional": false,
		"scalarValue": "100",
		"inheritedFrom": "inherited_from",
		"note":	"note"
	},
	"usedUserTypes": [
		"@user_type_1",
		"@user_type_2"
	],
	"usedUserEnums": [
		"@user_enum_1",
		"@user_enum_2"
	]
}`,
			},

			{
				Schema{
					Notation: notation.SchemaNotationRegex,
					ContentJSight: &SchemaContentJSight{
						TokenType:     "string",
						Type:          "bar",
						ScalarValue:   "100",
						InheritedFrom: "inherited_from",
						Note:          "note",
					},
					UsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
					UsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
					ContentRegexp: "/[A-z]*/",
				},
				`{
	"notation": "regex",
	"content": "/[A-z]*/"
}`,
			},
		}

		for _, c := range cc {
			t.Run(string(c.schema.Notation), func(t *testing.T) {
				b, err := c.schema.MarshalJSON()
				require.NoError(t, err)
				assert.JSONEq(t, c.expected, string(b))
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := []Schema{
			{
				Notation: notation.SchemaNotation("foo"),
				ContentJSight: &SchemaContentJSight{
					TokenType:     "string",
					Type:          "bar",
					ScalarValue:   "100",
					InheritedFrom: "inherited_from",
					Note:          "note",
				},
				UsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
				UsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
				ContentRegexp: "/[A-z]*/",
			},
		}

		for _, c := range cc {
			t.Run(string(c.Notation), func(t *testing.T) {
				_, err := c.MarshalJSON()
				require.Error(t, err)
			})
		}
	})
}

func TestSchema_MarshalJSON_Order(t *testing.T) {
	cc := []struct {
		schema   string
		expected []string
	}{
		{
			`{"id": 1, "name": "Tom", "age": 6}`,
			[]string{"id", "name", "age"},
		},
		{
			`{"id": 1, "age": 6, "name": "Tom"}`,
			[]string{"id", "age", "name"},
		},
		{
			`{"name": "Tom", "id": 1, "age": 6}`,
			[]string{"name", "id", "age"},
		},
		{
			`{"name": "Tom", "age": 6, "id": 1}`,
			[]string{"name", "age", "id"},
		},
		{
			`{"age": 6, "id": 1, "name": "Tom"}`,
			[]string{"age", "id", "name"},
		},
		{
			`{"age": 6, "name": "Tom", "id": 1}`,
			[]string{"age", "name", "id"},
		},
	}

	for _, c := range cc {
		t.Run(c.schema, func(t *testing.T) {
			s, err := UnmarshalSchema("", []byte(c.schema), &UserSchemas{})
			require.NoError(t, err)

			ss := make([]string, 0, len(s.ContentJSight.Children))
			for _, v := range s.ContentJSight.Children {
				ss = append(ss, v.Key)
			}
			assert.Equal(t, c.expected, ss)
		})
	}
}

func BenchmarkSchema_MarshalJSON(b *testing.B) {
	cc := []notation.SchemaNotation{
		notation.SchemaNotationJSight,
		notation.SchemaNotationRegex,
		notation.SchemaNotation("foo"),
	}

	for _, c := range cc {
		b.Run(string(c), func(b *testing.B) {
			s := Schema{
				Notation: c,
				ContentJSight: &SchemaContentJSight{
					TokenType:     "string",
					Type:          "bar",
					ScalarValue:   "100",
					InheritedFrom: "inherited_from",
					Note:          "note",
				},
				UsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
				UsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
				ContentRegexp: "/[A-z]*/",
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = s.MarshalJSON()
			}
		})
	}
}
