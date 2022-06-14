package catalog

import "github.com/jsightapi/jsight-api-go-library/directive"

type JsonRpcInteraction struct {
	Protocol      Protocol
	Path          Path
	Method        string
	PathVariables *PathVariables
	Tags          []TagName
	Annotation    *string
	Description   *string
	Params        *jsonRpcParams
	Result        *jsonRpcResult
}

func (ji JsonRpcInteraction) path() Path {
	return ji.Path
}

type jsonRpcParams struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

type jsonRpcResult struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}
