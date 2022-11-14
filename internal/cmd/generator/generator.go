package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
)

type byCommentGenerator struct {
	commonGenerator
}

type byCommentGeneratorVisitorFn = func(
	pkgName, path string,
	spec *ast.TypeSpec,
	imports map[string]string,
) error

func newByCommentGenerator(
	markerComment string,
	visitor byCommentGeneratorVisitorFn,
) byCommentGenerator {
	return byCommentGenerator{
		commonGenerator: commonGenerator{
			visitor: func(
				pkgName, path string,
				decl *ast.GenDecl,
				imports map[string]string,
			) error {
				if !strings.Contains(decl.Doc.Text(), markerComment) {
					return nil
				}

				for _, s := range decl.Specs {
					spec, ok := s.(*ast.TypeSpec)
					if !ok {
						continue
					}

					if err := visitor(pkgName, path, spec, imports); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}
}

type commonGenerator struct {
	visitor func(
		pkgName, path string,
		decl *ast.GenDecl,
		imports map[string]string,
	) error
}

func (g commonGenerator) Generate(p string) error {
	pkgName, dd, err := g.parseFile(p)
	if err != nil {
		return fmt.Errorf("faile to parse file: %w", err)
	}

	if err := g.generate(pkgName, p, dd); err != nil {
		return fmt.Errorf("failed to find target types: %w", err)
	}

	return nil
}

func (commonGenerator) parseFile(p string) (pkgName string, dd []ast.Decl, err error) {
	const flags = parser.ParseComments | parser.AllErrors
	f, err := parser.ParseFile(token.NewFileSet(), p, nil, flags)
	if err != nil {
		return "", nil, err
	}

	return f.Name.Name, f.Decls, nil
}

func (g commonGenerator) generate(pkgName, filePath string, dd []ast.Decl) error {
	imports := map[string]string{}

	for _, d := range dd {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			return nil
		}

		if _, ok := decl.Specs[0].(*ast.ImportSpec); ok {
			for _, s := range decl.Specs {
				spec, ok := s.(*ast.ImportSpec)
				if !ok {
					continue
				}

				p := strings.Trim(spec.Path.Value, `"`)
				imports[path.Base(p)] = p
			}
			continue
		}

		if err := g.visitor(pkgName, filePath, decl, imports); err != nil {
			return err
		}
	}
	return nil
}
