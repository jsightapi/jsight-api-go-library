package catalog

// TODO
// func TestSchema_MarshalJSON(t *testing.T) {
// 	t.Run("positive", func(t *testing.T) {
// 		cc := []struct {
// 			schema   ExchangeSchema
// 			expected string
// 		}{
// 			{
// 				schema: &ExchangeJSightSchema{
// 					exchangeContent: &exchangeContent{
// 						TokenType:     "string",
// 						Type:          "bar",
// 						ScalarValue:   "100",
// 						InheritedFrom: "inherited_from",
// 						Note:          "note",
// 					},
// 					exchangeUsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
// 					exchangeUsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
// 					ContentRegexp:         "/[A-z]*/",
// 				},
// 				`{
// 	"notation": "jsight",
// 	"content": {
// 		"tokenType": "string",
// 		"type": "bar",
// 		"optional": false,
// 		"scalarValue": "100",
// 		"inheritedFrom": "inherited_from",
// 		"note":	"note"
// 	},
// 	"usedUserTypes": [
// 		"@user_type_1",
// 		"@user_type_2"
// 	],
// 	"usedUserEnums": [
// 		"@user_enum_1",
// 		"@user_enum_2"
// 	]
// }`,
// 			},
//
// 			{
// 				Schema{
// 					Notation: notation.SchemaNotationRegex,
// 					ContentJSight: &SchemaContentJSight{
// 						TokenType:     "string",
// 						Type:          "bar",
// 						ScalarValue:   "100",
// 						InheritedFrom: "inherited_from",
// 						Note:          "note",
// 					},
// 					UsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
// 					UsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
// 					ContentRegexp: "/[A-z]*/",
// 				},
// 				`{
// 	"notation": "regex",
// 	"content": "/[A-z]*/"
// }`,
// 			},
// 		}
//
// 		for _, c := range cc {
// 			t.Run(string(c.schema.Notation), func(t *testing.T) {
// 				b, err := c.schema.MarshalJSON()
// 				require.NoError(t, err)
// 				assert.JSONEq(t, c.expected, string(b))
// 			})
// 		}
// 	})
//
// 	t.Run("negative", func(t *testing.T) {
// 		cc := []Schema{
// 			{
// 				Notation: notation.SchemaNotation("foo"),
// 				ContentJSight: &SchemaContentJSight{
// 					TokenType:     "string",
// 					Type:          "bar",
// 					ScalarValue:   "100",
// 					InheritedFrom: "inherited_from",
// 					Note:          "note",
// 				},
// 				UsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
// 				UsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
// 				ContentRegexp: "/[A-z]*/",
// 			},
// 		}
//
// 		for _, c := range cc {
// 			t.Run(string(c.Notation), func(t *testing.T) {
// 				_, err := c.MarshalJSON()
// 				require.Error(t, err)
// 			})
// 		}
// 	})
// }
//
// func TestSchema_MarshalJSON_Order(t *testing.T) {
// 	cc := []struct {
// 		schema   string
// 		expected []string
// 	}{
// 		{
// 			`{"id": 1, "name": "Tom", "age": 6}`,
// 			[]string{"id", "name", "age"},
// 		},
// 		{
// 			`{"id": 1, "age": 6, "name": "Tom"}`,
// 			[]string{"id", "age", "name"},
// 		},
// 		{
// 			`{"name": "Tom", "id": 1, "age": 6}`,
// 			[]string{"name", "id", "age"},
// 		},
// 		{
// 			`{"name": "Tom", "age": 6, "id": 1}`,
// 			[]string{"name", "age", "id"},
// 		},
// 		{
// 			`{"age": 6, "id": 1, "name": "Tom"}`,
// 			[]string{"age", "id", "name"},
// 		},
// 		{
// 			`{"age": 6, "name": "Tom", "id": 1}`,
// 			[]string{"age", "name", "id"},
// 		},
// 	}
//
// 	for _, c := range cc {
// 		t.Run(c.schema, func(t *testing.T) {
// 			s, err := UnmarshalJSightSchema("", []byte(c.schema), &UserSchemas{}, nil)
// 			require.NoError(t, err)
//
// 			ss := make([]string, 0, len(s.ContentJSight.Children))
// 			for _, v := range s.ContentJSight.Children {
// 				ss = append(ss, *(v.Key))
// 			}
// 			assert.Equal(t, c.expected, ss)
// 		})
// 	}
// }
//
// func BenchmarkSchema_MarshalJSON(b *testing.B) {
// 	cc := []notation.SchemaNotation{
// 		notation.SchemaNotationJSight,
// 		notation.SchemaNotationRegex,
// 		notation.SchemaNotation("foo"),
// 	}
//
// 	for _, c := range cc {
// 		b.Run(string(c), func(b *testing.B) {
// 			s := Schema{
// 				Notation: c,
// 				ContentJSight: &SchemaContentJSight{
// 					TokenType:     "string",
// 					Type:          "bar",
// 					ScalarValue:   "100",
// 					InheritedFrom: "inherited_from",
// 					Note:          "note",
// 				},
// 				UsedUserTypes: NewStringSet("@user_type_1", "@user_type_2"),
// 				UsedUserEnums: NewStringSet("@user_enum_1", "@user_enum_2"),
// 				ContentRegexp: "/[A-z]*/",
// 			}
//
// 			b.ReportAllocs()
// 			b.ResetTimer()
//
// 			for i := 0; i < b.N; i++ {
// 				_, _ = s.MarshalJSON()
// 			}
// 		})
// 	}
// }
