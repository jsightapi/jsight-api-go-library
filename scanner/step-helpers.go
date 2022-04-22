package scanner

func caseWhitespace(c byte) byte {
	if isWhitespace(c) {
		return c
	} else {
		return otherByte(c)
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t'
}

func caseNewLine(c byte) byte {
	if isNewLine(c) {
		return c
	} else {
		return otherByte(c)
	}
}

func isNewLine(c byte) bool {
	return c == '\n' || c == '\r'
}

// byte is uint8: 0-255. 0 is for EOF.
func otherByte(b byte) byte {
	if b == 255 {
		return 254
	} else {
		return b + 1
	}
}

var any = []byte("any")
var empty = []byte("empty")
var regex = []byte("regex")

func (s *Scanner) isDirectiveParameterHasTypeOrAnyOrEmpty() bool {
	for _, lex := range s.lastDirectiveParameters {
		v := lex.Value().Unquote().TrimSquareBrackets()
		switch {
		case v.Equals(any), v.Equals(empty), v.IsUserTypeName():
			return true
		}
	}
	return false
}

func (s *Scanner) isDirectiveParameterHasAnyOrEmpty() bool {
	for _, lex := range s.lastDirectiveParameters {
		v := lex.Value().Unquote().TrimSquareBrackets()
		switch {
		case v.Equals(any), v.Equals(empty):
			return false
		}
	}
	return true
}

func (s *Scanner) isDirectiveParameterHasRegexNotation() bool {
	for _, lex := range s.lastDirectiveParameters {
		v := lex.Value().Unquote()
		if v.Equals(regex) {
			return true
		}
	}
	return false
}
