package catalog

import (
	"strings"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/notations/regex"

	"github.com/jsightapi/jsight-api-go-library/notation"
)

func PrepareRegexSchema(name string, regexStr bytes.Bytes) (*regex.Schema, error) {
	var oo []regex.Option
	oo = append(oo, regex.WithGeneratorSeed(0))

	s := regex.New(name, regexStr, oo...)

	// n, err := s.GetAST()
	// if err != nil {
	// 	return Schema{}, err
	// }

	return s, nil
}

// TODO remove the following ???

func UnmarshalRegexSchema(name string, regexStr bytes.Bytes) (schema Schema, err error) {
	return mainRegexMarshaller.Marshal(name, regexStr)
}

type regexMarshaller struct {
	useFixedSeed bool
}

var mainRegexMarshaller = regexMarshaller{}

func (m regexMarshaller) Marshal(name string, regexStr bytes.Bytes) (schema Schema, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	var oo []regex.Option
	if m.useFixedSeed {
		oo = append(oo, regex.WithGeneratorSeed(0))
	}

	s := regex.New(name, regexStr, oo...)

	n, err := s.GetAST()
	if err != nil {
		return Schema{}, err
	}

	schema = NewSchema(notation.SchemaNotationRegex)
	schema.ContentRegexp = strings.TrimPrefix(n.Value, "/")
	schema.ContentRegexp = strings.TrimSuffix(schema.ContentRegexp, "/")
	schema.JSchema = s

	example, err := s.Example()
	if err != nil {
		return Schema{}, err
	}

	schema.Example = string(example)
	return schema, nil
}
