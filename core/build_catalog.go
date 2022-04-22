package core

import (
	"errors"
	"j/japi/directive"
	"j/japi/jerr"
	"j/japi/notation"
	jschemaLib "j/schema"
	jerrors "j/schema/errors"
	"j/schema/kit"
	"j/schema/notations/jschema"
	"j/schema/notations/regex"
)

func (core *JApiCore) buildCatalog() *jerr.JAPIError {
	if len(core.directivesWithPastes) != 0 && core.directivesWithPastes[0].Type() != directive.Jsight {
		return core.directivesWithPastes[0].KeywordError("JSIGHT should be the first directive")
	}

	return core.addDirectives()
}

func (core *JApiCore) compileUserTypes() *jerr.JAPIError {
	// Two-phase algorithm. On the first step we just create schema for each user
	// type. On the second step we will add all schema to all.
	// This is the simplest solution which allows us to skip building dependency
	// graph between types.

	if err := core.buildUserTypes(); err != nil {
		return err
	}

	for kv := range core.userTypes.Iterate() {
		if err := core.chekUserType(kv.Key); err != nil {
			return err
		}
	}
	return nil
}

func (core *JApiCore) buildUserTypes() *jerr.JAPIError {
	for kv := range core.catalog.GetRawUserTypes().Iterate() {
		switch notation.SchemaNotation(kv.Value.Parameter("SchemaNotation")) {
		case "", notation.SchemaNotationJSight:
			if kv.Value.BodyCoords.IsSet() {
				core.userTypes.Set(kv.Key, jschema.New(kv.Key, kv.Value.BodyCoords.Read()))
			}
		case notation.SchemaNotationRegex:
			core.userTypes.Set(kv.Key, regex.New(kv.Key, kv.Value.BodyCoords.Read()))
		default:
			// nothing
		}
	}

	for kv := range core.userTypes.Iterate() {
		if err := core.buildUserType(kv.Key); err != nil {
			return err
		}
	}
	return nil
}

func (core *JApiCore) buildUserType(name string) *jerr.JAPIError {
	if _, ok := core.processedUserTypes[name]; ok {
		// This user type already processed, skip.
		return nil
	}

	currUT := core.userTypes.GetValue(name)
	if currUT == nil {
		return nil
	}

	dd := core.catalog.GetRawUserTypes()

	tt, err := core.getUsedUserTypes(currUT)
	if err != nil {
		return jschemaToJAPIError(err, dd.GetValue(name))
	}

	core.processedUserTypes[name] = struct{}{}
	alreadyAddedTypes := map[string]struct{}{}

	for _, n := range tt {
		if n != name {
			if err := core.buildUserType(n); err != nil {
				return err
			}
		}

		ut := core.userTypes.GetValue(n)
		if ut == nil {
			continue
		}

		if _, ok := alreadyAddedTypes[n]; !ok {
			if err := safeAddType(currUT, n, ut); err != nil {
				return jschemaToJAPIError(err, dd.GetValue(n))
			}
			alreadyAddedTypes[n] = struct{}{}
		}
	}

	core.userTypes.Set(name, currUT)

	return nil
}

func safeAddType(curr jschemaLib.Schema, n string, ut jschemaLib.Schema) error {
	err := curr.AddType(n, ut)
	var e interface{ Code() jerrors.ErrorCode }
	if errors.As(err, &e) && e.Code() == jerrors.ErrDuplicationOfNameOfTypes {
		err = nil
	}
	return err
}

func (core *JApiCore) getUsedUserTypes(ut jschemaLib.Schema) ([]string, error) {
	alreadyProcessed := map[string]struct{}{}
	if err := core.fetchUsedUserTypes(ut, alreadyProcessed); err != nil {
		return nil, err
	}

	ss := make([]string, 0, len(alreadyProcessed))
	for s := range alreadyProcessed {
		ss = append(ss, s)
	}

	return ss, nil
}

func (core *JApiCore) fetchUsedUserTypes(
	ut jschemaLib.Schema,
	alreadyProcessed map[string]struct{},
) error {
	if ut == nil {
		return nil
	}

	tt, err := ut.UsedUserTypes()
	if err != nil {
		return err
	}

	if len(tt) == 0 {
		return nil
	}

	for _, t := range tt {
		if _, ok := alreadyProcessed[t]; ok {
			continue
		}

		alreadyProcessed[t] = struct{}{}
		if err := core.fetchUsedUserTypes(core.userTypes.GetValue(t), alreadyProcessed); err != nil {
			return err
		}
	}
	return nil
}

func (core *JApiCore) chekUserType(name string) *jerr.JAPIError {
	err := core.userTypes.GetValue(name).Check()
	if err == nil {
		return nil
	}

	d := core.catalog.GetRawUserTypes().GetValue(name)
	var e kit.Error
	if !errors.As(err, &e) {
		return d.KeywordError(err.Error())
	}

	if e.IncorrectUserType() != "" && e.IncorrectUserType() != name {
		return core.chekUserType(e.IncorrectUserType())
	}

	return d.BodyErrorIndex(e.Message(), e.Position())
}

func jschemaToJAPIError(err error, d *directive.Directive) *jerr.JAPIError {
	var e kit.Error
	if errors.As(err, &e) {
		return d.BodyErrorIndex(e.Message(), e.Position())
	}
	return d.KeywordError(err.Error())
}
