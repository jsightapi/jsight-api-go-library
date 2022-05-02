package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/jsightapi/jsight-api-go-library/generator"
)

func main() {
	app := &cli.App{
		Name:  "jsight",
		Usage: "is a tool for working with files in the JSight language",
		Commands: []*cli.Command{
			{
				Name:  "doc",
				Usage: "generate API documentation in various formats",
				Subcommands: []*cli.Command{
					{
						Name:      "html",
						Usage:     "generates documentation in HTML format",
						Flags:     []cli.Flag{},
						ArgsUsage: "<input>",
						Action:    generateDoc(generator.FormatHTML),
					},
					{
						Name:      "pdf",
						Usage:     "generates documentation in PDF format",
						ArgsUsage: "<input>",
						Action:    generateDoc(generator.FormatPDF),
					},
					{
						Name:      "docx",
						Usage:     "generates documentation in DOCX format",
						ArgsUsage: "<input>",
						Action:    generateDoc(generator.FormatDOCX),
					},
				},
			},

			{
				Name:  "convert",
				Usage: "generate API documentation in various formats",
				Description: `If the path defined in <input> contains a document in the JSight language, then a file in the Open API format will be written to stdout.

If the path defined in the <input> contains a document in the OpenAPI language, then a file in the JSight format will be written to the stdout.`,
				ArgsUsage: "<input>",
			},

			{
				Name:  "stub",
				Usage: "create stubs in popular languages",
				Description: `The path to the file containing the API in the JSight language described by input argument.

The result of the generation is output to stdout.`,
				ArgsUsage: "<generator> <input>",
			},
		},
		EnableBashCompletion: true,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func generateDoc(f generator.Format) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		aa := ctx.Args()
		if aa.Len() != 1 {
			return cli.ShowCommandHelp(ctx, string(f))
		}

		specPath := aa.First()
		if _, err := os.Stat(specPath); err != nil {
			return err
		}

		r, err := os.Open(specPath)
		if err != nil {
			return err
		}

		g, err := generator.New(f)
		if err != nil {
			return err
		}

		return g.Generate(ctx.Context, r, ctx.App.Writer)
	}
}
