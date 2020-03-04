package parse_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alistairjudson/cronparse/internal/numberer"
	"github.com/alistairjudson/cronparse/internal/parse"
)

type stubNumbererProvider struct {
	number func(numstr string) (numberer.Number, error)
	rnge   func(startString, endString string) (numberer.Range, error)
	any    func() numberer.Any
}

func (s stubNumbererProvider) Number(numstr string) (numberer.Number, error) {
	return s.number(numstr)
}

func (s stubNumbererProvider) Range(startString, endString string) (numberer.Range, error) {
	return s.rnge(startString, endString)
}

func (s stubNumbererProvider) Any() numberer.Any {
	return s.any()
}

func TestBaseNumbererFactory_NumbererNumber(t *testing.T) {
	timesCalled := 0
	np := stubNumbererProvider{
		number: func(numstr string) (number numberer.Number, err error) {
			timesCalled++
			return 0, nil
		},
	}
	tokens := parse.Part{
		{
			Type:  parse.TokenTypeNumber,
			Value: "1",
		},
	}
	_, err := parse.NewBaseNumbererFactory(np).Numberer(tokens)
	if err != nil {
		t.Fatal(err)
	}
	if timesCalled != 1 {
		t.Fatalf("expected timesCalled to be 1, got (%d)", timesCalled)
	}
}

func TestBaseNumbererFactory_NumbererRange(t *testing.T) {
	timesCalled := 0
	np := stubNumbererProvider{
		rnge: func(startString, endString string) (n numberer.Range, err error) {
			timesCalled++
			return numberer.Range{
				Start: 1,
				End:   5,
			}, nil
		},
	}
	tokens := parse.Part{
		{
			Type:  parse.TokenTypeNumber,
			Value: "1",
		},
		{
			Type:  parse.TokenTypeDash,
			Value: "-",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "5",
		},
	}
	_, err := parse.NewBaseNumbererFactory(np).Numberer(tokens)
	if err != nil {
		t.Fatal(err)
	}
	if timesCalled != 1 {
		t.Fatalf("expected timesCalled to be 1, got (%d)", timesCalled)
	}
}

func TestBaseNumbererFactory_NumbererAny(t *testing.T) {
	timesCalled := 0
	np := stubNumbererProvider{
		any: func() numberer.Any {
			timesCalled++
			return numberer.Any{}
		},
	}
	tokens := parse.Part{
		{
			Type:  parse.TokenTypeAny,
			Value: "*",
		},
	}
	_, err := parse.NewBaseNumbererFactory(np).Numberer(tokens)
	if err != nil {
		t.Fatal(err)
	}
	if timesCalled != 1 {
		t.Fatalf("expected timesCalled to be 1, got (%d)", timesCalled)
	}
}

func TestBaseNumbererFactory_NumbererFails(t *testing.T) {
	timesCalled := 0
	np := stubNumbererProvider{
		number: func(numstr string) (number numberer.Number, err error) {
			timesCalled++
			return 0, nil
		},
		rnge: func(startString, endString string) (n numberer.Range, err error) {
			timesCalled++
			return numberer.Range{
				Start: 1,
				End:   5,
			}, nil
		},
		any: func() numberer.Any {
			timesCalled++
			return numberer.Any{}
		},
	}
	tokens := parse.Part{}
	_, err := parse.NewBaseNumbererFactory(np).Numberer(tokens)
	if err == nil {
		t.Fatal("expected an error got none")
	}
	if timesCalled != 0 {
		t.Fatalf("expected timesCalled to be 0, got (%d)", timesCalled)
	}
}

type stubParseNumbererProvider struct {
	numberer func(parts parse.Part) (parse.Numberer, error)
}

func (s stubParseNumbererProvider) Numberer(parts parse.Part) (parse.Numberer, error) {
	return s.numberer(parts)
}

func TestStepNumbererFactory_NumbererSucceedsSteps(t *testing.T) {
	snp := stubParseNumbererProvider{
		numberer: func(parts parse.Part) (p parse.Numberer, err error) {
			return numberer.NewRange(0, 59)
		},
	}

	tokens := parse.Part{
		{
			Type:  parse.TokenTypeNumber,
			Value: "0",
		},
		{
			Type:  parse.TokenTypeDash,
			Value: "-",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "59",
		},
		{
			Type:  parse.TokenTypeSlash,
			Value: "/",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "15",
		},
	}

	stepNumberProvider := parse.NewStepNumbererFactory(snp)
	stepNumberer, err := stepNumberProvider.Numberer(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expectedNumbers := []int{0, 15, 30, 45}
	gotNumbers := stepNumberer.Numbers()
	if !reflect.DeepEqual(expectedNumbers, gotNumbers) {
		t.Fatalf("expected (%+v), got (%+v)", expectedNumbers, gotNumbers)
	}
}

func TestStepNumbererFactory_NumbererFailsInvalidStep(t *testing.T) {
	snp := stubParseNumbererProvider{
		numberer: func(parts parse.Part) (p parse.Numberer, err error) {
			return numberer.NewRange(0, 59)
		},
	}

	tokens := parse.Part{
		{
			Type:  parse.TokenTypeNumber,
			Value: "0",
		},
		{
			Type:  parse.TokenTypeDash,
			Value: "-",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "59",
		},
		{
			Type:  parse.TokenTypeSlash,
			Value: "/",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "zzzzz",
		},
	}

	stepNumberProvider := parse.NewStepNumbererFactory(snp)
	_, err := stepNumberProvider.Numberer(tokens)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}

func TestStepNumbererFactory_NumbererNoStep(t *testing.T) {
	snp := stubParseNumbererProvider{
		numberer: func(parts parse.Part) (p parse.Numberer, err error) {
			return numberer.NewRange(0, 6)
		},
	}

	tokens := parse.Part{
		{
			Type:  parse.TokenTypeNumber,
			Value: "0",
		},
		{
			Type:  parse.TokenTypeDash,
			Value: "-",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "59",
		},
	}

	stepNumberProvider := parse.NewStepNumbererFactory(snp)
	stepNumberer, err := stepNumberProvider.Numberer(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expectedNumbers := []int{0, 1, 2, 3, 4, 5, 6}
	gotNumbers := stepNumberer.Numbers()
	if !reflect.DeepEqual(expectedNumbers, gotNumbers) {
		t.Fatalf("expected (%+v), got (%+v)", expectedNumbers, gotNumbers)
	}
}

func TestStepNumbererFactory_NumbererBaseFails(t *testing.T) {
	snp := stubParseNumbererProvider{
		numberer: func(parts parse.Part) (p parse.Numberer, err error) {
			return nil, errors.New("an error")
		},
	}

	tokens := parse.Part{
		{
			Type:  parse.TokenTypeNumber,
			Value: "0",
		},
		{
			Type:  parse.TokenTypeDash,
			Value: "-",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "59",
		},
		{
			Type:  parse.TokenTypeSlash,
			Value: "/",
		},
		{
			Type:  parse.TokenTypeNumber,
			Value: "15",
		},
	}

	stepNumberProvider := parse.NewStepNumbererFactory(snp)
	_, err := stepNumberProvider.Numberer(tokens)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}
