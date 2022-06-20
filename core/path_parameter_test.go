package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_splitPath(t *testing.T) {
	cc := map[string][]string{
		"":             {},
		"/":            {},
		"////":         {},
		"aaa/bbb":      {"aaa", "bbb"},
		"/aaa/bbb":     {"aaa", "bbb"},
		"/aaa/bbb/":    {"aaa", "bbb"},
		"//aaa//bbb//": {"aaa", "bbb"},
	}

	for given, expected := range cc {
		t.Run(given, func(t *testing.T) {
			actual := splitPath(given)
			assert.Equal(t, expected, actual)
		})
	}
}

func Test_PathParameters(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string][]PathParameter{
			"":     {},
			"/":    {},
			"/aaa": {},
			"/aaa/{id}": {
				{"aaa/{id}", "id"},
			},
			"/aaa/{id}/bbb/{some}": {
				{"aaa/{id}", "id"},
				{"aaa/{id}/bbb/{some}", "some"},
			},
			"///aaa///{id}///": {
				{"aaa/{id}", "id"},
			},
		}

		for given, expected := range cc {
			t.Run(given, func(t *testing.T) {
				acual, err := PathParameters(given)
				require.NoError(t, err)
				assert.Equal(t, expected, acual)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		ss := []string{
			"/aaa/{id}/bbb/{id}",
			"/aaa/{}",
		}

		for _, s := range ss {
			t.Run(s, func(t *testing.T) {
				_, err := PathParameters(s)
				assert.Error(t, err)
			})
		}
	})
}

func Test_duplicatedPathParameters(t *testing.T) {
	cc := map[string]struct {
		args []PathParameter
		want string
	}{
		"empty": {
			[]PathParameter{},
			"",
		},
		"one": {
			[]PathParameter{
				{"any string", "id"},
			},
			"",
		},
		"two": {
			[]PathParameter{
				{"any string", "id"},
				{"any string", "name"},
			},
			"",
		},
		"three": {
			[]PathParameter{
				{"any string", "id"},
				{"any string", "name"},
				{"any string", "id"},
			},
			"id",
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			got := duplicatedPathParameters(c.args)
			assert.Equal(t, c.want, got)
		})
	}
}
