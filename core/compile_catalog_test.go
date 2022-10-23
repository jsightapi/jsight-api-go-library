package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/catalog"
)

func TestJApiCore_buildPathVariables(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given    *JApiCore
			expected *JApiCore
		}{
			"empty": {
				&JApiCore{
					catalog: &catalog.Catalog{
						Interactions: &catalog.Interactions{},
					},
				},
				&JApiCore{
					catalog: &catalog.Catalog{
						Interactions: &catalog.Interactions{},
					},
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				err := c.given.buildPathVariables()
				require.Nil(t, err)

				assert.Equal(t, c.expected, c.given)
			})
		}
	})

	// TODO t.Run("negative", func(t *testing.T) {
	// 	cc := map[string]*JApiCore{
	// 		`Has unused parameters "bar" in schema`: {
	// 			rawPathVariables: []rawPathVariable{
	// 				{
	// 					schema: catalog.Schema{
	// 						ExchangeContent: &catalog.ExchangeContent{
	// 							Children: []*catalog.ExchangeContent{
	// 								{Key: catalog.SrtPtr("foo")},
	// 								{Key: catalog.SrtPtr("bar")},
	// 							},
	// 						},
	// 					},
	// 					parameters: []PathParameter{
	// 						{
	// 							path:      "/path/{foo}/{bar}",
	// 							parameter: "foo",
	// 						},
	// 					},
	// 					pathDirective: *directive.New(
	// 						directive.Get,
	// 						directive.NewCoords(
	// 							fs.NewFile("foo", "/path/{foo}/{bar}"),
	// 							0,
	// 							5,
	// 						),
	// 					),
	// 				},
	// 			},
	//
	// 			catalog: &catalog.Catalog{
	// 				Interactions: &catalog.Interactions{},
	// 			},
	// 		},
	// 	}
	//
	// 	for expected, core := range cc {
	// 		t.Run(expected, func(t *testing.T) {
	// 			err := core.buildPathVariables()
	// 			assert.EqualError(t, err, expected)
	// 		})
	// 	}
	//
	// 	t.Run("nil", func(t *testing.T) {
	// 		assert.Panics(t, func() {
	// 			var core *JApiCore
	//
	// 			_ = core.buildPathVariables()
	// 		})
	// 	})
	// })
}

// TODO func TestCore_propertiesToMap(t *testing.T) {
// 	cc := map[string]struct {
// 		given    *catalog.Properties
// 		expected map[string]*catalog.ExchangeContent
// 	}{
// 		"nil": {},
//
// 		"empty": {
// 			given: &catalog.Properties{},
// 		},
//
// 		"with data": {
// 			given: catalog.NewProperties(
// 				map[string]*catalog.ExchangeContent{
// 					"foo": nil,
// 					"bar": {
// 						Note: "fake",
// 					},
// 				},
// 				[]string{"foo", "bar"},
// 			),
// 			expected: map[string]*catalog.ExchangeContent{
// 				"foo": nil,
// 				"bar": {
// 					Note: "fake",
// 				},
// 			},
// 		},
// 	}
//
// 	for n, c := range cc {
// 		t.Run(n, func(t *testing.T) {
// 			actual := (&JApiCore{}).propertiesToMap(c.given)
// 			assert.Equal(t, c.expected, actual)
// 		})
// 	}
// }

// TODO func TestCore_getPropertiesNames(t *testing.T) {
// 	cc := map[string]struct {
// 		given    map[string]*catalog.ExchangeContent
// 		expected string
// 	}{
// 		"nil": {},
//
// 		"empty": {
// 			given: map[string]*catalog.ExchangeContent{},
// 		},
//
// 		"with data": {
// 			given: map[string]*catalog.ExchangeContent{
// 				"foo":  nil,
// 				"bar":  nil,
// 				"fizz": nil,
// 				"buzz": nil,
// 			},
// 			expected: "foo, bar, fizz, buzz",
// 		},
// 	}
//
// 	split := func(s string) []string {
// 		return strings.Split(s, ", ")
// 	}
//
// 	for n, c := range cc {
// 		t.Run(n, func(t *testing.T) {
// 			actual := (&JApiCore{}).getPropertiesNames(c.given)
// 			assert.ElementsMatch(t, split(c.expected), split(actual))
// 		})
// 	}
// }

