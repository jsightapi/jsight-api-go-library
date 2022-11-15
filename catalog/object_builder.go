package catalog

import (
	schema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/loader"
	"github.com/jsightapi/jsight-schema-go-library/panics"
)

type ObjectBuilder struct {
	schema   *jschema.JSchema
	rootNode *ischema.ObjectNode
}

// NewObjectBuilder used only for building Path variables in the JSight API library
func NewObjectBuilder() ObjectBuilder {
	objNode := ischema.NewObjectNode(lexeme.LexEvent{})

	inner := ischema.New()
	inner.SetRootNode(objNode)

	s := jschema.New("", "")
	s.Inner = &inner

	return ObjectBuilder{
		schema:   s,
		rootNode: objNode,
	}
}

func (b ObjectBuilder) AddProperty(key string, node ischema.Node, types map[string]ischema.Type) {
	k := ischema.ObjectNodeKey{
		Key:        key,
		IsShortcut: false,
		Lex:        lexeme.LexEvent{},
	}
	b.rootNode.AddChild(k, node)
	for kk, vv := range types {
		b.schema.Inner.AddType(kk, vv)
	}
}

func (b ObjectBuilder) Len() int {
	return b.rootNode.Len()
}

func (b ObjectBuilder) UserTypeNames() []string {
	b.schema.CollectUserTypes()
	return b.schema.UsedUserTypes_.Data()
}

func (b ObjectBuilder) AddType(name string, sc schema.Schema) {
	if s, ok := sc.(*jschema.JSchema); ok {
		b.schema.Inner.AddType(name, ischema.Type{
			Schema:   s.Inner,
			RootFile: s.File,
		})
	}
}

func (b ObjectBuilder) Build() *jschema.JSchema {
	s := b.schema
	_ = s.LoadOnce.Do(func() (err error) {
		defer func() {
			err = panics.Handle(recover(), err)
		}()
		s.ASTNode = s.BuildASTNode()
		loader.CompileBasic(s.Inner, s.AreKeysOptionalByDefault)
		return nil
	})

	_ = s.Compile()

	return s
}
