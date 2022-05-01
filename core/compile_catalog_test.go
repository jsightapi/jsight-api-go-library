package core

import (
	"strings"
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
						ResourceMethods: &catalog.ResourceMethods{},
					},
				},
				&JApiCore{
					catalog: &catalog.Catalog{
						ResourceMethods: &catalog.ResourceMethods{},
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
								Properties: catalog.NewProperties(
									map[string]*catalog.SchemaContentJSight{
										"foo": nil,
										"bar": nil,
									},
									[]string{"foo", "bar"},
								),
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
					ResourceMethods: &catalog.ResourceMethods{},
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

func TestCore_propertiesToMap(t *testing.T) {
	cc := map[string]struct {
		given    *catalog.Properties
		expected map[string]*catalog.SchemaContentJSight
	}{
		"nil": {},

		"empty": {
			given: &catalog.Properties{},
		},

		"with data": {
			given: catalog.NewProperties(
				map[string]*catalog.SchemaContentJSight{
					"foo": nil,
					"bar": {
						Note: "fake",
					},
				},
				[]string{"foo", "bar"},
			),
			expected: map[string]*catalog.SchemaContentJSight{
				"foo": nil,
				"bar": {
					Note: "fake",
				},
			},
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			actual := (&JApiCore{}).propertiesToMap(c.given)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestCore_getPropertiesNames(t *testing.T) {
	cc := map[string]struct {
		given    map[string]*catalog.SchemaContentJSight
		expected string
	}{
		"nil": {},

		"empty": {
			given: map[string]*catalog.SchemaContentJSight{},
		},

		"with data": {
			given: map[string]*catalog.SchemaContentJSight{
				"foo":  nil,
				"bar":  nil,
				"fizz": nil,
				"buzz": nil,
			},
			expected: "foo, bar, fizz, buzz",
		},
	}

	split := func(s string) []string {
		return strings.Split(s, ", ")
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			actual := (&JApiCore{}).getPropertiesNames(c.given)
			assert.ElementsMatch(t, split(c.expected), split(actual))
		})
	}
}

func TestJApiCore_processSchemaContentJSightAllOf(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			core               *JApiCore
			given              *catalog.SchemaContentJSight
			expectedUUT        *catalog.StringSet
			expectedProperties *catalog.Properties
		}{
			"not object": {
				NewJApiCore(fs.NewFile("", []byte(`{}`))),
				&catalog.SchemaContentJSight{
					JsonType: "string",
				},
				catalog.NewStringSet(),
				nil,
			},

			"without allOf": {
				NewJApiCore(fs.NewFile("", []byte(`{}`))),
				&catalog.SchemaContentJSight{
					JsonType:   jschema.JSONTypeObject,
					Properties: &catalog.Properties{},
					Rules:      &catalog.Rules{},
				},
				catalog.NewStringSet(),
				catalog.NewProperties(nil, nil),
			},

			"with allOf, single": {
				&JApiCore{
					catalog: func() *catalog.Catalog {
						c := catalog.NewCatalog()
						c.UserTypes.Set("@foo", &catalog.UserType{
							Schema: catalog.Schema{
								ContentJSight: &catalog.SchemaContentJSight{
									JsonType: jschema.JSONTypeObject,
									Properties: catalog.NewProperties(
										map[string]*catalog.SchemaContentJSight{
											"foo": {},
											"bar": {},
										},
										[]string{"foo", "bar"},
									),
								},
							},
						})
						return c
					}(),
				},
				&catalog.SchemaContentJSight{
					JsonType:   jschema.JSONTypeObject,
					Properties: &catalog.Properties{},
					Rules: catalog.NewRules(
						map[string]catalog.Rule{
							"allOf": {
								JsonType:    jschema.JSONTypeString,
								ScalarValue: "@foo",
								Properties:  &catalog.Rules{},
							},
						},
						[]string{"allOf"},
					),
				},
				catalog.NewStringSet("@foo"),
				catalog.NewProperties(
					map[string]*catalog.SchemaContentJSight{
						"foo": {
							InheritedFrom: "@foo",
						},
						"bar": {
							InheritedFrom: "@foo",
						},
					},
					[]string{"foo", "bar"},
				),
			},

			"with allOf, array": {
				&JApiCore{
					catalog: func() *catalog.Catalog {
						c := catalog.NewCatalog()
						c.UserTypes.Set("@foo", &catalog.UserType{
							Schema: catalog.Schema{
								ContentJSight: &catalog.SchemaContentJSight{
									JsonType: jschema.JSONTypeObject,
									Properties: catalog.NewProperties(
										map[string]*catalog.SchemaContentJSight{
											"foo1": {},
											"foo2": {},
										},
										[]string{"foo1", "foo2"},
									),
								},
							},
						})
						c.UserTypes.Set("@bar", &catalog.UserType{
							Schema: catalog.Schema{
								ContentJSight: &catalog.SchemaContentJSight{
									JsonType: jschema.JSONTypeObject,
									Properties: catalog.NewProperties(
										map[string]*catalog.SchemaContentJSight{
											"bar1": {},
											"bar2": {},
										},
										[]string{"bar1", "bar2"},
									),
								},
							},
						})
						return c
					}(),
				},
				&catalog.SchemaContentJSight{
					JsonType:   jschema.JSONTypeObject,
					Properties: &catalog.Properties{},
					Rules: catalog.NewRules(
						map[string]catalog.Rule{
							"allOf": {
								JsonType:   jschema.JSONTypeArray,
								Properties: &catalog.Rules{},
								Items: []catalog.Rule{
									{
										JsonType:    jschema.JSONTypeString,
										ScalarValue: "@foo",
										Properties:  &catalog.Rules{},
									},
									{
										JsonType:    jschema.JSONTypeString,
										ScalarValue: "@bar",
										Properties:  &catalog.Rules{},
									},
								},
							},
						},
						[]string{"allOf"},
					),
				},
				catalog.NewStringSet("@bar", "@foo"),
				catalog.NewProperties(
					map[string]*catalog.SchemaContentJSight{
						"foo1": {
							InheritedFrom: "@foo",
						},
						"foo2": {
							InheritedFrom: "@foo",
						},
						"bar1": {
							InheritedFrom: "@bar",
						},
						"bar2": {
							InheritedFrom: "@bar",
						},
					},
					[]string{"foo1", "foo2", "bar1", "bar2"},
				),
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				uut := catalog.NewStringSet()

				err := c.core.processSchemaContentJSightAllOf(c.given, uut)
				require.NoError(t, err)

				assert.Equal(t, c.expectedUUT, uut)
				assert.Equal(t, c.expectedProperties, c.given.Properties)
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
			expectedProperties *catalog.Properties
		}{
			"sc.properties is nil": {
				&catalog.SchemaContentJSight{},
				catalog.NewStringSet("foo"),
				catalog.NewProperties(
					map[string]*catalog.SchemaContentJSight{
						"bar": {
							InheritedFrom: "foo",
						},
					},
					[]string{"bar"},
				),
			},

			"sc.properties isn't nil": {
				&catalog.SchemaContentJSight{
					Properties: catalog.NewProperties(
						map[string]*catalog.SchemaContentJSight{
							"fizz": {},
							"buzz": {},
						},
						[]string{"fizz", "buzz"},
					),
				},
				catalog.NewStringSet("foo"),
				catalog.NewProperties(
					map[string]*catalog.SchemaContentJSight{
						"bar": {
							InheritedFrom: "foo",
						},
						"fizz": {},
						"buzz": {},
					},
					[]string{"bar", "fizz", "buzz"},
				),
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
							JsonType: jschema.JSONTypeObject,
							Properties: catalog.NewProperties(
								map[string]*catalog.SchemaContentJSight{
									"bar": {},
								},
								[]string{"bar"},
							),
						},
					},
				})
				uut := &catalog.StringSet{}

				err := core.inheritPropertiesFromUserType(c.given, uut, "foo")
				require.NoError(t, err)

				assert.Equal(t, c.expectedUUT, uut)
				assert.Equal(t, c.expectedProperties, c.given.Properties)
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
								JsonType: jschema.JSONTypeString,
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
						JsonType: jschema.JSONTypeObject,
						Properties: catalog.NewProperties(
							map[string]*catalog.SchemaContentJSight{
								"bar": {},
							},
							[]string{"bar"},
						),
					},
				},
			})
			err := core.inheritPropertiesFromUserType(
				&catalog.SchemaContentJSight{
					Properties: catalog.NewProperties(
						map[string]*catalog.SchemaContentJSight{
							"bar": {},
						},
						[]string{"bar"},
					),
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
							JsonType: jschema.JSONTypeObject,
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
							JsonType: jschema.JSONTypeObject,
							Properties: catalog.NewProperties(
								map[string]*catalog.SchemaContentJSight{
									"foo": {},
								},
								[]string{"foo"},
							),
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
