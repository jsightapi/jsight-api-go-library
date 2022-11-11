package core

import (
	"testing"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/require"
)

func TestJApiCore_compileUserTypes(t *testing.T) {
	c := catalog.NewCatalog()

	ut := map[string]string{
		"@foo": "42",
		"@bar": `{
	"foo": @foo
}`,
		"@fizz": `{
		"bar": @bar
}`,
		"@buzz": `{
		"fizz": @fizz
}`,
	}

	core := JApiCore{
		rawUserTypes:       &directive.Directives{},
		userTypes:          &catalog.UserSchemas{},
		processedUserTypes: map[string]struct{}{},
		catalog:            c,
	}

	for n, p := range ut {
		d := directive.New(directive.Jsight, directive.Coords{})
		d.BodyCoords = directive.NewCoords(
			fs.NewFile(n, p),
			0,
			bytes.Index(len(p)-1),
		)
		require.NoError(t, d.SetNamedParameter("Name", n))
		core.AddRawUserType(d)
	}

	err := core.compileUserTypes()
	require.Nil(t, err)

	_ = core.userTypes.Each(func(k string, v jschemaLib.Schema) error {
		require.NoErrorf(t, v.Check(), "Check %q user type", k)
		return nil
	})
}
