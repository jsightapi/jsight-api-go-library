package jerr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsightapi/jsight-schema-go-library/fs"
)

func TestNewJApiError(t *testing.T) {
	const (
		msg   = "fake error"
		index = 1
	)
	file := fs.NewFile("foo", "123456")

	je := NewJApiError(msg, file, index)

	assert.Equal(t, msg, je.Msg)
	assert.Equal(t, NewLocation(file, index), je.Location)
	assert.Nil(t, je.wrapped)
}

func TestJApiError_OccurredInFile(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			je                 *JApiError
			expectedStackTrace []stackTraceItem
		}{
			"empty stack trace": {
				&JApiError{},
				[]stackTraceItem{
					{
						path:   "/foo/bar",
						atLine: 2,
					},
				},
			},

			"not empty stack trace": {
				&JApiError{
					includeTrace: []stackTraceItem{
						{
							path:   "/fizz",
							atLine: 100,
						},
						{
							path:   "/buzz",
							atLine: 500,
						},
					},
				},
				[]stackTraceItem{
					{
						path:   "/fizz",
						atLine: 100,
					},
					{
						path:   "/buzz",
						atLine: 500,
					},
					{
						path:   "/foo/bar",
						atLine: 2,
					},
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				c.je.OccurredInFile(fs.NewFile("/foo/bar", "12\n34\n56"), 3)

				assert.Equal(t, c.expectedStackTrace, c.je.includeTrace)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("nil error", func(t *testing.T) {
			assert.Panics(t, func() {
				var je *JApiError
				je.OccurredInFile(nil, 0)
			})
		})

		t.Run("nil file", func(t *testing.T) {
			assert.Panics(t, func() {
				(&JApiError{}).OccurredInFile(nil, 0)
			})
		})
	})
}

func TestJApiError_HasStackTrace(t *testing.T) {
	cc := map[string]struct {
		given    *JApiError
		expected bool
	}{
		"nil error": {},
		"nil stack trace": {
			given: &JApiError{},
		},
		"empty stack trace": {
			given: &JApiError{
				includeTrace: make([]stackTraceItem, 0, 10),
			},
		},
		"not empty stack trace": {
			given: &JApiError{
				includeTrace: []stackTraceItem{
					{},
				},
			},
			expected: true,
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			actual := c.given.HasStackTrace()
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestJApiError_Error(t *testing.T) {
	t.Run("without stack trace", func(t *testing.T) {
		const msg = "fake error"

		je := NewJApiError(msg, fs.NewFile("foo", "123456"), 1)

		assert.Equal(t, msg, je.Error())
	})

	t.Run("with stack trace", func(t *testing.T) {
		je := NewJApiError("fake error", fs.NewFile("foo", "123456"), 1)

		je.OccurredInFile(fs.NewFile("bar", "1\n2\n3\n4\n5\n6"), 4)
		je.OccurredInFile(fs.NewFile("fizz", "123456"), 4)
		je.OccurredInFile(fs.NewFile("buzz", "123\n456"), 4)

		assert.Equal(t, `fake error
foo:1
bar:3
fizz:1
buzz:2`, je.Error())
	})
}
