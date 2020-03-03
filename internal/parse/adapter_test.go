package parse_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alistairjudson/cronparse/internal/parse"
)

type stubPartsProvider struct {
	parts func(input string) ([]parse.Part, error)
}

func (s stubPartsProvider) Parts(input string) ([]parse.Part, error) {
	return s.parts(input)
}

func TestAdapter_ParseFailsPartsFactory(t *testing.T) {
	spp := stubPartsProvider{parts: func(input string) (parts []parse.Part, err error) {
		return nil, errors.New("an error")
	}}
	adapter := parse.NewAdapter(stubParseNumbererProvider{})
	adapter.PartsFactory = spp
	_, err := adapter.Parse("foo bar baz")
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}

func TestAdapter_ParseFailsNumbererFactory(t *testing.T) {
	expectedErr := errors.New("an error")
	spnp := stubParseNumbererProvider{numberer: func(parts parse.Part) (numberer parse.Numberer, err error) {
		return nil, expectedErr
	}}
	adapter := parse.NewAdapter(spnp)
	_, err := adapter.Parse("0-59/15")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to be (%s), got (%s)", expectedErr, err)
	}
}

func TestAdapter_ParseSucceeds(t *testing.T) {
	expectedNumbers := []int{0, 15, 30, 31, 45, 50, 51, 52, 53, 54, 55}
	numberer, err := parse.MinuteParser.Parse("*/15,30,31,50-55")
	if err != nil {
		t.Fatal(err)
	}

	gotNumbers := numberer.Numbers()
	if !reflect.DeepEqual(expectedNumbers, gotNumbers) {
		t.Fatalf("expected numbers to be (%+v), got (%+v)", expectedNumbers, gotNumbers)
	}
}
