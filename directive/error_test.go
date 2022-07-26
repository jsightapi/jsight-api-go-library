package directive

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/assert"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func TestDirective_KeywordError(t *testing.T) {
	coords := NewCoords(fs.NewFile("foo", "123456"), 1, 2)

	je := New(Jsight, coords).KeywordError("fake error")
	assert.Equal(t, jerr.NewJApiError("fake error", coords.file, coords.begin), je)
}

func TestDirective_BodyError(t *testing.T) {
	coords := NewCoords(fs.NewFile("foo", "123456"), 1, 2)

	t.Run("body coords is set", func(t *testing.T) {
		d := New(Jsight, NewCoords(fs.NewFile("bar", "123456"), 5, 6))
		d.BodyCoords = coords

		je := d.BodyError("fake error")
		assert.Equal(t, jerr.NewJApiError("fake error", coords.file, coords.begin), je)
	})

	t.Run("body coords isn't set", func(t *testing.T) {
		d := New(Jsight, coords)

		je := d.BodyError("fake error")
		assert.Equal(t, jerr.NewJApiError("fake error", coords.file, coords.begin), je)
	})
}

func TestDirective_BodyErrorIndex(t *testing.T) {
	const idx = 3
	coords := NewCoords(fs.NewFile("foo", "123456"), 1, 2)

	d := New(Jsight, Coords{})
	d.BodyCoords = coords

	je := d.BodyErrorIndex("fake error", idx)
	assert.Equal(t, jerr.NewJApiError("fake error", coords.file, coords.begin+idx), je)
}

func TestDirective_ParameterError(t *testing.T) {
	coords := NewCoords(fs.NewFile("foo", "123456"), 1, 2)

	je := New(Jsight, coords).ParameterError("fake error")
	assert.Equal(t, jerr.NewJApiError("fake error", coords.file, coords.begin), je)
}

type makeErrorCallStack struct {
	t          *testing.T
	expectedJE *jerr.JApiError
}

func (m makeErrorCallStack) AddIncludeTraceToError(je *jerr.JApiError) {
	assert.Equal(m.t, m.expectedJE, je)
}

func TestDirective_makeError(t *testing.T) {
	coords := NewCoords(fs.NewFile("foo", "123456"), 1, 2)
	expected := jerr.NewJApiError("fake error", coords.File(), coords.begin)

	d := New(Jsight, coords)
	d.includeTracer = makeErrorCallStack{t, expected}

	je := d.makeError("fake error", coords.File(), coords.begin)

	assert.Equal(t, expected, je)
}
