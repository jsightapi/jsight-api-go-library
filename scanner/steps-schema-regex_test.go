package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_stateRegex(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := Scanner{}
		assert.Nil(t, s.step)

		err := stateRegex(&s, '/')
		require.Nil(t, err)
		assert.NotNil(t, s.step)
	})

	t.Run("negative", func(t *testing.T) {
		s := newTestScanner("1")
		s.step = nil

		err := stateRegex(s, '1')
		require.NotNil(t, err)
		assert.Equal(t, "invalid character '1' in the regular expression, expecting '/' character", err.Msg)

		assert.Nil(t, s.step)
	})
}

func Test_stateRegexFirstChar(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := Scanner{}
		assert.Nil(t, s.step)

		err := stateRegexFirstChar(&s, '1')
		require.Nil(t, err)
		assert.NotNil(t, s.step)
	})

	t.Run("negative", func(t *testing.T) {
		s := newTestScanner("/")
		s.step = nil

		err := stateRegexFirstChar(s, '/')
		require.NotNil(t, err)
		assert.Equal(t, "invalid character '/' empty regex", err.Msg)

		assert.Nil(t, s.step)
	})
}

func Test_stateRegexBody(t *testing.T) {
	t.Run("/", func(t *testing.T) {
		s := Scanner{
			curIndex: 1,
		}

		err := stateRegexBody(&s, '/')
		require.Nil(t, err)

		assert.Equal(t, []LexemeEvent{
			{TextEnd, 1},
		}, s.finds)
		assert.NotNil(t, s.step)
	})

	t.Run("\\", func(t *testing.T) {
		s := Scanner{
			curIndex: 1,
		}

		err := stateRegexBody(&s, '\\')
		require.Nil(t, err)

		assert.Nil(t, s.finds)
		assert.NotNil(t, s.step)
	})

	t.Run("ordinal", func(t *testing.T) {
		s := Scanner{
			curIndex: 1,
		}

		err := stateRegexBody(&s, '1')
		require.Nil(t, err)

		assert.Nil(t, s.finds)
		assert.Nil(t, s.step)
	})
}

func Test_stateRegexBodyAfterSlash(t *testing.T) {
	s := Scanner{}
	err := stateRegexBodyAfterSlash(&s, '/')
	require.Nil(t, err)
}
