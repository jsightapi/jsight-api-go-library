package directive

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func TestNew(t *testing.T) {
	e := Info
	coords := Coords{}
	d := New(e, coords)

	assert.EqualValues(t, e, d.type_)
	assert.EqualValues(t, e, d.Type())

	assert.NotNil(t, d.namedParameters)
	assert.Equal(t, e.String(), d.Keyword)
	assert.Equal(t, coords, d.keywordCoords)
	assert.IsType(t, nopIncludeTracer{}, d.includeTracer)
}

type fakeCallStack struct{}

func (fakeCallStack) AddIncludeTraceToError(*jerr.JApiError) {}

func TestNewWithCallStack(t *testing.T) {
	e := Info
	coords := Coords{}
	d := NewWithCallStack(e, coords, fakeCallStack{})

	assert.EqualValues(t, e, d.type_)
	assert.EqualValues(t, e, d.Type())

	assert.NotNil(t, d.namedParameters)
	assert.Equal(t, e.String(), d.Keyword)
	assert.Equal(t, coords, d.keywordCoords)
	assert.IsType(t, fakeCallStack{}, d.includeTracer)
}

func TestDirective_String(t *testing.T) {
	e := Info
	actual := New(e, Coords{})

	assert.Equal(t, e.String(), actual.String())
}

func TestDirective_Equal(t *testing.T) {
	const content = "content"
	file := fs.NewFile("file", content)

	cc := map[string]struct {
		x, y     Directive
		expected bool
	}{
		"empty": {expected: true},
		"different file": {
			x: Directive{
				keywordCoords: NewCoords(fs.NewFile("foo", content), 0, 1),
			},
			y: Directive{
				keywordCoords: NewCoords(fs.NewFile("bar", content), 0, 1),
			},
		},
		"different begin": {
			x: Directive{
				keywordCoords: NewCoords(file, 0, 1),
			},
			y: Directive{
				keywordCoords: NewCoords(file, 10, 1),
			},
		},
		"equal": {
			x: Directive{
				keywordCoords: NewCoords(file, 0, 1),
			},
			y: Directive{
				keywordCoords: NewCoords(file, 0, 1),
			},
			expected: true,
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			actual := c.x.Equal(c.y)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestDirective_HasNamedParameters(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		actual := Directive{
			namedParameters: map[string]string{"foo": "bar"},
		}.HasNamedParameter()
		assert.True(t, actual)
	})

	t.Run("false", func(t *testing.T) {
		actual := Directive{}.HasNamedParameter()
		assert.False(t, actual)
	})
}

func TestDirective_NamedParameter(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		d := New(Jsight, Coords{})
		require.NoError(t, d.SetNamedParameter("foo", "bar"))

		actual := d.HasNamedParameter()
		assert.True(t, actual)
	})

	t.Run("not exists", func(t *testing.T) {
		actual := New(Jsight, Coords{}).HasNamedParameter()
		assert.False(t, actual)
	})
}

func TestDirective_SetNamedParameter(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		d := New(Info, Coords{})

		err := d.SetNamedParameter("foo", "bar")
		require.NoError(t, err)

		assert.Equal(t, map[string]string{
			"foo": "bar",
		}, d.namedParameters)
	})

	t.Run("negative", func(t *testing.T) {
		d := New(Info, Coords{})

		err := d.SetNamedParameter("foo", "bar")
		require.NoError(t, err)

		err = d.SetNamedParameter("foo", "bar")
		assert.EqualError(t, err, `the "foo" parameter is already defined for the "INFO" directive`)

		assert.Equal(t, map[string]string{
			"foo": "bar",
		}, d.namedParameters)
	})
}

func TestDirective_AppendChild(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		d := &Directive{}
		c := &Directive{}

		d.AppendChild(c)

		assert.Equal(t, []*Directive{c}, d.Children)
	})

	t.Run("filled", func(t *testing.T) {
		c1 := &Directive{}
		c2 := &Directive{}
		c3 := &Directive{}

		d := &Directive{Children: []*Directive{
			c1,
			c2,
		}}

		d.AppendChild(c3)

		assert.Equal(t, []*Directive{c1, c2, c3}, d.Children)
	})
}

func TestDirective_CopyWoParentAndChildren(t *testing.T) {
	c := &Directive{}

	d := New(Jsight, Coords{})
	d.Annotation = "annotation"
	d.HasExplicitContext = true
	d.namedParameters = map[string]string{"foo": "bar"}
	d.BodyCoords = NewCoords(nil, 0, 1)
	d.Children = []*Directive{c}
	d.includeTracer = fakeCallStack{}

	actual := d.CopyWoParentAndChildren()
	assert.Equal(t, Directive{
		type_:              d.type_,
		Annotation:         d.Annotation,
		Keyword:            d.Keyword,
		HasExplicitContext: d.HasExplicitContext,
		namedParameters:    d.namedParameters,
		keywordCoords:      d.keywordCoords,
		BodyCoords:         d.BodyCoords,
		includeTracer:      fakeCallStack{},
	}, actual)
}
