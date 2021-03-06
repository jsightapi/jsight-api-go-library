package catalog

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
)

type HttpInteractionId struct {
	protocol Protocol
	path     Path
	method   HttpMethod
}

func (h HttpInteractionId) Protocol() Protocol {
	return h.protocol
}

func (h HttpInteractionId) Path() Path {
	return h.path
}

func (h HttpInteractionId) String() string {
	return fmt.Sprintf("http %s %s", h.method.String(), h.path.String())
}

func (h HttpInteractionId) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func newHttpInteractionId(d directive.Directive) (HttpInteractionId, error) {
	h := HttpInteractionId{
		protocol: HTTP,
	}

	path, err := d.Path()
	if err != nil {
		return h, err
	}

	de, err := d.HttpMethod()
	if err != nil {
		return h, err
	}

	method, err := NewHttpMethod(de)
	if err != nil {
		return h, err
	}

	h.path = Path(path)
	h.method = method

	return h, nil
}
