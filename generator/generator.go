package generator

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/jsightapi/jsight-api-go-library/kit"
)

type Format string

const (
	FormatHTML Format = "html"
	FormatPDF  Format = "pdf"
	FormatDOCX Format = "docx"
)

// Generator an abstraction for generating documentation from the specification.
type Generator interface {
	// Generate generates documentation from the specification.
	// Specification will be read from in and print to out.
	Generate(ctx context.Context, in io.Reader, out io.Writer) error
}

var ErrUnsupportedFormat = errors.New("unsupported format")

func New(f Format) (Generator, error) {
	if f != FormatHTML {
		return nil, ErrUnsupportedFormat
	}
	return newHTML(), nil
}

type common struct {
	gen func(kit.JApi, io.Writer) error
}

func (c common) Generate(_ context.Context, in io.Reader, out io.Writer) error {
	b, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("failed to read specification: %w", err)
	}
	j := kit.NewJapiFromBytes(b)
	if err := j.ValidateJAPI(); err != nil {
		return fmt.Errorf("validate specification: %w", err)
	}
	return c.gen(j, out)
}
