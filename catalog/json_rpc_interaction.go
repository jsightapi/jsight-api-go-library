package catalog

import "github.com/jsightapi/jsight-api-go-library/directive"

type JsonRpcInteraction struct {
	Description   *string
	Annotation    *string
	Params        *jsonRpcParams
	Result        *jsonRpcResult
	PathVariables *PathVariables
	Protocol      Protocol
	PathVal       Path
	Method        string
	Tags          []TagName
}

type jsonRpcParams struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

type jsonRpcResult struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

func (j JsonRpcInteraction) Path() Path {
	return j.PathVal
}

func (j *JsonRpcInteraction) SetPathVariables(p *PathVariables) {
	j.PathVariables = p
}

func newJsonRpcInteraction(path Path, method string, annotation string, tn TagName) *JsonRpcInteraction {
	j := &JsonRpcInteraction{
		Description:   nil,
		Annotation:    nil,
		Params:        nil,
		Result:        nil,
		PathVariables: nil,
		Protocol:      jsonRpc,
		PathVal:       path,
		Method:        method,
		Tags:          []TagName{tn},
	}
	if annotation != "" {
		j.Annotation = &annotation
	}
	return j
}
