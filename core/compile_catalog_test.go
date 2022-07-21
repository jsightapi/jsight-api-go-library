package core

import (
	"testing"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
)

func TestJApiCore_BuildResourceMethodsPathVariables(t *testing.T) {
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
				err := c.given.BuildResourceMethodsPathVariables()
				require.Nil(t, err)

				assert.Equal(t, c.expected, c.given)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]*JApiCore{
			`Has unused parameters "bar" in schema`: {
				rawPathVariables: []rawPathVariable{
					{
						schema: catalog.Schema{
							ContentJSight: &catalog.SchemaContentJSight{
								Children: []*catalog.SchemaContentJSight{
									{Key: catalog.SrtPtr("foo")},
									{Key: catalog.SrtPtr("bar")},
								},
							},
						},
						parameters: []PathParameter{
							{
								path:      "/path/{foo}/{bar}",
								parameter: "foo",
							},
						},
						pathDirective: *directive.New(
							directive.Get,
							directive.NewCoords(
								fs.NewFile("foo", bytes.Bytes("/path/{foo}/{bar}")),
								0,
								5,
							),
						),
					},
				},

				catalog: &catalog.Catalog{
					Interactions: &catalog.Interactions{},
				},
			},
		}

		for expected, core := range cc {
			t.Run(expected, func(t *testing.T) {
				err := core.BuildResourceMethodsPathVariables()
				assert.EqualError(t, err, expected)
			})
		}

		t.Run("nil", func(t *testing.T) {
			assert.Panics(t, func() {
				var core *JApiCore

				_ = core.BuildResourceMethodsPathVariables()
			})
		})
	})
}

// TODO func TestCore_propertiesToMap(t *testing.T) {
// 	cc := map[string]struct {
// 		given    *catalog.Properties
// 		expected map[string]*catalog.SchemaContentJSight
// 	}{
// 		"nil": {},
//
// 		"empty": {
// 			given: &catalog.Properties{},
// 		},
//
// 		"with data": {
// 			given: catalog.NewProperties(
// 				map[string]*catalog.SchemaContentJSight{
// 					"foo": nil,
// 					"bar": {
// 						Note: "fake",
// 					},
// 				},
// 				[]string{"foo", "bar"},
// 			),
// 			expected: map[string]*catalog.SchemaContentJSight{
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
// 		given    map[string]*catalog.SchemaContentJSight
// 		expected string
// 	}{
// 		"nil": {},
//
// 		"empty": {
// 			given: map[string]*catalog.SchemaContentJSight{},
// 		},
//
// 		"with data": {
// 			given: map[string]*catalog.SchemaContentJSight{
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

