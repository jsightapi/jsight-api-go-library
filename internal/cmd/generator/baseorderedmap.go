package main

import (
	"errors"
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
)

type baseOrderedMapGenerator struct {
	template string
}

type orderedMap struct {
	Name        string
	PkgName     string
	KeyType     string
	ValueType   string
	UsedImports map[string]struct{}
}

func (baseOrderedMapGenerator) collectUsedTypes(data, order *ast.Field, m *orderedMap) error {
	mapType, ok := data.Type.(*ast.MapType)
	if !ok {
		return errors.New(`"data" field should be a map`)
	}

	var err error

	m.KeyType, err = typeToString(mapType.Key)
	if err != nil {
		return fmt.Errorf(`failed to get "data" map key type: %w`, err)
	}

	m.ValueType, err = typeToString(mapType.Value)
	if err != nil {
		return fmt.Errorf(`failed to get "data" map value type: %w`, err)
	}

	slice, ok := order.Type.(*ast.ArrayType)
	if !ok {
		return errors.New(`"order" field should be a slice`)
	}

	sliceType, err := typeToString(slice.Elt)
	if err != nil {
		return fmt.Errorf(`failed to get "order" slice item type: %w`, err)
	}

	if sliceType != m.KeyType {
		return fmt.Errorf(
			`"order" slice item type %q isn't equal to %q`,
			sliceType,
			m.KeyType,
		)
	}

	return nil
}

func (g baseOrderedMapGenerator) fillImports(om *orderedMap, imports map[string]string) error {
	if pkg := g.getTypePackage(om.ValueType); pkg != "" {
		p, ok := imports[pkg]
		if !ok {
			return fmt.Errorf("failed to find import for type %q", om.ValueType)
		}
		om.UsedImports[p] = struct{}{}
	}

	if pkg := g.getTypePackage(om.KeyType); pkg != "" {
		p, ok := imports[pkg]
		if !ok {
			return fmt.Errorf("failed to find import for type %q", om.KeyType)
		}
		om.UsedImports[p] = struct{}{}
	}
	return nil
}

func (baseOrderedMapGenerator) getTypePackage(t string) string {
	parts := strings.SplitN(t, ".", 2)
	if len(parts) != 2 {
		return ""
	}

	conversationMap := map[string]string{
		"jschema": "schema",
	}

	if p, ok := conversationMap[parts[0]]; ok {
		return p
	}

	return parts[0]
}

func (g baseOrderedMapGenerator) generateCode(om orderedMap, dirPath string) error {
	path := filepath.Join(dirPath, camelCaseToUnderscore(om.Name)+"_gen.go")
	return renderTemplateToFile(g.template, path, om)
}
