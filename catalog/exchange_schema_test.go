package catalog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchema_MarshalJSON_Order(t *testing.T) {
	cc := []struct {
		schema   string
		expected []string
	}{
		{
			`{"id": 1, "name": "Tom", "age": 6}`,
			[]string{"id", "name", "age"},
		},
		{
			`{"id": 1, "age": 6, "name": "Tom"}`,
			[]string{"id", "age", "name"},
		},
		{
			`{"name": "Tom", "id": 1, "age": 6}`,
			[]string{"name", "id", "age"},
		},
		{
			`{"name": "Tom", "age": 6, "id": 1}`,
			[]string{"name", "age", "id"},
		},
		{
			`{"age": 6, "id": 1, "name": "Tom"}`,
			[]string{"age", "id", "name"},
		},
		{
			`{"age": 6, "name": "Tom", "id": 1}`,
			[]string{"age", "name", "id"},
		},
	}

	for _, c := range cc {
		t.Run(c.schema, func(t *testing.T) {
			s, err := NewExchangeJSightSchema("", []byte(c.schema), &UserSchemas{}, nil, &UserTypes{})
			require.NoError(t, err)

			err = s.Compile()
			require.NoError(t, err)

			ss := make([]string, 0, len(s.exchangeContent.Children))
			for _, v := range s.exchangeContent.Children {
				ss = append(ss, *(v.Key))
			}
			assert.Equal(t, c.expected, ss)
		})
	}
}
