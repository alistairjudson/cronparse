package parse_test

import (
	"reflect"
	"testing"

	"github.com/alistairjudson/cronparse/internal/parse"
)

func TestNewTokenSourceFails(t *testing.T) {
	tests := []struct {
		name                 string
		field                string
		expectedErrorMessage string
	}{
		{
			name:                 "empty",
			field:                "",
			expectedErrorMessage: "input cannot be empty",
		},
		{
			name:                 "letter",
			field:                "a",
			expectedErrorMessage: "(a) is unexpected at the start of a statement, expected (* or [0-9]+)",
		},
		{
			name:                 "any + unexpected",
			field:                "*a",
			expectedErrorMessage: "(a) is unexpected after (*), only (/,) expected",
		},
		{
			name:                 "invalid step",
			field:                "*/a",
			expectedErrorMessage: "(a) is unexpected after a step, only numbers are expected",
		},
		{
			name:                 "invalid after step",
			field:                "*/12-",
			expectedErrorMessage: "(-) is unexpected after a step, only (,) is expected",
		},
		{
			name:                 "invalid after number",
			field:                "12a",
			expectedErrorMessage: "(a) is unexpected after a number, only (,-/) expected",
		},
		{
			name:                 "invalid in range",
			field:                "12-a",
			expectedErrorMessage: "(a) is unexpected in a range, only digits are expected",
		},
		{
			name:                 "invalid in range",
			field:                "12-13a",
			expectedErrorMessage: "(a) is unexpected after a range, only (/,) expected",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := parse.NewTokenSource(test.field).Tokens()
			var last parse.Token
			for token := range tokens {
				last = token
			}
			if last.Type != parse.TokenTypeError {
				t.Fatalf("expected token type to be error, got (%s)", last.Type)
			}
			gotMessage := last.Value
			if gotMessage != test.expectedErrorMessage {
				t.Fatalf("expected error message to be (%s), got (%s)", test.expectedErrorMessage, gotMessage)
			}
		})
	}
}

func TestNewTokenSourceSucceeds(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		expectedTypes parse.Types
	}{
		{
			name:          "any",
			field:         "*",
			expectedTypes: parse.Types{parse.TokenTypeAny},
		},
		{
			name:  "any with step",
			field: "*/12",
			expectedTypes: parse.Types{
				parse.TokenTypeAny,
				parse.TokenTypeSlash,
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "number",
			field: "12",
			expectedTypes: parse.Types{
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "any + any with step",
			field: "*,*/12",
			expectedTypes: parse.Types{
				parse.TokenTypeAny,
				parse.TokenTypeComma,
				parse.TokenTypeAny,
				parse.TokenTypeSlash,
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "number + number",
			field: "12,13",
			expectedTypes: parse.Types{
				parse.TokenTypeNumber,
				parse.TokenTypeComma,
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "range + any with step",
			field: "1-4,*/15",
			expectedTypes: parse.Types{
				parse.TokenTypeNumber,
				parse.TokenTypeDash,
				parse.TokenTypeNumber,
				parse.TokenTypeComma,
				parse.TokenTypeAny,
				parse.TokenTypeSlash,
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "number with step",
			field: "4/12",
			expectedTypes: parse.Types{
				parse.TokenTypeNumber,
				parse.TokenTypeSlash,
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "range with step",
			field: "1-5/2",
			expectedTypes: parse.Types{
				parse.TokenTypeNumber,
				parse.TokenTypeDash,
				parse.TokenTypeNumber,
				parse.TokenTypeSlash,
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "range with step and number",
			field: "1-5/2,1",
			expectedTypes: parse.Types{
				parse.TokenTypeNumber,
				parse.TokenTypeDash,
				parse.TokenTypeNumber,
				parse.TokenTypeSlash,
				parse.TokenTypeNumber,
				parse.TokenTypeComma,
				parse.TokenTypeNumber,
			},
		},
		{
			name:  "range",
			field: "1-5",
			expectedTypes: parse.Types{
				parse.TokenTypeNumber,
				parse.TokenTypeDash,
				parse.TokenTypeNumber,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := parse.NewTokenSource(test.field).Tokens()
			part := make(parse.Part, 0)
			for token := range tokens {
				part = append(part, token)
			}
			gotTypes := part.Types()
			if !reflect.DeepEqual(test.expectedTypes, gotTypes) {
				t.Fatalf("expected (%+v), got (%+v)", test.expectedTypes, gotTypes)
			}
		})
	}
}

func TestTypes_Contains(t *testing.T) {
	types := parse.Types{parse.TokenTypeNumber}
	if !types.Contains(parse.TokenTypeNumber) {
		t.Fatal("expected types to contain number")
	}
	if types.Contains(parse.TokenTypeAny) {
		t.Fatal("expected types not to contain any")
	}
}

func TestTypes_StartsWith(t *testing.T) {
	types := parse.Types{parse.TokenTypeNumber}
	if !types.StartsWith(parse.TokenTypeNumber) {
		t.Fatal("expected types to start with number")
	}
	if types.StartsWith(parse.TokenTypeAny) {
		t.Fatal("expected types not to start with any")
	}
	empty := parse.Types{}
	if empty.StartsWith(parse.TokenTypeAny) {
		t.Fatal("expected types not to contain anything")
	}
}

func TestTokenType_String(t *testing.T) {
	gotName := parse.TokenTypeAny.String()
	expectedName := "any"
	if expectedName != gotName {
		t.Fatalf("expected (%s), got (%s)", expectedName, gotName)
	}
}
