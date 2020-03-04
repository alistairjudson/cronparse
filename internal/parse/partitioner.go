package parse

import "errors"

// Part is a group of tokens, that represent a single item, within a field of a
// cron expression, e.g. (0-30/5) within an a field like "0-30/5,40,50"
type Part []Token

// Types returns you the types of all of the statement types contained within a pary
func (p Part) Types() Types {
	types := make(Types, 0, len(p))
	for _, token := range p {
		types = append(types, token.Type)
	}
	return types
}

// NewPartitioner is a type that can parse a field of a cron statement into it's
// constituent parts
func NewPartitioner() Partitioner {
	return Partitioner{
		TokenSource: NewTokenSource,
	}
}

// PartsProvider is an interface that represents a type that can parse a field of
// a cron expression into it's constituent parts
type PartsProvider interface {
	Parts(input string) ([]Part, error)
}

// Partitioner is a type that adapts a Tokeniser into a list of Parts
type Partitioner struct {
	TokenSource func(input string) TokenSource
}

// Parts will read all of the tokens from the Tokeniser and will split them based on the
// commas in the input into their parts
func (p Partitioner) Parts(input string) ([]Part, error) {
	tokens := p.TokenSource(input).Tokens()
	parts := make([]Part, 0)
	part := make(Part, 0)
	for token := range tokens {
		if token.Type == TokenTypeError {
			return nil, errors.New(token.Value)
		}
		if token.Type == TokenTypeComma {
			parts = append(parts, part)
			part = make(Part, 0)
			continue
		}
		part = append(part, token)
	}
	parts = append(parts, part)
	return parts, nil
}
