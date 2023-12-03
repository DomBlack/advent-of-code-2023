package day03

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/DomBlack/advent-of-code-2023/pkg/stream"
	"github.com/cockroachdb/errors"
)

type TokenType uint8

const (
	Number TokenType = iota
	Part
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int

	Numbers []string // If this is a part, the numbers which are adjacent to it
}

func (t Token) String() string {
	switch t.Type {
	case Number:
		return fmt.Sprintf("Number(%s) %d:%d", t.Value, t.Line, t.Column)
	case Part:
		return fmt.Sprintf("Part(%s) %d:%d", t.Value, t.Line, t.Column)
	}

	return fmt.Sprintf("Unknown(%s)", t.Value)
}

func (t Token) StartCol() int {
	return t.Column
}

func (t Token) EndCol() int {
	return t.Column + len(t.Value) - 1
}

// lexer implements a simple lexer for the day 3 problem
type lexer struct {
	reader     *bufio.Reader
	firstError error
	current    rune
	line       int
	column     int
}

func parseSchematic(input []byte) stream.Stream[Token] {
	p := &lexer{
		reader: bufio.NewReader(bytes.NewReader(input)),
		line:   1,
		column: 1,
	}
	p.current, _, p.firstError = p.reader.ReadRune()

	return p
}

func (p *lexer) Next() (Token, error) {
	if p.firstError != nil {
		return Token{}, p.firstError
	}

	if !p.consumeWhitespace() {
		return Token{}, p.firstError
	}

	switch {
	case p.current >= '0' && p.current <= '9':
		return p.consumeNumber()
	default:
		return p.consumePart()
	}
}

// consumeWhitespace consumes all whitespace and returns true if there is another rune to read
func (p *lexer) consumeWhitespace() bool {
	for p.current == '.' || p.current == '\n' {
		if !p.consume() {
			return false
		}
	}

	return true
}

func (p *lexer) consumeNumber() (Token, error) {
	var buf bytes.Buffer

	startLine := p.line
	startColumn := p.column

	for p.current >= '0' && p.current <= '9' {
		buf.WriteRune(p.current)

		if !p.consume() {
			break
		}
	}

	return Token{
		Type:   Number,
		Value:  buf.String(),
		Line:   startLine,
		Column: startColumn,
	}, nil
}

func (p *lexer) consumePart() (Token, error) {
	t := Token{
		Type:   Part,
		Value:  fmt.Sprintf("%c", p.current),
		Line:   p.line,
		Column: p.column,
	}

	p.consume()
	return t, nil
}

// consume consumes the current rune and returns true if there is another rune to read
func (p *lexer) consume() bool {
	// Read the next rune
	next, _, err := p.reader.ReadRune()
	if err != nil {
		if p.firstError == nil {
			p.firstError = errors.Wrap(err, "failed to read rune")
		}
		return false
	}

	// Track the position of the current rune
	p.current = next
	if p.current == '\n' {
		p.line++
		p.column = 0
	} else {
		p.column++
	}

	return true
}
