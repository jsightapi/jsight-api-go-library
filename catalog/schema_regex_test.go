package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/notation"
)

func TestRegexMarshaller_Marshal(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s, err := regexMarshaller{
			useFixedSeed: true,
		}.
			Marshal("foo", []byte("/bar-\\d/"))
		require.NoError(t, err)

		assert.Equal(t, Schema{
			Notation:      notation.SchemaNotationRegex,
			ContentRegexp: "bar-\\d",
			Example:       "bar-4",
			UsedUserTypes: &StringSet{},
			UsedUserEnums: &StringSet{},
		}, s)
	})

	t.Run("negative", func(t *testing.T) {
		_, err := regexMarshaller{
			useFixedSeed: true,
		}.
			Marshal("", []byte("invalid"))
		assert.EqualError(t, err, `ERROR (code 1500): Regex should starts with '/' character, but found 'i'
	in line 1 on file 
	> invalid
	--^`)
	})
}
