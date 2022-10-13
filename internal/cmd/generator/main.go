// This generator should be used for generating some source code.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	gg := []generator{
		newOrderMapGenerator(),
		newUnsafeOrderMapGenerator(),
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()
		if info.IsDir() {
			if shouldIgnoreDir(name) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldIgnoreFile(name) {
			return nil
		}

		log.Printf("Process file %q", path)
		for _, g := range gg {
			if err := g.Generate(path); err != nil {
				return fmt.Errorf("%q failed to process file %q: %w", g.Name(), path, err)
			}
		}
		return nil
	})
}

func shouldIgnoreDir(name string) bool {
	return (name == "vendor") ||
		(name[0] == '.') ||
		(name == "cmd") ||
		(name == "tests") ||
		(name == "testdata")
}

func shouldIgnoreFile(name string) bool {
	return !strings.HasSuffix(name, ".go") ||
		strings.HasSuffix(name, "_test.go") ||
		strings.HasSuffix(name, "_gen.go")
}

type generator interface {
	Name() string
	Generate(path string) error
}
