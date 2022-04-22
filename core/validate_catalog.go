package core

import (
	"fmt"
	"j/japi/catalog"
	"j/japi/jerr"
	"j/japi/notation"
)

func (core *JApiCore) validateCatalog() *jerr.JAPIError {
	if err := core.validateInfo(); err != nil {
		return err
	}

	if err := core.validateRequestBody(); err != nil {
		return err
	}

	if err := core.validateResponseBody(); err != nil {
		return err
	}

	if err := core.validateHeaders(); err != nil {
		return err
	}

	// if err := core.validatePathVariables(); err != nil {
	// 	return err
	// }

	return core.validateUsedUserTypes()
}

func (core *JApiCore) validateInfo() *jerr.JAPIError {
	if core.catalog.Info != nil && core.catalog.Info.Title == "" && core.catalog.Info.Version == "" && core.catalog.Info.Description == nil {
		return core.catalog.Info.Directive.KeywordError("empty info")
	}
	return nil
}

func (core *JApiCore) validateUsedUserTypes() *jerr.JAPIError {
	for ku := range core.catalog.UserTypes.Iterate() {
		u := ku.Value
		if err := core.findUserTypes(u.Schema.UsedUserTypes); err != nil {
			return u.Directive.BodyError(err.Error())
		}
	}

	for ks := range core.catalog.Servers.Iterate() {
		s := ks.Value
		if s.BaseUrlVariables != nil && s.BaseUrlVariables.Schema != nil {
			if err := core.findUserTypes(s.BaseUrlVariables.Schema.UsedUserTypes); err != nil {
				return s.BaseUrlVariables.Directive.BodyError(err.Error())
			}
		}
	}

	for kv := range core.catalog.ResourceMethods.Iterate() {
		v := kv.Value

		if v.Query != nil && v.Query.Schema != nil {
			if err := core.findUserTypes(v.Query.Schema.UsedUserTypes); err != nil {
				return v.Query.Directive.BodyError(err.Error())
			}
		}

		if v.Request != nil {
			if v.Request.HTTPRequestHeaders != nil && v.Request.HTTPRequestHeaders.Schema != nil {
				if err := core.findUserTypes(v.Request.HTTPRequestHeaders.Schema.UsedUserTypes); err != nil {
					return v.Request.HTTPRequestHeaders.Directive.BodyError(err.Error())
				}
			}

			if v.Request.HTTPRequestBody != nil && v.Request.HTTPRequestBody.Schema != nil {
				if err := core.findUserTypes(v.Request.HTTPRequestBody.Schema.UsedUserTypes); err != nil {
					return v.Request.HTTPRequestBody.Directive.BodyError(err.Error())
				}
			}
		}

		for _, r := range v.Responses {
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
	}
	return nil
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

func (core *JApiCore) validateRequestBody() *jerr.JAPIError {
	for kv := range core.catalog.ResourceMethods.Iterate() {
		r := kv.Value
		if r.Request != nil && r.Request.HTTPRequestBody == nil {
			return r.Request.Directive.KeywordError(fmt.Sprintf(`undefined request body for resource "%s"`, kv.Key.String()))
		}
	}
	return nil
}

func (core *JApiCore) validateResponseBody() *jerr.JAPIError {
	for kv := range core.catalog.ResourceMethods.Iterate() {
		for _, response := range kv.Value.Responses {
			if response.Body == nil {
				return response.Directive.KeywordError(fmt.Sprintf(`undefined response body for resource "%s", HTTP-code "%s"`, kv.Key.String(), response.Code))
			}
		}
	}
	return nil
}

func (core *JApiCore) isJsightCastToObject(schema *catalog.Schema) bool {
	if schema != nil && schema.ContentJSight != nil && schema.Notation == notation.SchemaNotationJSight {
		switch schema.ContentJSight.JsonType {
		case "object":
			return true
		case "shortcut":
			if userType, ok := core.catalog.UserTypes.Get(schema.ContentJSight.ScalarValue); ok {
				return core.isJsightCastToObject(&userType.Schema)
			}
		}
	}
	return false
}

func (core *JApiCore) validateHeaders() *jerr.JAPIError {
	for kv := range core.catalog.ResourceMethods.Iterate() {
		resourceMethod := kv.Value

		request := resourceMethod.Request
		if request != nil && request.HTTPRequestHeaders != nil && !core.isJsightCastToObject(request.HTTPRequestHeaders.Schema) {
			return resourceMethod.Request.HTTPRequestHeaders.Directive.BodyError(jerr.BodyMustBeObject)
		}
		for _, response := range resourceMethod.Responses {
			if response.Headers != nil && !core.isJsightCastToObject(response.Headers.Schema) {
				return response.Headers.Directive.BodyError(jerr.BodyMustBeObject)
			}
		}
	}
	return nil
}

// func (core *JApiCore) validatePathVariables() *jerr.JAPIError {
// 	for kv := range core.catalog.ResourceMethods.Iterate() {
// 		kv.Value.PathVariables
// 	}
// 	return nil
// }
