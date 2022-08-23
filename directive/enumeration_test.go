package directive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumeration_String(t *testing.T) {
	cc := map[string]Enumeration{
		"JSIGHT":             Jsight,
		"INFO":               Info,
		"Title":              Title,
		"Version":            Version,
		"Description":        Description,
		"SERVER":             Server,
		"BaseUrl":            BaseURL,
		"URL":                URL,
		"GET":                Get,
		"POST":               Post,
		"PUT":                Put,
		"PATCH":              Patch,
		"DELETE":             Delete,
		"Body":               Body,
		"Request":            Request,
		"HTTP-response-code": HTTPResponseCode,
		"Path":               Path,
		"Headers":            Headers,
		"Query":              Query,
		"TYPE":               Type,
		"ENUM":               Enum,
		"MACRO":              Macro,
		"PASTE":              Paste,
		"INCLUDE":            Include,
		"Protocol":           Protocol,
		"Method":             Method,
		"Params":             Params,
		"Result":             Result,
	}

	for expected, given := range cc {
		t.Run(expected, func(t *testing.T) {
			actual := given.String()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestEnumeration_IsHTTPRequestMethod(t *testing.T) {
	cc := map[Enumeration]bool{
		Jsight:           false,
		Info:             false,
		Title:            false,
		Version:          false,
		Description:      false,
		Server:           false,
		BaseURL:          false,
		URL:              false,
		Get:              true,
		Post:             true,
		Put:              true,
		Patch:            true,
		Delete:           true,
		Body:             false,
		Request:          false,
		HTTPResponseCode: false,
		Path:             false,
		Headers:          false,
		Query:            false,
		Type:             false,
		Enum:             false,
		Macro:            false,
		Paste:            false,
		Include:          false,
		Protocol:         false,
		Method:           false,
		Params:           false,
		Result:           false,
	}

	for e, expected := range cc {
		t.Run(e.String(), func(t *testing.T) {
			actual := e.IsHTTPRequestMethod()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestEnumeration_IsAllowedForRootContext(t *testing.T) {
	cc := map[Enumeration]bool{
		Jsight:           true,
		Info:             true,
		Title:            false,
		Version:          false,
		Description:      false,
		Server:           true,
		BaseURL:          false,
		URL:              true,
		Get:              true,
		Post:             true,
		Put:              true,
		Patch:            true,
		Delete:           true,
		Body:             false,
		Request:          false,
		HTTPResponseCode: false,
		Path:             false,
		Headers:          false,
		Query:            false,
		Type:             true,
		Enum:             true,
		Macro:            true,
		Paste:            true,
		Include:          false,
		Protocol:         false,
		Method:           false,
		Params:           false,
		Result:           false,
	}

	for e, expected := range cc {
		t.Run(e.String(), func(t *testing.T) {
			actual := e.IsAllowedForRootContext()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestEnumeration_IsAllowedForDirectiveContext(t *testing.T) {
	type testCase struct {
		base     Enumeration
		child    Enumeration
		expected bool
	}

	cc := []testCase{
		{URL, Jsight, false},
		{URL, Info, false},
		{URL, Title, false},
		{URL, Version, false},
		{URL, Description, false},
		{URL, Server, false},
		{URL, BaseURL, false},
		{URL, URL, false},
		{URL, Get, true},
		{URL, Post, true},
		{URL, Put, true},
		{URL, Patch, true},
		{URL, Delete, true},
		{URL, Body, false},
		{URL, Request, false},
		{URL, HTTPResponseCode, false},
		{URL, Path, true},
		{URL, Headers, false},
		{URL, Query, false},
		{URL, Type, false},
		{URL, Enum, false},
		{URL, Macro, false},
		{URL, Paste, true},
		{URL, Include, false},
		{URL, Protocol, true},
		{URL, Method, true},
		{URL, Params, false},
		{URL, Result, false},

		{Info, Jsight, false},
		{Info, Info, false},
		{Info, Title, true},
		{Info, Version, true},
		{Info, Description, true},
		{Info, Server, false},
		{Info, BaseURL, false},
		{Info, URL, false},
		{Info, Get, false},
		{Info, Post, false},
		{Info, Put, false},
		{Info, Patch, false},
		{Info, Delete, false},
		{Info, Body, false},
		{Info, Request, false},
		{Info, HTTPResponseCode, false},
		{Info, Path, false},
		{Info, Headers, false},
		{Info, Query, false},
		{Info, Type, false},
		{Info, Enum, false},
		{Info, Macro, false},
		{Info, Paste, true},
		{Info, Include, false},
		{Info, Protocol, false},
		{Info, Method, false},
		{Info, Params, false},
		{Info, Result, false},

		{Server, Jsight, false},
		{Server, Info, false},
		{Server, Title, false},
		{Server, Version, false},
		{Server, Description, false},
		{Server, Server, false},
		{Server, BaseURL, true},
		{Server, URL, false},
		{Server, Get, false},
		{Server, Post, false},
		{Server, Put, false},
		{Server, Patch, false},
		{Server, Delete, false},
		{Server, Body, false},
		{Server, Request, false},
		{Server, HTTPResponseCode, false},
		{Server, Path, false},
		{Server, Headers, false},
		{Server, Query, false},
		{Server, Type, false},
		{Server, Enum, false},
		{Server, Macro, false},
		{Server, Paste, true},
		{Server, Include, false},
		{Server, Protocol, false},
		{Server, Method, false},
		{Server, Params, false},
		{Server, Result, false},

		{Macro, Jsight, false},
		{Macro, Info, true},
		{Macro, Title, true},
		{Macro, Version, true},
		{Macro, Description, true},
		{Macro, Server, true},
		{Macro, BaseURL, true},
		{Macro, URL, true},
		{Macro, Get, true},
		{Macro, Post, true},
		{Macro, Put, true},
		{Macro, Patch, true},
		{Macro, Delete, true},
		{Macro, Body, true},
		{Macro, Request, true},
		{Macro, HTTPResponseCode, true},
		{Macro, Path, true},
		{Macro, Headers, true},
		{Macro, Query, true},
		{Macro, Type, true},
		{Macro, Enum, true},
		{Macro, Macro, false},
		{Macro, Paste, true},
		{Macro, Include, false},
		{Macro, Protocol, false},
		{Macro, Method, false},
		{Macro, Params, false},
		{Macro, Result, false},
	}

	all := []Enumeration{
		Jsight,
		Info,
		Title,
		Version,
		Description,
		Server,
		BaseURL,
		URL,
		Get,
		Post,
		Put,
		Patch,
		Delete,
		Body,
		Request,
		HTTPResponseCode,
		Path,
		Headers,
		Query,
		Type,
		Enum,
		Macro,
		Paste,
		Include,
		Protocol,
		Method,
		Params,
		Result,
	}

	notAllowedAtAll := []Enumeration{
		Jsight,
		Title,
		Version,
		Description,
		BaseURL,
		Body,
		Path,
		Headers,
		Query,
		Type,
		Enum,
		Paste,
		Include,
	}
	for _, base := range notAllowedAtAll {
		for _, child := range all {
			cc = append(cc, testCase{base, child, false})
		}
	}

	crud := []Enumeration{Get, Post, Put, Patch, Delete}
	for _, base := range crud {
		for _, child := range all {
			expected := false
			switch child { //nolint:exhaustive // False-positive.
			case Description, Query, Path, Request, HTTPResponseCode, Paste:
				expected = true
			}
			cc = append(cc, testCase{base, child, expected})
		}
	}

	responseAndRequests := []Enumeration{
		HTTPResponseCode,
		Request,
	}
	for _, base := range responseAndRequests {
		for _, child := range all {
			expected := false
			switch child { //nolint:exhaustive // False-positive.
			case Body, Headers, Paste:
				expected = true
			}
			cc = append(cc, testCase{base, child, expected})
		}
	}

	for _, c := range cc {
		t.Run(fmt.Sprintf("%s <- %s = %t", c.base, c.child, c.expected), func(t *testing.T) {
			actual := c.base.IsAllowedForDirectiveContext(c.child)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func BenchmarkEnumeration_IsAllowedForDirectiveContext(b *testing.B) {
	b.ReportAllocs()

	b.Run("first match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			URL.IsAllowedForDirectiveContext(Path)
		}
	})

	b.Run("last match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Macro.IsAllowedForDirectiveContext(Paste)
		}
	})

	b.Run("miss", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Paste.IsAllowedForDirectiveContext(Include)
		}
	})
}
