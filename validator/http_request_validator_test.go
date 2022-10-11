package validator

import (
	"bytes"
	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/core"
	"github.com/jsightapi/jsight-api-go-library/test"
	"github.com/jsightapi/jsight-schema-go-library/reader"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"
)

func NewCatalog(t *testing.T, jstFilename string) *catalog.Catalog {
	f := reader.Read(jstFilename) // can panic
	c := core.NewJApiCore(f)
	je := c.BuildCatalog()
	if je != nil {
		t.Fatal(je.Error())
	}
	return c.Catalog()
}

func NewRequest(t *testing.T, method string, url string, body []byte) *http.Request {
	r, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err.Error())
	}
	return r
}

func TestHTTPRequestValidator_Process(t *testing.T) {
	type fields struct {
		c *catalog.Catalog
		r *http.Request
	}

	jstFilename := filepath.Join(test.GetTestDataDir(), "jsight_0.3", "others", "full.jst")
	c := NewCatalog(t, jstFilename)

	tests := []struct {
		name   string
		fields fields
		want   catalog.InteractionID
	}{
		{
			"Test #1",
			fields{
				c: c,
				r: NewRequest(t, "GET", "/users/123", []byte{}),
			},
			catalog.HTTPInteractionID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := HTTPRequestValidator{
				c: tt.fields.c,
				r: tt.fields.r,
			}
			if got := v.Process(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
