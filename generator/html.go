package generator

import (
	"errors"
	"io"

	"github.com/jsightapi/jsight-api-go-library/kit"
)

type html struct {
	common
}

var _ Generator = html{}

func newHTML() html {
	h := html{}
	h.gen = h.generate
	return h
}

func (html) generate(_ kit.JApi, _ io.Writer) error {
	return errors.New("not implemented yet")
}
