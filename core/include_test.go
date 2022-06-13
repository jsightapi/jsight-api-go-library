package core

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/assert"

	"github.com/jsightapi/jsight-api-go-library/scanner"
)

func Test_isIncludeKeyword(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given    *scanner.Lexeme
			expected bool
		}{
			"nil": {},
			"not a keyword": {
				given: scanner.NewLexeme(scanner.Schema, 0, 0, nil),
			},
			"not an include keyword": {
				given: scanner.NewLexeme(scanner.Keyword, 0, 1, fs.NewFile("", []byte("12"))),
			},
			"an include keyword": {
				given:    scanner.NewLexeme(scanner.Keyword, 0, 6, fs.NewFile("", []byte("INCLUDE"))),
				expected: true,
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				actual := isIncludeKeyword(c.given)
				assert.Equal(t, c.expected, actual)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			isIncludeKeyword(&scanner.Lexeme{})
		})
	})
}

func Test_validateIncludeFileName(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ss := []string{
			"foo",
			".foo",
			"foo/bar",
			".foo/bar",
			"..foo/bar",
			"foo/bar.jst",
		}

		for _, s := range ss {
			t.Run(s, func(t *testing.T) {
				err := validateIncludeFileName(s)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]string{
			"/foo/bar":   "mustn't starts with '/'",
			"/./bar":     "mustn't starts with '/'",
			"/../bar":    "mustn't starts with '/'",
			"./foo/bar":  "mustn't include '..' or '.'",
			"../foo/bar": "mustn't include '..' or '.'",
			"foo/.":      "mustn't include '..' or '.'",
			"foo/./":     "mustn't include '..' or '.'",
			"foo/./bar":  "mustn't include '..' or '.'",
			"foo/..":     "mustn't include '..' or '.'",
			"foo/../":    "mustn't include '..' or '.'",
			"foo/../bar": "mustn't include '..' or '.'",
			"foo/./bar/../fizz/../buzz/./fizzbuzz/..": "mustn't include '..' or '.'",
			"foo\\bar":           "the separator for directories and files should be the symbol '/'",
			"foo/bar\\fizz/buzz": "the separator for directories and files should be the symbol '/'",
		}

		for given, expected := range cc {
			t.Run(given, func(t *testing.T) {
				err := validateIncludeFileName(given)
				assert.EqualError(t, err, expected)
			})
		}
	})
}
