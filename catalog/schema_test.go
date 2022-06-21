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
						JsonType:    "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},
			"42 # a comment": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},
			`42 // {type: "integer"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules: &Rules{
							data: map[string]Rule{
								"type": {
									JsonType:    "string",
									ScalarValue: "integer",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
						Properties: &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`"foo"`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "string",
						Type:        "string",
						ScalarValue: "foo",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`"foo" // {type: "string"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "string",
						Type:        "string",
						ScalarValue: "foo",
						Rules: &Rules{
							data: map[string]Rule{
								"type": {
									JsonType:    "string",
									ScalarValue: "string",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
						Properties: &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"3.14": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "number",
						Type:        "float",
						ScalarValue: "3.14",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`3.14 // {type: "float"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "number",
						Type:        "float",
						ScalarValue: "3.14",
						Rules: &Rules{
							data: map[string]Rule{
								"type": {
									JsonType:    "string",
									ScalarValue: "float",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
						Properties: &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"true": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "boolean",
						Type:        "boolean",
						ScalarValue: "true",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`true // {type: "boolean"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "boolean",
						Type:        "boolean",
						ScalarValue: "true",
						Rules: &Rules{
							data: map[string]Rule{
								"type": {
									JsonType:    "string",
									ScalarValue: "boolean",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
						Properties: &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"false": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "boolean",
						Type:        "boolean",
						ScalarValue: "false",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`false // {type: "boolean"}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "boolean",
						Type:        "boolean",
						ScalarValue: "false",
						Rules: &Rules{
							data: map[string]Rule{
								"type": {
									JsonType:    "string",
									ScalarValue: "boolean",
									Properties:  &Rules{},
								},
							},
							order: []string{"type"},
						},
						Properties: &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"null": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "null",
						Type:        "null",
						ScalarValue: "null",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"42 //   some note   ": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "number",
						Type:        "integer",
						ScalarValue: "42",
						Note:        "some note",
						Rules:       &Rules{},
						Properties:  &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"42 // {min: 1, max: 100}": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules: &Rules{
							data: map[string]Rule{
								"min": {
									JsonType:    "number",
									ScalarValue: "1",
									Properties:  &Rules{},
								},
								"max": {
									JsonType:    "number",
									ScalarValue: "100",
									Properties:  &Rules{},
								},
							},
							order: []string{"min", "max"},
						},
						Properties: &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			"42 // {min: 1, max: 100} -\tsome note\t": {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "number",
						Type:        "integer",
						ScalarValue: "42",
						Rules: &Rules{
							data: map[string]Rule{
								"min": {
									JsonType:    "number",
									ScalarValue: "1",
									Properties:  &Rules{},
								},
								"max": {
									JsonType:    "number",
									ScalarValue: "100",
									Properties:  &Rules{},
								},
							},
							order: []string{"min", "max"},
						},
						Properties: &Properties{},
						Note:       "some note",
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			`"fizz" // {enum: ["fizz", "buzz"]}`: {
				expected: Schema{
					Notation: notation.SchemaNotationJSight,
					ContentJSight: &SchemaContentJSight{
						JsonType:    "string",
						Type:        "enum",
						ScalarValue: "fizz",
						Rules: &Rules{
							data: map[string]Rule{
								"enum": {
									JsonType:   "array",
									Properties: &Rules{},
									Items: []Rule{
										{
											JsonType:    "string",
											ScalarValue: "fizz",
											Properties:  &Rules{},
										},
										{
											JsonType:    "string",
											ScalarValue: "buzz",
											Properties:  &Rules{},
										},
									},
								},
							},
							order: []string{"enum"},
						},
						Properties: &Properties{},
					},
					UsedUserTypes: &StringSet{},
					UsedUserEnums: &StringSet{},
				},
			},

			// todo ERROR (code 806): An array was expected as a value for the "enum"
			// `"foo" // {enum: "@AnEnum"}`: {
			//	Notation: notation.SchemaNotationJSight,
			//	ContentJSight: &SchemaContentJSight{
			//		JsonType: "string",
			//		Type:     "enum",
			//		ScalarValue: "foo",
			//		Rules: map[string]Rule{
			//			"enum": {
			//				JsonType: "string",
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
						JsonType:   "array",
						Type:       "array",
						Rules:      &Rules{},
						Properties: &Properties{},
						Items: []*SchemaContentJSight{
							{
								JsonType:    "number",
								Type:        "integer",
								ScalarValue: "1",
								Optional:    true,
								Rules:       &Rules{},
								Properties:  &Properties{},
							},
							{
								JsonType:    "number",
								Type:        "integer",
								ScalarValue: "2",
								Optional:    true,
								Rules:       &Rules{},
								Properties:  &Properties{},
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
						JsonType:   "array",
						Type:       "array",
						Rules:      &Rules{},
						Properties: &Properties{},
						Items: []*SchemaContentJSight{
							{
								JsonType:    "shortcut",
								Type:        "@foo",
								ScalarValue: "@foo",
								Optional:    true,
								Rules:       &Rules{},
								Properties:  &Properties{},
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
						JsonType: "object",
						Type:     "object",
						Rules:    &Rules{},
						Properties: &Properties{
							data: map[string]*SchemaContentJSight{
								"foo": {
									JsonType:    "string",
									Type:        "string",
									ScalarValue: "bar",
									Rules:       &Rules{},
									Properties:  &Properties{},
								},
							},
							order: []string{"foo"},
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
						JsonType: "object",
						Type:     "object",
						Rules: &Rules{
							data: map[string]Rule{
								"additionalProperties": {
									JsonType:    "boolean",
									ScalarValue: "true",
									Properties:  &Rules{},
								},
							},
							order: []string{"additionalProperties"},
						},
						Properties: &Properties{
							data: map[string]*SchemaContentJSight{
								"foo": {
									JsonType:    "string",
									Type:        "string",
									ScalarValue: "bar",
									Optional:    true,
									Rules: &Rules{
										data: map[string]Rule{
											"optional": {
												JsonType:    "boolean",
												ScalarValue: "true",
												Properties:  &Rules{},
											},
										},
										order: []string{"optional"},
									},
									Properties: &Properties{},
								},
							},
							order: []string{"foo"},
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
						JsonType:    "number",
						Type:        "mixed",
						ScalarValue: "42",
						Rules: &Rules{
							data: map[string]Rule{
								"type": {
									JsonType:    "string",
									ScalarValue: "mixed",
									Properties:  &Rules{},
								},
								"or": {
									JsonType:   "array",
									Properties: &Rules{},
									Items: []Rule{
										{
											JsonType: "object",
											Properties: &Rules{
												data: map[string]Rule{
													"type": {
														JsonType:    "string",
														ScalarValue: "@foo",
														Properties:  &Rules{},
													},
												},
												order: []string{"type"},
											},
										},
										{
											JsonType: "object",
											Properties: &Rules{
												data: map[string]Rule{
													"type": {
														JsonType:    "string",
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
						Properties: &Properties{},
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
						JsonType:    "shortcut",
						Type:        "@foo",
						ScalarValue: "@foo",
						Rules:       &Rules{},
						Properties:  &Properties{},
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
						JsonType:    "shortcut",
						Type:        "mixed",
						ScalarValue: "@foo | @bar",
						Rules:       &Rules{},
						Properties:  &Properties{},
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
						JsonType: "object",
						Type:     "object",
						Rules: &Rules{
							data: map[string]Rule{
								"allOf": {
									JsonType:    "string",
									ScalarValue: "@foo",
									Properties:  &Rules{},
								},
							},
							order: []string{"allOf"},
						},
						Properties:    &Properties{},
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
						JsonType: "object",
						Type:     "object",
						Rules: &Rules{
							data: map[string]Rule{
								"allOf": {
									JsonType:   "array",
									Properties: &Rules{},
									Items: []Rule{
										{
											JsonType:    "string",
											ScalarValue: "@foo",
											Properties:  &Rules{},
										},
										{
											JsonType:    "string",
											ScalarValue: "@bar",
											Properties:  &Rules{},
										},
									},
								},
							},
							order: []string{"allOf"},
						},
						Properties:    &Properties{},
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
				actual, err := UnmarshalSchema("foo", []byte(b), tt, nil)
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
		//					JsonType: "object",
		//					Type:     "object",
		//					Properties: map[string]SchemaContentJSight{
		//						"page": {
		//							JsonType:    "number",
		//							Type:        "integer",
		//							ScalarValue: 1,
		//						},
		//						"per_page": {
		//							JsonType:    "number",
		//							Type:        "integer",
		//							Optional:    true,
		//							ScalarValue: 50,
		//							Rules: map[string]Rule{
		//								"optional": {
		//									JsonType:    "boolean",
		//									ScalarValue: true,
		//								},
		//								"max": {
		//									JsonType:    "number",
		//									ScalarValue: 100,
		//								},
		//							},
		//						},
		//						"filter": {
		//							JsonType: "object",
		//							Type:     "object",
		//							Optional: true,
		//							Rules: map[string]Rule{
		//								"optional": {
		//									JsonType:    "boolean",
		//									ScalarValue: true,
		//								},
		//							},
		//							Properties: map[string]SchemaContentJSight{
		//								"size": {
		//									JsonType:    "string",
		//									Type:        "enum",
		//									Optional:    true,
		//									ScalarValue: "XXL",
		//									Note:        "Filter by cat size",
		//									Rules: map[string]Rule{
		//										"optional": {
		//											JsonType:    "boolean",
		//											ScalarValue: true,
		//										},
		//										"enum": {
		//											JsonType:    "string",
		//											ScalarValue: "@catSizeEnum",
		//										},
		//									},
		//								},
		//								"age": {
		//									JsonType:    "number",
		//									Type:        "integer",
		//									Optional:    true,
		//									ScalarValue: 12,
		//									Note:        "Filter by cat's age",
		//									Rules: map[string]Rule{
		//										"optional": {
		//											JsonType:    "boolean",
		//											ScalarValue: true,
		//										},
		//									},
		//								},
		//							},
		//						},
		//						"parameters": {
		//							JsonType: "object",
		//							Type:     "object",
		//							Optional: true,
		//							Rules: map[string]Rule{
		//								"type": {
		//									JsonType:    "string",
		//									ScalarValue: "@catParameters",
		//								},
		//							},
		//						},
		//						"foo": {
		//							JsonType: "shortcut",
		//							Type:     "@foo",
		//						},
		//						"bar": {
		//							JsonType: "shortcut",
		//							Type:     "@foo | @buzz",
		//						},
		//						"@keyName": {
		//							IsKeyShortcut: true,
		//							JsonType:      "shortcut",
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
						JsonType:      "string",
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
		"jsonType": "string",
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
						JsonType:      "string",
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
					JsonType:      "string",
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
			s, err := UnmarshalSchema("", []byte(c.schema), &UserSchemas{}, nil)
			require.NoError(t, err)
			assert.Equal(t, c.expected, s.ContentJSight.Properties.order)
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
					JsonType:      "string",
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
