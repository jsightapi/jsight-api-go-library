package scanner

import (
	"j/japi/jerr"
	"j/schema/bytes"
	"j/schema/fs"
	"j/schema/kit"
)

// Pass rest of the file to jsc scanner to find out where jschema ends
func stateJSchema(s *Scanner, _ byte) *jerr.JAPIError {
	s.found(SchemaBegin)
	schemaLength, je := s.readSchemaWithJsc()
	if je != nil {
		return je
	}
	if schemaLength > 0 {
		s.curIndex += bytes.Index(schemaLength - 1)
	}
	s.step = stateSchemaClosed
	return nil
}

func (s *Scanner) readSchemaWithJsc() (uint, *jerr.JAPIError) {
	b := s.file.Content()
	bb := b.Slice(s.curIndex, bytes.Index(b.Len()-1))
	f := fs.NewFile("", bb)
	schemaLength, err := kit.LengthOfSchema(f)
	if err != nil {
		return 0, s.japiError(err.Message(), s.curIndex+bytes.Index(err.Position()))
	}
	return schemaLength, nil
}
