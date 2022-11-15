package core

import (
	"errors"
	"testing"

	schema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/internal/mocks"
)

func Test_fetchUsedUserTypes(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		userTypes := &catalog.UserSchemas{}

		s := mocks.NewSchema(t)
		s.On("UsedUserTypes").Return([]string{"foo", "bar"}, nil)

		foo := mocks.NewSchema(t)
		foo.On("UsedUserTypes").Return(nil, nil)

		bar := mocks.NewSchema(t)
		bar.On("UsedUserTypes").Return(nil, nil)

		userTypes.Set("foo", foo)
		userTypes.Set("bar", bar)

		actual, err := fetchUsedUserTypes(s, userTypes)
		require.NoError(t, err)
		assert.Equal(t, []string{"foo", "bar"}, actual)
	})

	t.Run("negative", func(t *testing.T) {
		s := mocks.NewSchema(t)
		s.On("UsedUserTypes").Return(nil, errors.New("fake error"))

		_, err := fetchUsedUserTypes(s, nil)
		assert.EqualError(t, err, "fake error")
	})
}

func TestUsedUserTypeFetcher_fetch(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given          func(*testing.T) schema.Schema
			setupUserTypes func(*testing.T) *catalog.UserSchemas

			expectedAlreadyProcessed map[string]struct{}
			expectedUsedUserTypes    []string
		}{
			"nil schema": {
				given:          func(*testing.T) schema.Schema { return nil },
				setupUserTypes: func(*testing.T) *catalog.UserSchemas { return nil },
			},

			"there are no any user types": {
				given: func(t *testing.T) schema.Schema {
					s := mocks.NewSchema(t)
					s.On("UsedUserTypes").Return(nil, nil)
					return s
				},
				setupUserTypes: func(*testing.T) *catalog.UserSchemas { return nil },
			},

			"there are user types": {
				given: func(t *testing.T) schema.Schema {
					s := mocks.NewSchema(t)
					s.On("UsedUserTypes").Return([]string{"foo", "bar"}, nil)
					return s
				},
				setupUserTypes: func(t *testing.T) *catalog.UserSchemas {
					us := &catalog.UserSchemas{}

					foo := mocks.NewSchema(t)
					foo.On("UsedUserTypes").Return([]string{"fizz", "buzz"}, nil)

					bar := mocks.NewSchema(t)
					bar.On("UsedUserTypes").Return(nil, nil)

					fizz := mocks.NewSchema(t)
					fizz.On("UsedUserTypes").Return([]string{"foo"}, nil)

					buzz := mocks.NewSchema(t)
					buzz.On("UsedUserTypes").Return(nil, nil)

					us.Set("foo", foo)
					us.Set("bar", bar)
					us.Set("fizz", fizz)
					us.Set("buzz", buzz)
					return us
				},
				expectedAlreadyProcessed: map[string]struct{}{
					"foo":  {},
					"bar":  {},
					"fizz": {},
					"buzz": {},
				},
				expectedUsedUserTypes: []string{
					"foo",
					"fizz",
					"buzz",
					"bar",
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				f := usedUserTypeFetcher{
					userTypes:        c.setupUserTypes(t),
					alreadyProcessed: map[string]struct{}{},
				}

				err := f.fetch(c.given(t))
				require.NoError(t, err)

				assert.Len(t, f.alreadyProcessed, len(c.expectedAlreadyProcessed))
				for k := range c.expectedAlreadyProcessed {
					assert.Contains(t, f.alreadyProcessed, k)
				}
				assert.Equal(t, c.expectedUsedUserTypes, f.usedUserTypes)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]struct {
			newSchema      func(*testing.T) schema.Schema
			setupUserTypes func(*testing.T) *catalog.UserSchemas
		}{
			"fake error": {
				newSchema: func(t *testing.T) schema.Schema {
					s := mocks.NewSchema(t)
					s.On("UsedUserTypes").Return(nil, errors.New("fake error"))
					return s
				},
				setupUserTypes: func(*testing.T) *catalog.UserSchemas {
					return nil
				},
			},

			`process type "foo": fake error`: {
				newSchema: func(t *testing.T) schema.Schema {
					s := mocks.NewSchema(t)
					s.On("UsedUserTypes").Return([]string{"foo"}, nil)
					return s
				},
				setupUserTypes: func(*testing.T) *catalog.UserSchemas {
					us := &catalog.UserSchemas{}

					s := mocks.NewSchema(t)
					s.On("UsedUserTypes").Return(nil, errors.New("fake error"))

					us.Set("foo", s)
					return us
				},
			},
		}

		for expected, c := range cc {
			t.Run(expected, func(t *testing.T) {
				f := usedUserTypeFetcher{
					userTypes:        c.setupUserTypes(t),
					alreadyProcessed: map[string]struct{}{},
				}

				err := f.fetch(c.newSchema(t))
				assert.EqualError(t, err, expected)
			})
		}
	})
}
