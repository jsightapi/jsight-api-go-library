package core

import (
	"fmt"
	"j/japi/directive"
	"j/japi/jerr"
)

func (core *JApiCore) processPaste() *jerr.JAPIError {
	core.directivesWithPastes = make([]*directive.Directive, 0, 200)
	core.currentContextDirective = nil
	return core.processDirectiveList(core.directives)
}

func (core *JApiCore) processDirectiveList(list []*directive.Directive) *jerr.JAPIError {
	for i := 0; i != len(list); i++ {
		if je := core.processDirective(list[i]); je != nil {
			return je
		}
	}
	return nil
}

func (core *JApiCore) processDirective(d *directive.Directive) *jerr.JAPIError {
	if d.Type() == directive.Paste {
		if je := core.processPasteDirective(d); je != nil {
			return d.KeywordError(je.Error())
		}
		return nil
	}

	dd := d.CopyWoParentAndChildren()
	if je := core.processContext(&dd, &core.directivesWithPastes); je != nil {
		return je
	}

	if d.Children != nil {
		if je := core.processDirectiveList(d.Children); je != nil {
			return je
		}
	}

	if d.HasExplicitContext {
		core.currentContextDirective = dd.Parent
	}

	return nil
}

func (core *JApiCore) processPasteDirective(paste *directive.Directive) *jerr.JAPIError {
	if paste.Annotation != "" {
		return paste.KeywordError(jerr.AnnotationIsForbiddenForTheDirective)
	}
	if paste.Parameter("Name") == "" {
		return paste.KeywordError(fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Name"))
	}

	macro, ok := core.macro[paste.Parameter("Name")]
	if !ok {
		return paste.KeywordError("macro not found")
	}

	// macro.Children != nil - checked above
	return core.processDirectiveList(macro.Children)
}
