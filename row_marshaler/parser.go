package row_marshaler

import (
	"strings"
	"unicode"
)

// Token represents a parsed value from a row literal.
type Token struct {
	Value    string // The parsed value (with escapes processed)
	RawValue string // The original raw value from input
	Quoted   bool   // Whether the value was quoted in the input
	Position int    // Position in the input where this token started
}

// ParseRowLiteral parses a row literal string into a slice of tokens.
// A row literal has the format: (value1,value2,"quoted value",...)
func ParseRowLiteral(input string) ([]Token, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, ErrEmptyInput
	}

	// Check for opening parenthesis
	if len(input) == 0 || input[0] != '(' {
		return nil, NewParseError(0, "expected opening parenthesis '('", input)
	}

	// Check for closing parenthesis
	if input[len(input)-1] != ')' {
		return nil, NewParseError(len(input)-1, "expected closing parenthesis ')'", input)
	}

	// Extract content between parentheses
	content := input[1 : len(input)-1]

	// Handle empty row
	if strings.TrimSpace(content) == "" {
		return []Token{}, nil
	}

	return parseTokens(content, 1) // Start at position 1 (after opening paren)
}

// parseTokens parses the content inside the parentheses into tokens.
func parseTokens(content string, startPos int) ([]Token, error) {
	var tokens []Token
	pos := 0
	n := len(content)

	for pos < n {
		// Skip leading whitespace
		for pos < n && unicode.IsSpace(rune(content[pos])) {
			pos++
		}

		if pos >= n {
			break
		}

		var token Token
		var err error

		if content[pos] == '"' {
			// Parse quoted string
			token, pos, err = parseQuotedString(content, pos, startPos)
		} else {
			// Parse unquoted value
			token, pos, err = parseUnquotedValue(content, pos, startPos)
		}

		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)

		// Skip whitespace after value
		for pos < n && unicode.IsSpace(rune(content[pos])) {
			pos++
		}

		// Check for comma separator or end
		if pos < n {
			if content[pos] == ',' {
				pos++
				// Allow trailing comma with no value after
				// by continuing the loop
			} else {
				return nil, NewParseError(startPos+pos, "expected ',' or end of row", content)
			}
		}
	}

	return tokens, nil
}

// parseQuotedString parses a quoted string starting at pos.
func parseQuotedString(content string, pos int, startPos int) (Token, int, error) {
	if pos >= len(content) || content[pos] != '"' {
		return Token{}, pos, NewParseError(startPos+pos, "expected opening quote", content)
	}

	tokenStart := pos
	pos++ // Skip opening quote

	var value strings.Builder
	escaped := false

	for pos < len(content) {
		ch := content[pos]

		if escaped {
			// Handle escape sequences
			switch ch {
			case 'n':
				value.WriteByte('\n')
			case 'r':
				value.WriteByte('\r')
			case 't':
				value.WriteByte('\t')
			case '\\':
				value.WriteByte('\\')
			case '"':
				value.WriteByte('"')
			case ',':
				value.WriteByte(',')
			case '0':
				value.WriteByte(0) // Null byte
			default:
				// For unknown escapes, include the backslash and character
				value.WriteByte('\\')
				value.WriteByte(ch)
			}
			escaped = false
			pos++
			continue
		}

		if ch == '\\' {
			escaped = true
			pos++
			continue
		}

		if ch == '"' {
			// Check for doubled quote (PostgreSQL escape for literal quote)
			if pos+1 < len(content) && content[pos+1] == '"' {
				value.WriteByte('"')
				pos += 2 // Skip both quotes
				continue
			}
			// End of quoted string
			rawValue := content[tokenStart : pos+1]
			pos++ // Skip closing quote
			return Token{
				Value:    value.String(),
				RawValue: rawValue,
				Quoted:   true,
				Position: startPos + tokenStart,
			}, pos, nil
		}

		value.WriteByte(ch)
		pos++
	}

	// Reached end without closing quote
	return Token{}, pos, NewParseError(startPos+tokenStart, "unclosed quoted string", content)
}

// parseUnquotedValue parses an unquoted value starting at pos.
func parseUnquotedValue(content string, pos int, startPos int) (Token, int, error) {
	tokenStart := pos
	var value strings.Builder

	for pos < len(content) {
		ch := content[pos]

		// End of unquoted value
		if ch == ',' || ch == ')' || unicode.IsSpace(rune(ch)) {
			break
		}

		// Quote in middle of unquoted value is an error
		if ch == '"' {
			return Token{}, pos, NewParseError(startPos+pos, "unexpected quote in unquoted value", content)
		}

		value.WriteByte(ch)
		pos++
	}

	rawValue := content[tokenStart:pos]

	return Token{
		Value:    strings.TrimSpace(value.String()),
		RawValue: rawValue,
		Quoted:   false,
		Position: startPos + tokenStart,
	}, pos, nil
}

// EscapeString escapes a string for inclusion in a row literal.
// It adds surrounding quotes and escapes special characters.
func EscapeString(s string) string {
	var b strings.Builder
	b.WriteByte('"')

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case 0:
			b.WriteString(`\0`)
		default:
			b.WriteByte(ch)
		}
	}

	b.WriteByte('"')
	return b.String()
}

// BuildRowLiteral builds a row literal string from a slice of string values.
// Values are expected to already be properly formatted (quoted strings, unquoted numbers).
func BuildRowLiteral(values []string) string {
	var b strings.Builder
	b.WriteByte('(')
	for i, v := range values {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(v)
	}
	b.WriteByte(')')
	return b.String()
}
