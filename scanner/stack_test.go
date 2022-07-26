package scanner

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func TestIncludedFileStack_Push(t *testing.T) {
	newFile := func(path string) *fs.File {
		return fs.NewFile(path, "123456")
	}

	newStack := func(pp ...string) func(*testing.T) *Stack {
		return func(t *testing.T) *Stack {
			s := &Stack{}

			for _, p := range pp {
				require.NoError(t, s.Push(NewJApiScanner(newFile(p)), 0))
			}

			return s
		}
	}

	t.Run("positive", func(t *testing.T) {
		file := newFile("/foo/bar")

		cc := map[string]func(*testing.T) *Stack{
			"empty": newStack(),
			"with some files": newStack(
				"/fizz/buzz",
				"/foo/bar/0",
			),
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				err := c(t).Push(NewJApiScanner(file), 0)
				require.NoError(t, err)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			var s *Stack
			_ = s.Push(nil, 0)
		})

		cc := []func(*testing.T) *Stack{
			newStack("/foo/bar"),
			newStack(
				"/foo/bar1",
				"/foo/bar2",
				"/foo/bar",
				"/foo/bar3",
			),
			newStack(
				"/foo/bar",
				"/foo/bar1",
				"/foo/bar2",
				"/foo/bar3",
			),
		}

		for _, c := range cc {
			t.Run("", func(t *testing.T) {
				err := c(t).Push(NewJApiScanner(newFile("/foo/bar")), 0)
				assert.ErrorIs(t, err, ErrRecursionDetected)
			})
		}
	})
}

func TestStack_computeScannerHash(t *testing.T) {
	t.Run("stability", func(t *testing.T) {
		const times = 1000
		scanner := &Scanner{file: fs.NewFile("/foo/bar", "content")}

		origin, err := (&Stack{}).computeScannerHash(scanner)
		require.NoError(t, err)

		for i := 0; i < times; i++ {
			n, err := (&Stack{}).computeScannerHash(scanner)
			require.NoError(t, err)

			assert.Equal(t, origin, n)
		}
	})
}

func BenchmarkStack_computeScannerHash(b *testing.B) {
	scanner := &Scanner{file: fs.NewFile("/foo/bar", "content")}
	stack := &Stack{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := stack.computeScannerHash(scanner)
		require.NoError(b, err)
	}
}

func TestStack_AddIncludeTraceToError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		(&Stack{}).AddIncludeTraceToError(nil)
	})

	file := fs.NewFile("/foo/bar", "12\n34\n56")

	cc := map[string]struct {
		je       *jerr.JApiError
		stack    []stackItem
		expected string
	}{
		"error without trace, empty stack": {
			je:       jerr.NewJApiError("fake error", file, 4),
			expected: "fake error",
		},
		"error with trace, empty stack": {
			je: func() *jerr.JApiError {
				je := jerr.NewJApiError("fake error", file, 4)
				je.OccurredInFile(fs.NewFile("/fizz/buzz", "foo"), 1)
				return je
			}(),
			expected: `fake error
/foo/bar:2
/fizz/buzz:1`,
		},
		"error without trace, filled stack": {
			je: jerr.NewJApiError("fake error", file, 4),
			stack: []stackItem{
				{NewJApiScanner(fs.NewFile("/from/tracer", "aa\nbb\ncc")), 7},
			},
			expected: `fake error
/foo/bar:2
/from/tracer:3`,
		},
		"error with trace, filled stack": {
			je: func() *jerr.JApiError {
				je := jerr.NewJApiError("fake error", file, 4)
				je.OccurredInFile(fs.NewFile("/fizz/buzz", "foo"), 1)
				return je
			}(),
			stack: []stackItem{
				{NewJApiScanner(fs.NewFile("/from/tracer", "aa\nbb\ncc")), 7},
			},
			expected: `fake error
/foo/bar:2
/fizz/buzz:1`,
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			(&Stack{stack: c.stack}).AddIncludeTraceToError(c.je)

			assert.Equal(t, c.expected, c.je.Error())
		})
	}
}

func TestStack_ToDirectiveIncludeTracer(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		tracer := (&Stack{}).ToDirectiveIncludeTracer()
		assert.Equal(t, emptyTracer, tracer)
	})

	t.Run("filled", func(t *testing.T) {
		const (
			at1 = 42
			at2 = 100500
		)
		fs1 := fs.NewFile("/foo", "content")
		fs2 := fs.NewFile("/bar", "content")

		s := &Stack{}
		require.NoError(t, s.Push(NewJApiScanner(fs1), at1))
		require.NoError(t, s.Push(NewJApiScanner(fs2), at2))

		tracer := s.ToDirectiveIncludeTracer()
		require.IsType(t, directiveIncludeTracer{}, tracer)

		assert.Equal(t, []directiveIncludeTracerItem{
			{fs1, at1},
			{fs2, at2},
		}, tracer.(directiveIncludeTracer).stack)
	})
}

func Test_newDirectiveIncludeTracer(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		tracer := newDirectiveIncludeTracer(nil)
		assert.Equal(t, []directiveIncludeTracerItem{}, tracer.stack)
	})

	t.Run("filled", func(t *testing.T) {
		const (
			at1 = 42
			at2 = 100500
		)
		fs1 := fs.NewFile("/foo", "content")
		fs2 := fs.NewFile("/bar", "content")

		tracer := newDirectiveIncludeTracer([]stackItem{
			{NewJApiScanner(fs1), at1},
			{NewJApiScanner(fs2), at2},
		})
		assert.Equal(t, []directiveIncludeTracerItem{
			{fs1, at1},
			{fs2, at2},
		}, tracer.stack)
	})
}

func TestDirectiveIncludeTracer_AddIncludeTraceToError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		directiveIncludeTracer{}.AddIncludeTraceToError(nil)
	})

	file := fs.NewFile("/foo/bar", "12\n34\n56")

	cc := map[string]struct {
		je       *jerr.JApiError
		stack    []directiveIncludeTracerItem
		expected string
	}{
		"error without trace, empty stack": {
			je:       jerr.NewJApiError("fake error", file, 4),
			expected: "fake error",
		},
		"error with trace, empty stack": {
			je: func() *jerr.JApiError {
				je := jerr.NewJApiError("fake error", file, 4)
				je.OccurredInFile(fs.NewFile("/fizz/buzz", "foo"), 1)
				return je
			}(),
			expected: `fake error
/foo/bar:2
/fizz/buzz:1`,
		},
		"error without trace, filled stack": {
			je: jerr.NewJApiError("fake error", file, 4),
			stack: []directiveIncludeTracerItem{
				{fs.NewFile("/from/tracer", "aa\nbb\ncc"), 7},
			},
			expected: `fake error
/foo/bar:2
/from/tracer:3`,
		},
		"error with trace, filled stack": {
			je: func() *jerr.JApiError {
				je := jerr.NewJApiError("fake error", file, 4)
				je.OccurredInFile(fs.NewFile("/fizz/buzz", "foo"), 1)
				return je
			}(),
			stack: []directiveIncludeTracerItem{
				{fs.NewFile("/from/tracer", "aa\nbb\ncc"), 7},
			},
			expected: `fake error
/foo/bar:2
/fizz/buzz:1`,
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			directiveIncludeTracer{stack: c.stack}.AddIncludeTraceToError(c.je)

			assert.Equal(t, c.expected, c.je.Error())
		})
	}
}
