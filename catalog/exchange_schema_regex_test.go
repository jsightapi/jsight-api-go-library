package catalog

import (
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexMarshaller_Marshal(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s, err := PrepareRegexSchema("foo", []byte("/bar-\\d/"))
		require.NoError(t, err)

		n, err := s.GetAST()
		require.NoError(t, err)
		assert.Equal(t, jschemaLib.ASTNode{
			TokenType:  jschemaLib.TokenTypeString,
			SchemaType: jschemaLib.TokenTypeString,
			Value:      "/bar-\\d/",
		}, n)

		example, err := s.Example()
		require.NoError(t, err)
		assert.Equal(t, []byte("bar-4"), example)
	})

	t.Run("negative", func(t *testing.T) {
		s, err := PrepareRegexSchema("", []byte("invalid"))
		require.NoError(t, err)

		_, err = s.GetAST()
		assert.EqualError(t, err, `ERROR (code 1500): Regex should starts with '/' character, but found 'i'
	in line 1 on file 
	> invalid
	--^`)
	})
}
