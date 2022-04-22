package core

import (
	"j/japi/catalog"
	"j/japi/directive"
	"j/japi/jerr"
	"j/japi/scanner"
	"j/schema/fs"
)

// JApiCore the Brain and heart of jApi. Collects lexemes from scanner, validates document logic and structure,
// builds catalog, renders documentation
type JApiCore struct {
	// userTypes represent all user types.
	userTypes *catalog.UserSchemas

	// processedUserTypes a "set" of already processed user types. Used for resolving
	// circular dependencies between user types.
	processedUserTypes map[string]struct{}

	// representation of Api data.
	catalog *catalog.Catalog

	// scanner for parsing .jst files.
	scanner *scanner.Scanner

	// current file that core works with.
	file *fs.File

	// currentContextDirective is current context for adding a child directive.
	currentContextDirective *directive.Directive

	// to accumulate directive data pieces from scanner.
	currentDirective *directive.Directive

	// directives from loaded from project.
	directives []*directive.Directive

	// directives after processing the MACRO and PASTE directives.
	directivesWithPastes []*directive.Directive

	// uniqURLPath used for checking the uniqueness URL paths.
	uniqURLPath map[catalog.Path]bool

	// similarPaths to check the forbidding of "similar" paths.
	similarPaths map[string]string

	// rawPathVariables contains properties of the Path directive.
	rawPathVariables []rawPathVariable

	// macro contains list of all project macros.
	macro map[string]*directive.Directive
}

func NewJApiCore(file *fs.File) *JApiCore {
	c := &JApiCore{
		userTypes:          &catalog.UserSchemas{},
		processedUserTypes: make(map[string]struct{}, 30),
		file:               file,
		catalog:            catalog.NewCatalog(),
		currentDirective:   nil,
		directives:         make([]*directive.Directive, 0, 200),
		uniqURLPath:        make(map[catalog.Path]bool, 20),
		similarPaths:       make(map[string]string, 20),
		rawPathVariables:   make([]rawPathVariable, 0, 40),
		macro:              make(map[string]*directive.Directive, 20),
	}
	return c
}

// Main internal method of Core. Only when done, core is ready to handle data validation or doc rendering.
func (core *JApiCore) processJapiProject() *jerr.JAPIError {
	if err := core.scanProject(); err != nil {
		return err
	}

	if err := core.compileCore(); err != nil {
		return err
	}

	if err := core.buildCatalog(); err != nil {
		return err
	}

	if err := core.compileCatalog(); err != nil {
		return err
	}

	return core.validateCatalog()
}