func TestJApiCore_processSchemaContentJSightAllOf(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			core               *JApiCore
			given              *catalog.SchemaContentJSight
			expectedUUT        *catalog.StringSet
			expectedProperties []*catalog.SchemaContentJSight
		}{
			"not object": {
				NewJApiCore(fs.NewFile("", []byte(`{}`))),
				&catalog.SchemaContentJSight{
					TokenType: "string",
				},
				catalog.NewStringSet(),
				nil,
			},

			"without allOf": {
				NewJApiCore(fs.NewFile("", []byte(`{}`))),
				&catalog.SchemaContentJSight{
					TokenType: jschema.JSONTypeObject,
					Children:  []*catalog.SchemaContentJSight{},
				},
				catalog.NewStringSet(),
				[]*catalog.SchemaContentJSight{},
			},

			"with allOf, single": {
				&JApiCore{
					catalog: func() *catalog.Catalog {
						c := catalog.NewCatalog()
						c.UserTypes.Set("@foo", &catalog.UserType{
							Schema: catalog.Schema{
								ContentJSight: &catalog.SchemaContentJSight{
									TokenType: jschema.JSONTypeObject,
									Children: []*catalog.SchemaContentJSight{
										{Key: catalog.SrtPtr("foo")},
										{Key: catalog.SrtPtr("bar")},
									},
								},
							},
						})
						return c
					}(),
				},
				&catalog.SchemaContentJSight{
					TokenType: jschema.JSONTypeObject,
					Children:  []*catalog.SchemaContentJSight{},
					Rules: catalog.NewRules(
						[]catalog.Rule{
							{
								Key:         "allOf",
								TokenType:   catalog.RuleTokenTypeString,
								ScalarValue: "@foo",
							},
						},
					),
				},
				catalog.NewStringSet("@foo"),
				[]*catalog.SchemaContentJSight{
					{
						Key:           catalog.SrtPtr("foo"),
						InheritedFrom: "@foo",
					},
					{
						Key:           catalog.SrtPtr("bar"),
						InheritedFrom: "@foo",
					},
				},
			},

			"with allOf, array": {
				&JApiCore{
					catalog: func() *catalog.Catalog {
						c := catalog.NewCatalog()
						c.UserTypes.Set("@foo", &catalog.UserType{
							Schema: catalog.Schema{
								ContentJSight: &catalog.SchemaContentJSight{
									TokenType: jschema.JSONTypeObject,
									Children: []*catalog.SchemaContentJSight{
										{Key: catalog.SrtPtr("foo1")},
										{Key: catalog.SrtPtr("foo2")},
									},
								},
							},
						})
						c.UserTypes.Set("@bar", &catalog.UserType{
							Schema: catalog.Schema{
								ContentJSight: &catalog.SchemaContentJSight{
									TokenType: jschema.JSONTypeObject,
									Children: []*catalog.SchemaContentJSight{
										{Key: catalog.SrtPtr("bar1")},
										{Key: catalog.SrtPtr("bar2")},
									},
								},
							},
						})
						return c
					}(),
				},
				&catalog.SchemaContentJSight{
					TokenType: jschema.JSONTypeObject,
					Children:  []*catalog.SchemaContentJSight{},
					Rules: catalog.NewRules(
						[]catalog.Rule{
							{
								Key:       "allOf",
								TokenType: catalog.RuleTokenTypeArray,
								Children: []catalog.Rule{
									{
										TokenType:   catalog.RuleTokenTypeString,
										ScalarValue: "@foo",
									},
									{
										TokenType:   catalog.RuleTokenTypeString,
										ScalarValue: "@bar",
									},
								},
							},
						},
					),
				},
				catalog.NewStringSet("@bar", "@foo"),
				[]*catalog.SchemaContentJSight{
					{
						Key:           catalog.SrtPtr("foo1"),
						InheritedFrom: "@foo",
					},
					{
						Key:           catalog.SrtPtr("foo2"),
						InheritedFrom: "@foo",
					},
					{
						Key:           catalog.SrtPtr("bar1"),
						InheritedFrom: "@bar",
					},
					{
						Key:           catalog.SrtPtr("bar2"),
						InheritedFrom: "@bar",
					},
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				uut := catalog.NewStringSet()

				err := c.core.processSchemaContentJSightAllOf(c.given, uut)
				require.NoError(t, err)

				assert.Equal(t, c.expectedUUT, uut)
				assert.Equal(t, c.expectedProperties, c.given.Children)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			var core *JApiCore

			core.processSchemaContentJSightAllOf(nil, nil)
		})
	})
}

