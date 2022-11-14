package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"os"
	"strings"
	"text/template"
	"unicode"
)

func typeToString(expr ast.Expr) (string, error) {
	switch n := expr.(type) {
	case *ast.StarExpr:
		s, err := typeToString(n.X)
		if err != nil {
			return "", err
		}
		return "*" + s, nil

	case *ast.Ident:
		return n.Name, nil

	case *ast.StructType:
		// Handle only empty struct without fields.
		return "struct{}", nil

	case *ast.SelectorExpr:
		pkg, err := typeToString(n.X)
		if err != nil {
			return "", err
		}

		typ, err := typeToString(n.Sel)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s.%s", pkg, typ), nil
	}

	return "", fmt.Errorf("unhanled expression %#v", expr)
}

func camelCaseToUnderscore(s string) string {
	// Assume we have no more than 3 capital letters except the first one in the
	// type name.
	buf := strings.Builder{}
	buf.Grow(len(s) + 3)

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i != 0 {
				buf.WriteRune('_')
			}
			r = unicode.ToLower(r)
		}
		buf.WriteRune(r)
	}

	// Replace common abbreviations.
	return strings.ReplaceAll(buf.String(), "a_s_t", "ast")
}

func renderTemplateToFile(tmpl, path string, data interface{}) error {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 2048))
	if err = t.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	code, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to gofmt: %w", err)
	}

	return os.WriteFile(path, code, 0644) //nolint:gosec // It's okay, we save a code here.
}