// TODO move to catalog
// func TestJApiCore_processSchemaContentJSightAllOf(t *testing.T) {
// 	t.Run("positive", func(t *testing.T) {
// 		cc := map[string]struct {
// 			core               *JApiCore
// 			given              *catalog.ExchangeContent
// 			expectedUUT        *catalog.StringSet
// 			expectedProperties []*catalog.ExchangeContent
// 		}{
// 			"not object": {
// 				NewJApiCore(fs.NewFile("", `{}`)),
// 				&catalog.ExchangeContent{
// 					TokenType: "string",
// 				},
// 				catalog.NewStringSet(),
// 				nil,
// 			},
//
// 			"without allOf": {
// 				NewJApiCore(fs.NewFile("", `{}`)),
// 				&catalog.ExchangeContent{
// 					TokenType: jschema.TokenTypeObject,
// 					Children:  []*catalog.ExchangeContent{},
// 				},
// 				catalog.NewStringSet(),
// 				[]*catalog.ExchangeContent{},
// 			},
//
// 			"with allOf, single": {
// 				func() *JApiCore {
// 					core := NewJApiCore(fs.NewFile("", `{}`))
// 					core.catalog = catalog.NewCatalog()
// 					core.catalog.UserTypes.Set("@foo", &catalog.UserType{
// 						Schema: &catalog.ExchangeJSightSchema{
// 							ExchangeContent: &catalog.ExchangeContent{
// 								TokenType: jschema.TokenTypeObject,
// 								Children: []*catalog.ExchangeContent{
// 									{Key: catalog.SrtPtr("foo")},
// 									{Key: catalog.SrtPtr("bar")},
// 								},
// 							},
// 						},
// 					})
// 					return core
// 				}(),
// 				&catalog.ExchangeContent{
// 					TokenType: jschema.TokenTypeObject,
// 					Children:  []*catalog.ExchangeContent{},
// 					Rules: catalog.NewRules(
// 						[]catalog.Rule{
// 							{
// 								Key:         "allOf",
// 								TokenType:   catalog.RuleTokenTypeReference,
// 								ScalarValue: "@foo",
// 							},
// 						},
// 					),
// 				},
// 				catalog.NewStringSet("@foo"),
// 				[]*catalog.ExchangeContent{
// 					{
// 						Key:           catalog.SrtPtr("foo"),
// 						InheritedFrom: "@foo",
// 					},
// 					{
// 						Key:           catalog.SrtPtr("bar"),
// 						InheritedFrom: "@foo",
// 					},
// 				},
// 			},
//
// 			"with allOf, array": {
// 				func() *JApiCore {
// 					core := NewJApiCore(fs.NewFile("", `{}`))
// 					core.catalog = catalog.NewCatalog()
// 					core.catalog.UserTypes.Set("@foo", &catalog.UserType{
// 						Schema: &catalog.ExchangeJSightSchema{
// 							ExchangeContent: &catalog.ExchangeContent{
// 								TokenType: jschema.TokenTypeObject,
// 								Children: []*catalog.ExchangeContent{
// 									{Key: catalog.SrtPtr("foo1")},
// 									{Key: catalog.SrtPtr("foo2")},
// 								},
// 							},
// 						},
// 					})
// 					core.catalog.UserTypes.Set("@bar", &catalog.UserType{
// 						Schema: &catalog.ExchangeJSightSchema{
// 							ExchangeContent: &catalog.ExchangeContent{
// 								TokenType: jschema.TokenTypeObject,
// 								Children: []*catalog.ExchangeContent{
// 									{Key: catalog.SrtPtr("bar1")},
// 									{Key: catalog.SrtPtr("bar2")},
// 								},
// 							},
// 						},
// 					})
// 					return core
// 				}(),
// 				&catalog.ExchangeContent{
// 					TokenType: jschema.TokenTypeObject,
// 					Children:  []*catalog.ExchangeContent{},
// 					Rules: catalog.NewRules(
// 						[]catalog.Rule{
// 							{
// 								Key:       "allOf",
// 								TokenType: catalog.RuleTokenTypeArray,
// 								Children: []catalog.Rule{
// 									{
// 										TokenType:   catalog.RuleTokenTypeString,
// 										ScalarValue: "@foo",
// 									},
// 									{
// 										TokenType:   catalog.RuleTokenTypeString,
// 										ScalarValue: "@bar",
// 									},
// 								},
// 							},
// 						},
// 					),
// 				},
// 				catalog.NewStringSet("@bar", "@foo"),
// 				[]*catalog.ExchangeContent{
// 					{
// 						Key:           catalog.SrtPtr("foo1"),
// 						InheritedFrom: "@foo",
// 					},
// 					{
// 						Key:           catalog.SrtPtr("foo2"),
// 						InheritedFrom: "@foo",
// 					},
// 					{
// 						Key:           catalog.SrtPtr("bar1"),
// 						InheritedFrom: "@bar",
// 					},
// 					{
// 						Key:           catalog.SrtPtr("bar2"),
// 						InheritedFrom: "@bar",
// 					},
// 				},
// 			},
// 		}
//
// 		for n, c := range cc {
// 			t.Run(n, func(t *testing.T) {
// 				uut := catalog.NewStringSet()
//
// 				err := c.core.processExchangeContentAllOf(c.given, uut)
// 				require.NoError(t, err)
//
// 				assert.Equal(t, c.expectedUUT, uut)
// 				assert.Equal(t, c.expectedProperties, c.given.Children)
// 			})
// 		}
// 	})
//
// 	t.Run("negative", func(t *testing.T) {
// 		assert.Panics(t, func() {
// 			var core *JApiCore
//
// 			_ = core.processExchangeContentAllOf(nil, nil)
// 		})
// 	})
// }

