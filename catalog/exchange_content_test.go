package catalog

import (
	"testing"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchangeContent_processAllOf(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given              *ExchangeContent
			catalogUserTypes   *UserTypes
			expectedUUT        *StringSet
			expectedProperties []*ExchangeContent
		}{
			"not object": {
				&ExchangeContent{
					TokenType: "string",
				},
				&UserTypes{},
				NewStringSet(),
				nil,
			},

			"without allOf": {
				&ExchangeContent{
					TokenType: jschemaLib.TokenTypeObject,
					Children:  []*ExchangeContent{},
				},
				&UserTypes{},
				NewStringSet(),
				[]*ExchangeContent{},
			},

			// TODO

			// "with allOf, single": {
			// 	&ExchangeContent{
			// 		TokenType: jschemaLib.TokenTypeObject,
			// 		Children:  []*ExchangeContent{},
			// 		Rules: NewRules(
			// 			[]Rule{
			// 				{
			// 					Key:         "allOf",
			// 					TokenType:   RuleTokenTypeReference,
			// 					ScalarValue: "@foo",
			// 				},
			// 			},
			// 		),
			// 	},
			// 	func() *UserTypes {
			// 		ut := &UserTypes{}
			// 		ut.Set("@foo", &UserType{
			// 			Schema: &ExchangeJSightSchema{
			// 				ExchangeContent: &ExchangeContent{
			// 					TokenType: jschemaLib.TokenTypeObject,
			// 					Children: []*ExchangeContent{
			// 						{Key: StrPtr("foo")},
			// 						{Key: StrPtr("bar")},
			// 					},
			// 				},
			// 			},
			// 		})
			// 		return ut
			// 	}(),
			// 	NewStringSet("@foo"),
			// 	[]*ExchangeContent{
			// 		{
			// 			Key:           StrPtr("foo"),
			// 			InheritedFrom: "@foo",
			// 		},
			// 		{
			// 			Key:           StrPtr("bar"),
			// 			InheritedFrom: "@foo",
			// 		},
			// 	},
			// },

			// "with allOf, array": {
			// 	&ExchangeContent{
			// 		TokenType: jschemaLib.TokenTypeObject,
			// 		Children:  []*ExchangeContent{},
			// 		Rules: NewRules(
			// 			[]Rule{
			// 				{
			// 					Key:       "allOf",
			// 					TokenType: RuleTokenTypeArray,
			// 					Children: []Rule{
			// 						{
			// 							TokenType:   RuleTokenTypeString,
			// 							ScalarValue: "@foo",
			// 						},
			// 						{
			// 							TokenType:   RuleTokenTypeString,
			// 							ScalarValue: "@bar",
			// 						},
			// 					},
			// 				},
			// 			},
			// 		),
			// 	},
			// 	func() *UserTypes {
			// 		ut := &UserTypes{}
			// 		ut.Set("@foo", &UserType{
			// 			Schema: &ExchangeJSightSchema{
			// 				ExchangeContent: &ExchangeContent{
			// 					TokenType: jschemaLib.TokenTypeObject,
			// 					Children: []*ExchangeContent{
			// 						{Key: StrPtr("foo1")},
			// 						{Key: StrPtr("foo2")},
			// 					},
			// 				},
			// 			},
			// 		})
			// 		ut.Set("@bar", &UserType{
			// 			Schema: &ExchangeJSightSchema{
			// 				ExchangeContent: &ExchangeContent{
			// 					TokenType: jschemaLib.TokenTypeObject,
			// 					Children: []*ExchangeContent{
			// 						{Key: StrPtr("bar1")},
			// 						{Key: StrPtr("bar2")},
			// 					},
			// 				},
			// 			},
			// 		})
			// 		return ut
			// 	}(),
			// 	NewStringSet("@bar", "@foo"),
			// 	[]*ExchangeContent{
			// 		{
			// 			Key:           StrPtr("foo1"),
			// 			InheritedFrom: "@foo",
			// 		},
			// 		{
			// 			Key:           StrPtr("foo2"),
			// 			InheritedFrom: "@foo",
			// 		},
			// 		{
			// 			Key:           StrPtr("bar1"),
			// 			InheritedFrom: "@bar",
			// 		},
			// 		{
			// 			Key:           StrPtr("bar2"),
			// 			InheritedFrom: "@bar",
			// 		},
			// 	},
			// },
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				uut := NewStringSet()

				err := c.given.processAllOf(uut, c.catalogUserTypes)
				require.NoError(t, err)

				assert.Equal(t, c.expectedUUT, uut)
				assert.Equal(t, c.expectedProperties, c.given.Children)
			})
		}
	})
}

