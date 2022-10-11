package validator

import (
	"net/http"

	"github.com/jsightapi/jsight-api-go-library/catalog"
)

type HTTPRequestValidator struct {
	c *catalog.Catalog
	r *http.Request
}

func NewHTTPRequestValidator(c *catalog.Catalog, r *http.Request) HTTPRequestValidator {
	return HTTPRequestValidator{
		c: c,
		r: r,
	}
}

func (v HTTPRequestValidator) Process() error {
	_, err := v.possibleInteractions()
	return err
}

func (v HTTPRequestValidator) possibleInteractions() ([]catalog.Interaction, error) {
	found := make([]catalog.Interaction, 0, 5)

	method, err := catalog.NewHTTPMethodFromString(v.r.Method)
	if err != nil {
		return found, err
	}

	err = v.c.Interactions.Each(func(id catalog.InteractionID, interaction catalog.Interaction) error {
		if id.Protocol() == catalog.HTTP {
			kk := id.(catalog.HTTPInteractionID)
			if kk.Method() == method && id.Path().Similar(v.r.URL.Path) {
				found = append(found, interaction)
			}
		}
		return nil
	})

	return found, err
}
