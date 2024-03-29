package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

// processContext resolves context according to incoming directive
func (core *JApiCore) processContext(d *directive.Directive, root *[]*directive.Directive) *jerr.JApiError {
	for {
		if core.currentContextDirective == nil { // root context
			if d.Type().IsAllowedForRootContext() {
				*root = append(*root, d)
				core.currentContextDirective = d
				return nil
			}
			return d.KeywordError(fmt.Sprintf("%s %q", jerr.IncorrectContextOfDirective, d.String()))
		}

		// not the root context
		if core.currentContextDirective.Type().IsAllowedForDirectiveContext(d.Type()) {
			isURL := d.Type().IsHTTPRequestMethod() &&
				d.NamedParameter("Path") != "" &&
				core.currentContextDirective.Type() == directive.URL

			if isURL {
				if core.currentContextDirective.HasExplicitContext {
					return d.KeywordError(fmt.Sprintf(
						"%s %q with the \"Path\" parameter",
						jerr.IncorrectContextOfDirective,
						d.String(),
					))
				}
				*root = append(*root, d)
				core.currentContextDirective = d
				return nil
			}

			d.Parent = core.currentContextDirective
			core.currentContextDirective.AppendChild(d)
			core.currentContextDirective = d
			return nil
		}

		if core.currentContextDirective.HasExplicitContext {
			return d.KeywordError(fmt.Sprintf("%s %q", jerr.IncorrectContextOfDirective, d.String()))
		}
		core.currentContextDirective = core.currentContextDirective.Parent
	}
}
