package directive

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/assert"
)

func TestNewCoords(t *testing.T) {
	file := fs.NewFile("foo", "123456")
	const begin = 1
	const end = 2

	coords := NewCoords(file, begin, end)

	assert.Same(t, file, coords.file)
	assert.Same(t, file, coords.File())
	assert.EqualValues(t, begin, coords.begin)
	assert.EqualValues(t, begin, coords.Begin())
	assert.EqualValues(t, end, coords.end)
}

func TestCoords_Read(t *testing.T) {
	bytes := NewCoords(fs.NewFile("foo", "123456"), 1, 3).Read()
	assert.EqualValues(t, "234", bytes)
}

func TestCoords_IsSet(t *testing.T) {
	cc := map[string]struct {
		coords   Coords
		expected bool
	}{
		"without file, without end": {Coords{begin: 1}, false},
		"without file, with end":    {Coords{begin: 1, end: 2}, false},
		"with file, without end":    {Coords{file: fs.NewFile("", "content"), begin: 1}, false},
		"with file, with end":       {Coords{file: fs.NewFile("", "content"), begin: 1, end: 2}, true},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			actual := c.coords.IsSet()
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestCoords_String(t *testing.T) {
	bytes := NewCoords(fs.NewFile("foo", "123456"), 1, 3).String()
	assert.EqualValues(t, "[1:3]", bytes)
}