// func TestJApiCore_inheritPropertiesFromUserType(t *testing.T) {
// 	t.Run("positive", func(t *testing.T) {
// 		cc := map[string]struct {
// 			given              *ExchangeContent
// 			expectedUUT        *StringSet
// 			expectedProperties []*ExchangeContent
// 		}{
// 			"sc.properties is nil": {
// 				&ExchangeContent{},
// 				NewStringSet("foo"),
// 				[]*ExchangeContent{
// 					{
// 						Key:           StrPtr("bar"),
// 						InheritedFrom: "foo",
// 					},
// 				},
// 			},
//
// 			"sc.properties isn't nil": {
// 				&ExchangeContent{
// 					Children: []*ExchangeContent{
// 						{Key: StrPtr("fizz")},
// 						{Key: StrPtr("buzz")},
// 					},
// 				},
// 				NewStringSet("foo"),
// 				[]*ExchangeContent{
// 					{
// 						Key:           StrPtr("bar"),
// 						InheritedFrom: "foo",
// 					},
// 					{
// 						Key: StrPtr("fizz"),
// 					},
// 					{
// 						Key: StrPtr("buzz"),
// 					},
// 				},
// 			},
// 		}
//
// 		for n, c := range cc {
// 			t.Run(n, func(t *testing.T) {
// 				core := core.NewJApiCore(fs.NewFile("", `{}`))
// 				core.catalog = NewCatalog()
// 				core.UserTypes.Set("foo", &UserType{
// 					Schema: &ExchangeJSightSchema{
// 						ExchangeContent: &ExchangeContent{
// 							TokenType: jschema.TokenTypeObject,
// 							Children: []*ExchangeContent{
// 								{Key: StrPtr("bar")},
// 							},
// 						},
// 					},
// 				})
// 				uut := &StringSet{}
//
// 				err := core.inheritPropertiesFromUserType(c.given, uut, "foo")
// 				require.NoError(t, err)
//
// 				assert.Equal(t, c.expectedUUT, uut)
// 				assert.Equal(t, c.expectedProperties, c.given.Children)
// 			})
// 		}
// 	})
//
// 	t.Run("negative", func(t *testing.T) {
// 		cc := map[string]*JApiCore{
// 			`the user type "foo" not found`: {
// 				catalog: NewCatalog(),
// 			},
//
// 			`the user type "foo" is not an object`: {
// 				catalog: func() *Catalog {
// 					c := NewCatalog()
// 					c.UserTypes.Set("foo", &UserType{
// 						Schema: &ExchangeJSightSchema{
// 							ExchangeContent: &ExchangeContent{
// 								TokenType: jschema.TokenTypeString,
// 							},
// 						},
// 					})
// 					return c
// 				}(),
// 			},
// 		}
//
// 		for expected, core := range cc {
// 			t.Run(expected, func(t *testing.T) {
// 				err := core.inheritPropertiesFromUserType(nil, nil, "foo")
// 				assert.EqualError(t, err, expected)
// 			})
// 		}
//
// 		t.Run("property already set", func(t *testing.T) {
// 			core := NewJApiCore(fs.NewFile("", `{}`))
// 			core.catalog = NewCatalog()
// 			core.UserTypes.Set("foo", &UserType{
// 				Schema: &ExchangeJSightSchema{
// 					ExchangeContent: &ExchangeContent{
// 						TokenType: jschema.TokenTypeObject,
// 						Children: []*ExchangeContent{
// 							{Key: StrPtr("bar")},
// 						},
// 					},
// 				},
// 			})
// 			err := core.inheritPropertiesFromUserType(
// 				&ExchangeContent{
// 					Children: []*ExchangeContent{
// 						{Key: StrPtr("bar")},
// 					},
// 				},
// 				&StringSet{},
// 				"foo",
// 			)
// 			assert.EqualError(t, err, `it is not allowed to override the "bar" property from the user type "foo"`)
// 		})
//
// 		t.Run("core is nil", func(t *testing.T) {
// 			assert.Panics(t, func() {
// 				var core *JApiCore
// 				_ = core.inheritPropertiesFromUserType(
// 					&ExchangeContent{},
// 					&StringSet{},
// 					"foo",
// 				)
// 			})
// 		})
//
// 		t.Run("sc is nil", func(t *testing.T) {
// 			assert.Panics(t, func() {
// 				core := &JApiCore{
// 					catalog: NewCatalog(),
// 				}
// 				core.UserTypes.Set("foo", &UserType{
// 					Schema: &ExchangeJSightSchema{
// 						ExchangeContent: &ExchangeContent{
// 							TokenType: jschema.TokenTypeObject,
// 						},
// 					},
// 				})
// 				_ = core.inheritPropertiesFromUserType(
// 					nil,
// 					&StringSet{},
// 					"foo",
// 				)
// 			})
// 		})
//
// 		t.Run("uut is nil", func(t *testing.T) {
// 			assert.Panics(t, func() {
// 				core := &JApiCore{
// 					catalog: NewCatalog(),
// 				}
// 				core.UserTypes.Set("foo", &UserType{
// 					Schema: &ExchangeJSightSchema{
// 						ExchangeContent: &ExchangeContent{
// 							TokenType: jschema.TokenTypeObject,
// 							Children: []*ExchangeContent{
// 								{Key: StrPtr("foo")},
// 							},
// 						},
// 					},
// 				})
// 				_ = core.inheritPropertiesFromUserType(
// 					&ExchangeContent{},
// 					nil,
// 					"foo",
// 				)
// 			})
// 		})
// 	})
// }
