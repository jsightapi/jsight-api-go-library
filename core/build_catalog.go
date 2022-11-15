package core

import (
	stdErrors "errors"
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	schema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/kit"
)

func (core *JApiCore) buildCatalog() *jerr.JApiError {
	if len(core.directivesWithPastes) != 0 && core.directivesWithPastes[0].Type() != directive.Jsight {
		return core.directivesWithPastes[0].KeywordError("JSIGHT should be the first directive")
	}

	return core.addDirectives()
}

func adoptError(err error) (e *jerr.JApiError) {
	if err == nil {
		return nil
	}

	if stdErrors.As(err, &e) {
		return e
	}

	panic(fmt.Sprintf("Invalid error was given: %#v", err))
}

func safeAddType(curr schema.Schema, n string, ut schema.Schema) error {
	err := curr.AddType(n, ut)
	var e interface{ Code() errors.ErrorCode }
	if stdErrors.As(err, &e) && e.Code() == errors.ErrDuplicationOfNameOfTypes {
		err = nil
	}
	return err
}

func (core *JApiCore) checkUserType(name string) *jerr.JApiError {
	err := core.userTypes.GetValue(name).Check()
	if err == nil {
		return nil
	}

	d := core.rawUserTypes.GetValue(name)
	var e kit.Error
	if !stdErrors.As(err, &e) {
		return d.KeywordError(err.Error())
	}

	if e.IncorrectUserType() != "" && e.IncorrectUserType() != name {
		return core.checkUserType(e.IncorrectUserType())
	}

	return d.BodyErrorIndex(e.Message(), e.Position())
}

func jschemaToJAPIError(err error, d *directive.Directive) *jerr.JApiError {
	var e kit.Error
	if stdErrors.As(err, &e) {
		return d.BodyErrorIndex(e.Message(), e.Position())
	}
	return d.KeywordError(err.Error())
}