func TestJApiCore_inheritPropertiesFromUserType(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given              *catalog.SchemaContentJSight
			expectedUUT        *catalog.StringSet
			expectedProperties []*catalog.SchemaContentJSight
		}{
			"sc.properties is nil": {
				&catalog.SchemaContentJSight{},
				catalog.NewStringSet("foo"),
				[]*catalog.SchemaContentJSight{
					{
						Key:           catalog.SrtPtr("bar"),
						InheritedFrom: "foo",
					},
				},
			},

			"sc.properties isn't nil": {
				&catalog.SchemaContentJSight{
					Children: []*catalog.SchemaContentJSight{
						{Key: catalog.SrtPtr("fizz")},
						{Key: catalog.SrtPtr("buzz")},
					},
				},
				catalog.NewStringSet("foo"),
				[]*catalog.SchemaContentJSight{
					{
						Key:           catalog.SrtPtr("bar"),
						InheritedFrom: "foo",
					},
					{
						Key: catalog.SrtPtr("fizz"),
					},
					{
						Key: catalog.SrtPtr("buzz"),
					},
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				core := &JApiCore{
					catalog: catalog.NewCatalog(),
				}
				core.catalog.UserTypes.Set("foo", &catalog.UserType{
					Schema: catalog.Schema{
						ContentJSight: &catalog.SchemaContentJSight{
							TokenType: jschema.JSONTypeObject,
							Children: []*catalog.SchemaContentJSight{
								{Key: catalog.SrtPtr("bar")},
							},
						},
					},
				})
				uut := &catalog.StringSet{}

				err := core.inheritPropertiesFromUserType(c.given, uut, "foo")
				require.NoError(t, err)

				assert.Equal(t, c.expectedUUT, uut)
				assert.Equal(t, c.expectedProperties, c.given.Children)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]*JApiCore{
			`the user type "foo" not found`: {
				catalog: catalog.NewCatalog(),
			},

			`the user type "foo" is not an object`: {
				catalog: func() *catalog.Catalog {
					c := catalog.NewCatalog()
					c.UserTypes.Set("foo", &catalog.UserType{
						Schema: catalog.Schema{
							ContentJSight: &catalog.SchemaContentJSight{
								TokenType: jschema.JSONTypeString,
							},
						},
					})
					return c
				}(),
			},
		}

		for expected, core := range cc {
			t.Run(expected, func(t *testing.T) {
				err := core.inheritPropertiesFromUserType(nil, nil, "foo")
				assert.EqualError(t, err, expected)
			})
		}

		t.Run("property already set", func(t *testing.T) {
			core := &JApiCore{
				catalog: catalog.NewCatalog(),
			}
			core.catalog.UserTypes.Set("foo", &catalog.UserType{
				Schema: catalog.Schema{
					ContentJSight: &catalog.SchemaContentJSight{
						TokenType: jschema.JSONTypeObject,
						Children: []*catalog.SchemaContentJSight{
							{Key: catalog.SrtPtr("bar")},
						},
					},
				},
			})
			err := core.inheritPropertiesFromUserType(
				&catalog.SchemaContentJSight{
					Children: []*catalog.SchemaContentJSight{
						{Key: catalog.SrtPtr("bar")},
					},
				},
				&catalog.StringSet{},
				"foo",
			)
			assert.EqualError(t, err, `it is not allowed to override the "bar" property from the user type "foo"`) //nolint:lll
		})

		t.Run("core is nil", func(t *testing.T) {
			assert.Panics(t, func() {
				var core *JApiCore
				core.inheritPropertiesFromUserType(
					&catalog.SchemaContentJSight{},
					&catalog.StringSet{},
					"foo",
				)
			})
		})

		t.Run("sc is nil", func(t *testing.T) {
			assert.Panics(t, func() {
				core := &JApiCore{
					catalog: catalog.NewCatalog(),
				}
				core.catalog.UserTypes.Set("foo", &catalog.UserType{
					Schema: catalog.Schema{
						ContentJSight: &catalog.SchemaContentJSight{
							TokenType: jschema.JSONTypeObject,
						},
					},
				})
				core.inheritPropertiesFromUserType(
					nil,
					&catalog.StringSet{},
					"foo",
				)
			})
		})

		t.Run("uut is nil", func(t *testing.T) {
			assert.Panics(t, func() {
				core := &JApiCore{
					catalog: catalog.NewCatalog(),
				}
				core.catalog.UserTypes.Set("foo", &catalog.UserType{
					Schema: catalog.Schema{
						ContentJSight: &catalog.SchemaContentJSight{
							TokenType: jschema.JSONTypeObject,
							Children: []*catalog.SchemaContentJSight{
								{Key: catalog.SrtPtr("foo")},
							},
						},
					},
				})
				core.inheritPropertiesFromUserType(
					&catalog.SchemaContentJSight{},
					nil,
					"foo",
				)
			})
		})
	})
}
