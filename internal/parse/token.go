package parse

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

var typeMap = map[TokenType]string{
	TokenTypeError:  "error",
	TokenTypeAny:    "any",
	TokenTypeComma:  "comma",
	TokenTypeDash:   "dash",
	TokenTypeSlash:  "slash",
	TokenTypeNumber: "number",
}

// TokenType is a type that represents the type of a value within
// the field of a cron statement
type TokenType int

// String implements fmt.Stringer returning a string representation of
// the type of token
func (t TokenType) String() string {
	return typeMap[t]
}

// Types is a slice of token types, allows methods to be added to allow
// for pattern matching different sequences of tokens
type Types []TokenType

// StartsWith tells you whether the first type within a slice of TokenTypes
// is the same as the type provided
func (t Types) StartsWith(typ TokenType) bool {
	if len(t) == 0 {
		return false
	}
	return t[0] == typ
}

// Contains tells you whether a specific TokenType is contained within this
// array of types
func (t Types) Contains(search TokenType) bool {
	for _, typ := range t {
		if typ == search {
			return true
		}
	}
	return false
}

// Types of tokens in a cron expression
const (
	TokenTypeError TokenType = iota
	TokenTypeAny
	TokenTypeComma
	TokenTypeDash
	TokenTypeSlash
	TokenTypeNumber
)

const eof = -1

// StateFunc is a function that represents a specific state within the
// state machine
type StateFunc func(t *Tokeniser) StateFunc

// Token is a single token emitted from the state machine
type Token struct {
	Type  TokenType
	Value string
}

// NewTokenSource will return, and start the Tokeniser for a given input string
func NewTokenSource(input string) TokenSource {
	tokeniser := NewTokeniser(input)
	go tokeniser.Run()
	return tokeniser
}

// TokenSource is an interface that represents a type that can emmit a stream of
// tokens
type TokenSource interface {
	Tokens() chan Token
}

// NewTokeniser will return a new *Tokeniser with the start state as lexField
func NewTokeniser(input string) *Tokeniser {
	return &Tokeniser{
		Input:      input,
		tokens:     make(chan Token),
		StartState: lexField,
	}
}

// Tokeniser is a type that provides utilities for creating a state machine to
// tokenise an input string
type Tokeniser struct {
	Input             string
	Start, Pos, Width int
	tokens            chan Token
	StartState        StateFunc
}

// Tokens implements the TokenSource interface an returns the channel that Tokens will
// be streamed out on
func (t *Tokeniser) Tokens() chan Token {
	return t.tokens
}

// Run will run the Tokeniser until the state returned is nil
func (t *Tokeniser) Run() {
	for state := t.StartState; state != nil; {
		state = state(t)
	}
	close(t.tokens)
}

// Next will move to the next rune in the input string
func (t *Tokeniser) Next() rune {
	if t.Pos >= len(t.Input) {
		t.Width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(t.Input[t.Pos:])
	t.Width = w
	t.Pos += w
	return r
}

// Backup will move back one character in the input stream
func (t *Tokeniser) Backup() {
	t.Pos -= t.Width
}

// Emit will emit a token of the given type, and set the start
// of the current section to the current position
func (t *Tokeniser) Emit(typ TokenType) {
	t.tokens <- Token{
		Type:  typ,
		Value: t.Input[t.Start:t.Pos],
	}
	t.Start = t.Pos
}

// AcceptNumber will call next for as long as the rune is a numeric value
func (t *Tokeniser) AcceptNumber() {
	for r := t.Next(); r != eof && unicode.IsDigit(r); r = t.Next() {
	}
	t.Backup()
}

// Errorf will return an error token, and will return nil as a StateFunc,
// ending tokenisation
func (t *Tokeniser) Errorf(format string, args ...interface{}) StateFunc {
	t.tokens <- Token{
		Type:  TokenTypeError,
		Value: fmt.Sprintf(format, args...),
	}
	return nil
}