// TODO move to catalog
// func TestJApiCore_inheritPropertiesFromUserType(t *testing.T) {
// 	t.Run("positive", func(t *testing.T) {
// 		cc := map[string]struct {
// 			given              *catalog.ExchangeContent
// 			expectedUUT        *catalog.StringSet
// 			expectedProperties []*catalog.ExchangeContent
// 		}{
// 			"sc.properties is nil": {
// 				&catalog.ExchangeContent{},
// 				catalog.NewStringSet("foo"),
// 				[]*catalog.ExchangeContent{
// 					{
// 						Key:           catalog.SrtPtr("bar"),
// 						InheritedFrom: "foo",
// 					},
// 				},
// 			},
//
// 			"sc.properties isn't nil": {
// 				&catalog.ExchangeContent{
// 					Children: []*catalog.ExchangeContent{
// 						{Key: catalog.SrtPtr("fizz")},
// 						{Key: catalog.SrtPtr("buzz")},
// 					},
// 				},
// 				catalog.NewStringSet("foo"),
// 				[]*catalog.ExchangeContent{
// 					{
// 						Key:           catalog.SrtPtr("bar"),
// 						InheritedFrom: "foo",
// 					},
// 					{
// 						Key: catalog.SrtPtr("fizz"),
// 					},
// 					{
// 						Key: catalog.SrtPtr("buzz"),
// 					},
// 				},
// 			},
// 		}
//
// 		for n, c := range cc {
// 			t.Run(n, func(t *testing.T) {
// 				core := NewJApiCore(fs.NewFile("", `{}`))
// 				core.catalog = catalog.NewCatalog()
// 				core.catalog.UserTypes.Set("foo", &catalog.UserType{
// 					Schema: &catalog.ExchangeJSightSchema{
// 						ExchangeContent: &catalog.ExchangeContent{
// 							TokenType: jschema.TokenTypeObject,
// 							Children: []*catalog.ExchangeContent{
// 								{Key: catalog.SrtPtr("bar")},
// 							},
// 						},
// 					},
// 				})
// 				uut := &catalog.StringSet{}
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
// 				catalog: catalog.NewCatalog(),
// 			},
//
// 			`the user type "foo" is not an object`: {
// 				catalog: func() *catalog.Catalog {
// 					c := catalog.NewCatalog()
// 					c.UserTypes.Set("foo", &catalog.UserType{
// 						Schema: &catalog.ExchangeJSightSchema{
// 							ExchangeContent: &catalog.ExchangeContent{
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
// 			core.catalog = catalog.NewCatalog()
// 			core.catalog.UserTypes.Set("foo", &catalog.UserType{
// 				Schema: &catalog.ExchangeJSightSchema{
// 					ExchangeContent: &catalog.ExchangeContent{
// 						TokenType: jschema.TokenTypeObject,
// 						Children: []*catalog.ExchangeContent{
// 							{Key: catalog.SrtPtr("bar")},
// 						},
// 					},
// 				},
// 			})
// 			err := core.inheritPropertiesFromUserType(
// 				&catalog.ExchangeContent{
// 					Children: []*catalog.ExchangeContent{
// 						{Key: catalog.SrtPtr("bar")},
// 					},
// 				},
// 				&catalog.StringSet{},
// 				"foo",
// 			)
// 			assert.EqualError(t, err, `it is not allowed to override the "bar" property from the user type "foo"`)
// 		})
//
// 		t.Run("core is nil", func(t *testing.T) {
// 			assert.Panics(t, func() {
// 				var core *JApiCore
// 				_ = core.inheritPropertiesFromUserType(
// 					&catalog.ExchangeContent{},
// 					&catalog.StringSet{},
// 					"foo",
// 				)
// 			})
// 		})
//
// 		t.Run("sc is nil", func(t *testing.T) {
// 			assert.Panics(t, func() {
// 				core := &JApiCore{
// 					catalog: catalog.NewCatalog(),
// 				}
// 				core.catalog.UserTypes.Set("foo", &catalog.UserType{
// 					Schema: &catalog.ExchangeJSightSchema{
// 						ExchangeContent: &catalog.ExchangeContent{
// 							TokenType: jschema.TokenTypeObject,
// 						},
// 					},
// 				})
// 				_ = core.inheritPropertiesFromUserType(
// 					nil,
// 					&catalog.StringSet{},
// 					"foo",
// 				)
// 			})
// 		})
//
// 		t.Run("uut is nil", func(t *testing.T) {
// 			assert.Panics(t, func() {
// 				core := &JApiCore{
// 					catalog: catalog.NewCatalog(),
// 				}
// 				core.catalog.UserTypes.Set("foo", &catalog.UserType{
// 					Schema: &catalog.ExchangeJSightSchema{
// 						ExchangeContent: &catalog.ExchangeContent{
// 							TokenType: jschema.TokenTypeObject,
// 							Children: []*catalog.ExchangeContent{
// 								{Key: catalog.SrtPtr("foo")},
// 							},
// 						},
// 					},
// 				})
// 				_ = core.inheritPropertiesFromUserType(
// 					&catalog.ExchangeContent{},
// 					nil,
// 					"foo",
// 				)
// 			})
// 		})
// 	})
// }
