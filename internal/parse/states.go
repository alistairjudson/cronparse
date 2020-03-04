package parse

import "unicode"

// This files represents the all of the states that there can possibly be
// within a field of a cronparse expression. This is a state machine, it does not
// validate the values of the items, only that the syntax of the cronparse statement
// field matches the grammar of a cronparse field.

func lexField(t *Tokeniser) StateFunc {
	curr := t.Next()
	switch {
	case curr == '*':
		return lexAny
	case unicode.IsDigit(curr):
		return lexNumber
	case curr == eof:
		return t.Errorf("input cannot be empty")
	}
	return t.Errorf(`(%c) is unexpected at the start of a statement, expected (* or [0-9]+)`, curr)
}

func lexComma(t *Tokeniser) StateFunc {
	t.Emit(TokenTypeComma)
	return lexField
}

func lexAny(t *Tokeniser) StateFunc {
	t.Emit(TokenTypeAny)
	next := t.Next()
	switch next {
	case ',':
		return lexComma
	case '/':
		return lexStep
	case eof:
		return nil
	}
	return t.Errorf(`(%c) is unexpected after (*), only (/,) expected`, next)
}

func lexStep(t *Tokeniser) StateFunc {
	t.Emit(TokenTypeSlash)
	next := t.Next()
	if !unicode.IsDigit(next) {
		return t.Errorf("(%c) is unexpected after a step, only numbers are expected", next)
	}
	t.AcceptNumber()
	t.Emit(TokenTypeNumber)
	next = t.Next()
	switch next {
	case ',':
		return lexComma
	case eof:
		return nil
	}
	return t.Errorf("(%c) is unexpected after a step, only (,) is expected", next)
}

func lexNumber(t *Tokeniser) StateFunc {
	t.AcceptNumber()
	t.Emit(TokenTypeNumber)
	next := t.Next()
	switch next {
	case ',':
		return lexComma
	case '-':
		return lexRange
	case '/':
		return lexStep
	case eof:
		return nil
	}
	return t.Errorf("(%c) is unexpected after a number, only (,-/) expected", next)
}

func lexRange(t *Tokeniser) StateFunc {
	t.Emit(TokenTypeDash)
	next := t.Next()
	if !unicode.IsDigit(next) {
		return t.Errorf("(%c) is unexpected in a range, only digits are expected", next)
	}
	t.AcceptNumber()
	t.Emit(TokenTypeNumber)
	next = t.Next()
	switch next {
	case '/':
		return lexStep
	case ',':
		return lexComma
	case eof:
		return nil
	}
	return t.Errorf(`(%c) is unexpected after a range, only (/,) expected`, next)
}
