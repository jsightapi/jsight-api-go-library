package catalog

import (
	"sync"

	"github.com/jsightapi/jsight-api-go-library/directive"
)

type Server struct {
	BaseUrlVariables *baseUrlVariables `json:"baseUrlVariables,omitempty"`
	Annotation       string            `json:"annotation,omitempty"`
	BaseUrl          string            `json:"baseUrl"`
}

type baseUrlVariables struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

// Servers represent available servers.
// gen:OrderedMap
type Servers struct {
	data  map[string]*Server
	order []string
	mx    sync.RWMutex
}
