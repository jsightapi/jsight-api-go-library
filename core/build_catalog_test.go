package core

import (
	jschema "j/schema"
	"j/schema/bytes"
	"j/schema/fs"
	"testing"

	"j/japi/catalog"
	"j/japi/directive"

	"github.com/stretchr/testify/require"
)

func TestJApiCore_compileUserTypes(t *testing.T) {
	c := catalog.NewCatalog()

	ut := map[string][]byte{
		"@foo": []byte("42"),
		"@bar": []byte(`{
	"foo": @foo
}`),
		"@fizz": []byte(`{
		"bar": @bar
}`),
		"@buzz": []byte(`{
		"fizz": @fizz
}`),
	}

	for n, p := range ut {
		d := directive.New(directive.Jsight, directive.Coords{})
		d.BodyCoords = directive.NewCoords(
			fs.NewFile(n, p),
			0,
			bytes.Index(len(p)-1),
		)
		require.NoError(t, d.SetParameter("Name", n))
		c.AddRawUserType(d)
	}

	core := JApiCore{
		userTypes:          &catalog.UserSchemas{},
		processedUserTypes: map[string]struct{}{},
		catalog:            c,
	}

	err := core.compileUserTypes()
	require.Nil(t, err)

	_ = core.userTypes.Each(func(k string, v jschema.Schema) error {
		require.NoErrorf(t, v.Check(), "Check %q user type", k)
		return nil
	})
}
