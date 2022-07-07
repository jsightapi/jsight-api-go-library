package catalog

import "github.com/jsightapi/jsight-api-go-library/directive"

type JsonRpcInteraction struct {
	Id            string         `json:"id"`
	Protocol      Protocol       `json:"protocol"`
	PathVal       Path           `json:"path"`
	Method        string         `json:"method"`
	PathVariables *PathVariables `json:"pathVariables"`
	Tags          []TagName      `json:"tags"`
	Description   *string        `json:"annotation"`
	Annotation    *string        `json:"description"`
	Params        *jsonRpcParams `json:"params"`
	Result        *jsonRpcResult `json:"result"`
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

func newJsonRpcInteraction(id JsonRpcInteractionId, method string, annotation string, tn TagName) *JsonRpcInteraction {
	j := &JsonRpcInteraction{
		Id:            id.String(),
		Protocol:      JsonRpc,
		PathVal:       id.path,
		PathVariables: nil,
		Method:        method,
		Tags:          []TagName{tn},
		Annotation:    nil,
		Description:   nil,
		Params:        nil,
		Result:        nil,
	}
	if annotation != "" {
		j.Annotation = &annotation
	}
	return j
}
