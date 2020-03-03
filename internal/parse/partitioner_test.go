package parse_test

import (
	"testing"

	"github.com/alistairjudson/cronparse/internal/parse"
)

type stubTokenSource struct {
	tokens chan parse.Token
}

func (s stubTokenSource) Tokens() chan parse.Token {
	return s.tokens
}

func TestPartitioner_PartsFails(t *testing.T) {
	sts := stubTokenSource{
		tokens: make(chan parse.Token, 1),
	}
	sts.tokens <- parse.Token{
		Type:  parse.TokenTypeError,
		Value: "an error",
	}
	close(sts.tokens)
	partitioner := parse.NewPartitioner()
	partitioner.TokenSource = func(input string) parse.TokenSource {
		return sts
	}

	_, err := partitioner.Parts("foo")
	if err == nil {
		t.Fatal(err)
	}
}

func TestPartitioner_PartsSucceds(t *testing.T) {
	sts := stubTokenSource{
		tokens: make(chan parse.Token, 3),
	}
	sts.tokens <- parse.Token{
		Type:  parse.TokenTypeNumber,
		Value: "1",
	}
	sts.tokens <- parse.Token{
		Type:  parse.TokenTypeComma,
		Value: ",",
	}
	sts.tokens <- parse.Token{
		Type:  parse.TokenTypeNumber,
		Value: "1",
	}
	close(sts.tokens)
	partitioner := parse.NewPartitioner()
	partitioner.TokenSource = func(input string) parse.TokenSource {
		return sts
	}

	parts, err := partitioner.Parts("foo")
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got (%d)", len(parts))
	}
}
