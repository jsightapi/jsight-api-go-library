package catalog

type HttpInteraction struct {
	PathVal       Path           `json:"path"`
	Tags          []TagName      `json:"tags"`
	PathVariables *PathVariables `json:"pathVariables,omitempty"`
	Annotation    *string        `json:"annotation,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Query         *Query         `json:"query,omitempty"`
	Request       *HTTPRequest   `json:"request,omitempty"`
	Responses     []HTTPResponse `json:"responses,omitempty"`
	HttpMethod    HttpMethod     `json:"httpMethod"`
}

func (h HttpInteraction) Path() Path {
	return h.PathVal
}

func (h *HttpInteraction) SetPathVariables(p *PathVariables) {
	h.PathVariables = p
}

func newHttpInteraction(path Path, method HttpMethod, annotation string, tn TagName) *HttpInteraction {
	h := &HttpInteraction{
		HttpMethod:    method,
		PathVal:       path,
		Tags:          []TagName{tn},
		PathVariables: nil,
		Annotation:    nil,
		Description:   nil,
		Query:         nil,
		Request:       nil,
		Responses:     []HTTPResponse{},
	}
	if annotation != "" {
		h.Annotation = &annotation
	}
	return h
}
