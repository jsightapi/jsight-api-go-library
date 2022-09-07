package main

import (
	"errors"
	"fmt"
	"go/ast"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// orderedMapGenerator generator will search for `// gen:OrderedMap` comments for
// custom types and generate all necessary code to this type.
//
// Requirements:
//
//   - Custom type should be a struct with exactly two fields: "data" anf "order";
//   - Field "data" should be a map;
//   - Field "order" should be a slice of map keys;
//   - Field "mx" should a sync.RWMutex.
//
// Known limitations:
//   - Unfortunately "omitempty" tag won't work as expected for ordered maps, so
//     you should add specific code for marshaling.
//
// Added because we should preserve keys order in the map, but Golang don't
// gave to us such ability out of the box. And we want more efficient and clean
// way to manipulate such maps, so empty interface isn't an option.
type orderedMapGenerator struct {
	baseOrderedMapGenerator
	byCommentGenerator
}

var _ generator = orderedMapGenerator{}

func newOrderMapGenerator() orderedMapGenerator {
	g := orderedMapGenerator{
		baseOrderedMapGenerator: baseOrderedMapGenerator{
			template: `// Autogenerated code!
// DO NOT EDIT!
//
// Generated by OrderedMap generator from the internal/cmd/generator command.

package {{ .PkgName }}

import (
	"bytes"
	"encoding/json"
{{ range $k, $v := .UsedImports }}
	"{{ $k }}"
{{ end }}
)

// Set sets a value with specified key.
func (m *{{ .Name }}) Set(k {{ .KeyType }}, v {{ .ValueType }}) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[{{ .KeyType }}]{{ .ValueType }}{}
	}
	if !m.has(k) {
		m.order = append(m.order, k)
	}
	m.data[k] = v
}

// SetToTop do the same as Set, but new key will be placed on top of the order
// map.
func (m *{{ .Name }}) SetToTop(k {{ .KeyType }}, v {{ .ValueType }}) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[{{ .KeyType }}]{{ .ValueType }}{}
	}
	if !m.has(k) {
		m.order = append([]{{ .KeyType }}{k}, m.order...)
	}
	m.data[k] = v
}

// Update updates a value with specified key.
func (m *{{ .Name }}) Update(k {{ .KeyType }}, fn func(v {{ .ValueType }}) {{ .ValueType }}) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if !m.has(k) {
		// Prevent from possible nil pointer dereference if map value type is a
		// pointer.
		return
	}

	m.data[k] = fn(m.data[k])
}

// GetValue gets a value by key.
func (m *{{ .Name }}) GetValue(k {{ .KeyType }}) {{ .ValueType }} {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.data[k]
}

// Get gets a value by key.
func (m *{{ .Name }}) Get(k {{ .KeyType }}) ({{ .ValueType }}, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	v, ok := m.data[k]
	return v, ok
}

// Has checks that specified key is set.
func (m *{{ .Name }}) Has(k {{ .KeyType }}) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.has(k)
}

func (m *{{ .Name }}) has(k {{ .KeyType }}) bool {
	_, ok := m.data[k]
	return ok
}

// Len returns count of values.
func (m *{{ .Name }}) Len() int {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return len(m.data)
}

// Find finds first matched item from the map.
func (m *{{ .Name }}) Find(fn find{{ .CapitalizedName }}Func) ({{ .Name }}Item, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		if fn(k, m.data[k]) {
			return {{ .Name }}Item{
				Key:   k,
				Value: m.data[k],
			}, true
		}
	}
	return {{ .Name }}Item{}, false
}

type find{{ .CapitalizedName }}Func = func(k {{ .KeyType }}, v {{ .ValueType }}) bool

// Each iterates and perform given function on each item in the map.
func (m *{{ .Name }}) Each(fn each{{ .CapitalizedName }}Func) error {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		if err := fn(k, m.data[k]); err != nil {
			return err
		}
	}
	return nil
}

// EachReverse act almost the same as Each but in reverse order.
func (m *{{ .Name }}) EachReverse(fn each{{ .CapitalizedName }}Func) error {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for i := len(m.order) - 1; i >= 0; i-- {
		k := m.order[i]
		if err := fn(k, m.data[k]); err != nil {
			return err
		}
	}
	return nil
}

type each{{ .CapitalizedName }}Func = func(k {{ .KeyType }}, v {{ .ValueType }}) error

func (m *{{ .Name }}) EachSafe(fn eachSafe{{ .CapitalizedName }}Func) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		fn(k, m.data[k])
	}
}

type eachSafe{{ .CapitalizedName }}Func = func(k {{ .KeyType }}, v {{ .ValueType }})

// Map iterates and changes values in the map.
func (m *{{ .Name }}) Map(fn map{{ .CapitalizedName }}Func) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	for _, k := range m.order {
		v, err := fn(k, m.data[k])
		if err != nil {
			return err
		}
		m.data[k] = v
	}
	return nil
}

type map{{ .CapitalizedName }}Func = func(k {{ .KeyType }}, v {{ .ValueType }}) ({{ .ValueType }}, error)

// {{ .Name }}Item represent single data from the {{ .Name }}.
type {{ .Name }}Item struct {
	Key   {{ .KeyType }}
	Value {{ .ValueType }}
}

var _ json.Marshaler = &{{ .Name }}{}

func (m *{{ .Name }}) MarshalJSON() ([]byte, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	var buf bytes.Buffer
	buf.WriteRune('{')

	for i, k := range m.order {
		if i != 0 {
			buf.WriteRune(',')
		}

		// marshal key
		key, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteRune(':')

		// marshal value
		val, err := json.Marshal(m.data[k])
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}
`,
		},
	}
	g.byCommentGenerator = newByCommentGenerator(
		"gen:OrderedMap",
		func(pkgName, path string, spec *ast.TypeSpec, imports map[string]string) error {
			om, err := g.collectOrderMap(pkgName, spec, imports)
			if err != nil {
				return fmt.Errorf("failed to process type %q: %w", spec.Name, err)
			}

			if err := g.generateCode(om, filepath.Dir(path)); err != nil {
				return fmt.Errorf("failed to generate code for type %q: %w", om.Name, err)
			}
			return nil
		},
	)
	return g
}

func (orderedMapGenerator) Name() string { return "OrderedMap" }

func (g orderedMapGenerator) collectOrderMap(
	pkgName string,
	spec *ast.TypeSpec,
	imports map[string]string,
) (orderedMap, error) {
	strct, ok := spec.Type.(*ast.StructType)
	if !ok {
		return orderedMap{}, nil
	}

	if strct.Fields.NumFields() != 3 {
		return orderedMap{}, errors.New(`OrderedMap should have exactly three fields: "data", "order", and "mutex"`)
	}

	var (
		dataField  *ast.Field
		orderField *ast.Field
		mutexField *ast.Field
	)

	for _, f := range strct.Fields.List {
		switch f.Names[0].Name {
		case dataPropertyName:
			dataField = f

		case orderPropertyName:
			orderField = f

		case mutexPropertyName:
			mutexField = f
		}
	}

	om := orderedMap{
		Name:            spec.Name.Name,
		CapitalizedName: cases.Title(language.English, cases.NoLower).String(spec.Name.Name),
		PkgName:         pkgName,
		UsedImports:     map[string]struct{}{},
	}

	if err := g.checkMutexField(mutexField); err != nil {
		return orderedMap{}, err
	}

	if err := g.collectUsedTypes(dataField, orderField, &om); err != nil {
		return orderedMap{}, err
	}

	if err := g.fillImports(&om, imports); err != nil {
		return orderedMap{}, err
	}

	return om, nil
}

func (orderedMapGenerator) checkMutexField(f *ast.Field) error {
	if f == nil {
		return errors.New(`"mutex" field didn't present'`)
	}

	se, ok := f.Type.(*ast.SelectorExpr)
	if !ok {
		return errors.New(`"mutex" field should be *sync.RWMutex`)
	}

	if x, ok := se.X.(*ast.Ident); !ok || x.Name != "sync" || se.Sel.Name != "RWMutex" {
		return errors.New(`"mutex" field should be *sync.RWMutex`)
	}

	return nil
}
