package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/notation"
)

func (core *JApiCore) validateCatalog() *jerr.JApiError {
	if je := core.validateInfo(); je != nil {
		return je
	}

	if je := core.validateRequestBody(); je != nil {
		return je
	}

	if je := core.validateResponseBody(); je != nil {
		return je
	}

	if je := core.validateHeaders(); je != nil {
		return je
	}

	return core.validateUsedUserTypes()
}

func (core *JApiCore) validateInfo() *jerr.JApiError {
	isEmpty := core.catalog.Info != nil &&
		core.catalog.Info.Title == "" &&
		core.catalog.Info.Version == "" &&
		core.catalog.Info.Description == nil

	if isEmpty {
		return core.catalog.Info.Directive.KeywordError("empty info")
	}
	return nil
}

func (core *JApiCore) validateUsedUserTypes() *jerr.JApiError {
	err := core.catalog.UserTypes.Each(func(k string, v *catalog.UserType) error {
		if err := core.findUserTypes(v.Schema.UsedUserTypes); err != nil {
			return v.Directive.BodyError(err.Error())
		}
		return nil
	})
	if err != nil {
		return adoptError(err)
	}

	err = core.catalog.Servers.Each(func(k string, v *catalog.Server) error {
		s := v.BaseUrlVariables
		if s != nil && s.Schema != nil {
			if err := core.findUserTypes(s.Schema.UsedUserTypes); err != nil {
				return s.Directive.BodyError(err.Error())
			}
		}
		return nil
	})
	if err != nil {
		return adoptError(err)
	}

	return adoptError(core.catalog.Interactions.Each(func(k catalog.InteractionID, v catalog.Interaction) error {
		hi, ok := v.(*catalog.HTTPInteraction)
		if !ok {
			return nil
		}
		if hi.Query != nil && hi.Query.Schema != nil {
			if err := core.findUserTypes(hi.Query.Schema.UsedUserTypes); err != nil {
				return hi.Query.Directive.BodyError(err.Error())
			}
		}

		if hi.Request != nil {
			if hi.Request.HTTPRequestHeaders != nil && hi.Request.HTTPRequestHeaders.Schema != nil {
				if err := core.findUserTypes(hi.Request.HTTPRequestHeaders.Schema.UsedUserTypes); err != nil {
					return hi.Request.HTTPRequestHeaders.Directive.BodyError(err.Error())
				}
			}

			if hi.Request.HTTPRequestBody != nil && hi.Request.HTTPRequestBody.Schema != nil {
				if err := core.findUserTypes(hi.Request.HTTPRequestBody.Schema.UsedUserTypes); err != nil {
					return hi.Request.HTTPRequestBody.Directive.BodyError(err.Error())
				}
			}
		}

		for _, r := range hi.Responses {
			if r.Headers != nil && r.Headers.Schema != nil {
				if err := core.findUserTypes(r.Headers.Schema.UsedUserTypes); err != nil {
					return r.Headers.Directive.BodyError(err.Error())
				}
			}

			if r.Body != nil && r.Body.Schema != nil {
				if err := core.findUserTypes(r.Body.Schema.UsedUserTypes); err != nil {
					return r.Body.Directive.BodyError(err.Error())
				}
			}
		}
		return nil
	}))
}

// findUserTypes returns an error if a user type cannot be found
func (core *JApiCore) findUserTypes(uu *catalog.StringSet) error {
	for _, u := range uu.Data() {
		if !core.catalog.UserTypes.Has(u) {
			return fmt.Errorf("user type %q wasn't found", u)
		}
	}
	return nil
}

func (core *JApiCore) validateRequestBody() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(k catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			r := hi.Request
			if r != nil && r.HTTPRequestBody == nil {
				return r.Directive.KeywordError(fmt.Sprintf(`undefined request body for resource "%s"`, k.String()))
			}
		}
		return nil
	}))
}

func (core *JApiCore) validateResponseBody() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(k catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			for _, response := range hi.Responses {
				if response.Body == nil {
					return response.Directive.KeywordError(fmt.Sprintf(
						"undefined response body for resource %q, HTTP-code %q",
						k.String(),
						response.Code,
					))
				}
			}
		}
		return nil
	}))
}

func (core *JApiCore) isJsightCastToObject(schema *catalog.Schema) bool {
	if schema != nil && schema.ContentJSight != nil && schema.Notation == notation.SchemaNotationJSight {
		switch schema.ContentJSight.TokenType {
		case "object":
			return true
		case "reference":
			if userType, ok := core.catalog.UserTypes.Get(schema.ContentJSight.ScalarValue); ok {
				return core.isJsightCastToObject(&userType.Schema)
			}
		}
	}
	return false
}

func (core *JApiCore) validateHeaders() *jerr.JApiError {
	return adoptError(core.catalog.Interactions.Each(func(_ catalog.InteractionID, v catalog.Interaction) error {
		if hi, ok := v.(*catalog.HTTPInteraction); ok {
			request := hi.Request
			isNotAnObject := request != nil &&
				request.HTTPRequestHeaders != nil &&
				!core.isJsightCastToObject(request.HTTPRequestHeaders.Schema)
			if isNotAnObject {
				return request.HTTPRequestHeaders.Directive.BodyError(jerr.BodyMustBeObject)
			}
			for _, response := range hi.Responses {
				if response.Headers != nil && !core.isJsightCastToObject(response.Headers.Schema) {
					return response.Headers.Directive.BodyError(jerr.BodyMustBeObject)
				}
			}
		}
		return nil
	}))
}
