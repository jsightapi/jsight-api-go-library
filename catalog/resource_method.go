package catalog

type ResourceMethod struct {
	HttpMethod    Method         `json:"httpMethod"`
	Path          Path           `json:"path"`
	Tags          []TagName      `json:"tags"`
	PathVariables *PathVariables `json:"pathVariables,omitempty"`
	Annotation    string         `json:"annotation,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Query         *Query         `json:"query,omitempty"`
	Request       *HTTPRequest   `json:"request,omitempty"`
	Responses     []HTTPResponse `json:"responses,omitempty"`
}

func initResourceMethod(path Path, method Method, annotation string, tn TagName) ResourceMethod {
	return ResourceMethod{
		HttpMethod:    method,
		Path:          path,
		Tags:          []TagName{tn},
		PathVariables: nil,
		Annotation:    annotation,
		Description:   nil,
		Query:         nil,
		Request:       nil,
		Responses:     []HTTPResponse{},
	}
}
