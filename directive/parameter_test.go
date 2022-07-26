package directive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

func Test_unescapeParameter(t *testing.T) {
	tests := map[string]string{
		"123":              "123",
		"abc":              "abc",
		`aa \" bb \" cc`:   `aa \" bb \" cc`,
		`""`:               "",
		`"abc"`:            "abc",
		`"aa \" bb \" cc"`: `aa " bb " cc`,
		`"aa \\ bb \\ cc"`: `aa \ bb \ cc`,
	}

	for given, expected := range tests {
		t.Run(given, func(t *testing.T) {
			actual := unescapeParameter(bytes.Bytes(given))
			assert.Equal(t, expected, string(actual))
		})
	}
}

func Benchmark_unescapeParameter(b *testing.B) {
	str := bytes.Bytes(`"aa \" bb \" \" bb \" \" bb \" \" bb \" \" bb \" cc"`)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		unescapeParameter(str)
	}
}

func TestIsArrayOfTypes(t *testing.T) {
	tests := map[string]bool{
		"[@a]":  true,
		"[@世界]": false,
		"[@@]":  false,
		"[@]":   false,
		"[]":    false,
		"[":     false,
		"]":     false,
		"[@aaa": false,
		"@aaa]": false,
	}

	for given, expected := range tests {
		t.Run(given, func(t *testing.T) {
			actual := IsArrayOfTypes(bytes.Bytes(given))
			assert.Equal(t, expected, actual)
		})
	}
}

func TestDirective_AppendParameter(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		type testCase struct {
			directive          *Directive
			given              string
			expectedParameters map[string]string
		}
		cc := map[string]testCase{
			// Type
			fmt.Sprintf("%s schema notation", Type): {
				New(Type, Coords{}),
				"jsight",
				map[string]string{
					"SchemaNotation": "jsight",
				},
			},
			fmt.Sprintf("%s array of types", Type): {
				New(Type, Coords{}),
				"[@foo]",
				map[string]string{
					"Name": "[@foo]",
				},
			},
			fmt.Sprintf("%s type", Type): {
				New(Type, Coords{}),
				"@foo",
				map[string]string{
					"Name": "@foo",
				},
			},

			// Query
			fmt.Sprintf("%s htmlFormEncoded", Query): {
				New(Query, Coords{}),
				"htmlFormEncoded",
				map[string]string{
					"Format": "htmlFormEncoded",
				},
			},
			fmt.Sprintf("%s noFormat", Query): {
				New(Query, Coords{}),
				"noFormat",
				map[string]string{
					"Format": "noFormat",
				},
			},
			Query.String(): {
				New(Query, Coords{}),
				"foo",
				map[string]string{
					"QueryExample": "foo",
				},
			},

			Jsight.String(): {
				New(Jsight, Coords{}),
				"foo",
				map[string]string{
					"Version": "foo",
				},
			},

			Title.String(): {
				New(Title, Coords{}),
				"foo",
				map[string]string{
					"Title": "foo",
				},
			},

			Version.String(): {
				New(Version, Coords{}),
				"foo",
				map[string]string{
					"Version": "foo",
				},
			},

			BaseUrl.String(): {
				New(BaseUrl, Coords{}),
				"foo",
				map[string]string{
					"Path": "foo",
				},
			},
		}

		for _, typ := range []Enumeration{Url, Get, Post, Put, Patch, Delete} {
			cc[typ.String()] = testCase{
				New(typ, Coords{}),
				"foo",
				map[string]string{
					"Path": "foo",
				},
			}
		}

		for _, typ := range []Enumeration{Request, HTTPResponseCode, Body} {
			cc[fmt.Sprintf("%s schema notation", typ)] = testCase{
				New(typ, Coords{}),
				"jsight",
				map[string]string{
					"SchemaNotation": "jsight",
				},
			}

			cc[fmt.Sprintf("%s array of types", typ)] = testCase{
				New(typ, Coords{}),
				"[@foo]",
				map[string]string{
					"Type": "[@foo]",
				},
			}

			cc[fmt.Sprintf("%s type", typ)] = testCase{
				New(typ, Coords{}),
				"@foo",
				map[string]string{
					"Type": "@foo",
				},
			}
		}

		for _, typ := range []Enumeration{Server, Enum, Macro, Paste} {
			cc[typ.String()] = testCase{
				New(typ, Coords{}),
				"@foo",
				map[string]string{
					"Name": "@foo",
				},
			}
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				err := c.directive.AppendParameter(bytes.Bytes(c.given))
				require.NoError(t, err)

				assert.Equal(t, c.expectedParameters, c.directive.namedParameters)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := []Enumeration{
			Type,
			Request,
			HTTPResponseCode,
			Body,
			Server,
			Enum,
			Macro,
			Paste,
			Path,
		}

		for _, typ := range cc {
			t.Run(typ.String(), func(t *testing.T) {
				err := New(typ, Coords{}).AppendParameter(bytes.Bytes("foo"))
				assert.EqualError(t, err, `incorrect parameter "foo"`)
			})
		}
	})
}
