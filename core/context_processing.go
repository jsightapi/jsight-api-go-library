package core

import (
	"fmt"
	"j/japi/directive"
	"j/japi/jerr"
)

// processContext resolves context according to incoming directive
func (core *JApiCore) processContext(d *directive.Directive, root *[]*directive.Directive) *jerr.JAPIError {
	for {
		if core.currentContextDirective == nil { // root context
			if directive.IsAllowedForRootContext(d.Type()) {
				// d.Parent = nil
				*root = append(*root, d)
				core.currentContextDirective = d
				return nil
			} else {
				return d.KeywordError(fmt.Sprintf("%s %q", jerr.IncorrectContextOfDirective, d.String()))
			}
		} else {
			if core.currentContextDirective.Type().IsAllowedForDirectiveContext(d.Type()) {
				d.Parent = core.currentContextDirective
				core.currentContextDirective.AppendChild(d)
				core.currentContextDirective = d
				return nil
			} else {
				if core.currentContextDirective.HasExplicitContext {
					return d.KeywordError(fmt.Sprintf("%s %q", jerr.IncorrectContextOfDirective, d.String()))
				}
				core.currentContextDirective = core.currentContextDirective.Parent
			}
		}
	}
}
