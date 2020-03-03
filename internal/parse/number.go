package parse

import (
	"errors"
	"strconv"

	"github.com/alistairjudson/cronparse/internal/numberer"
)

// Numberer is an interface that represents the ability to get a set of numbers
// from a statement/numberer of a statement
type Numberer interface {
	Numbers() []int
}

// Provider is a type that represents a type that can give you numbers for
// a numberer
type NumbererProvider interface {
	Numberer(parts Part) (Numberer, error)
}

// NewBaseNumbererFactory is a type that can give you a new instance of the BaseNumbererFactory
func NewBaseNumbererFactory(factory numberer.Provider) BaseNumbererFactory {
	return BaseNumbererFactory{
		Factory: factory,
	}
}

// BaseNumbererFactory is a type that wraps numberer.Provider and will return the correct
// numberer for a given parsed part
type BaseNumbererFactory struct {
	Factory numberer.Provider
}

// Numberer implements the Provider interface and will provide the correct Numberer
// based on pattern matching specific patterns of tokens
func (b BaseNumbererFactory) Numberer(parts Part) (Numberer, error) {
	types := parts.Types()
	switch {
	case types.StartsWith(TokenTypeAny):
		return b.Factory.Any(), nil
	case types.Contains(TokenTypeDash):
		return b.Factory.Range(parts[0].Value, parts[2].Value)
	case types.StartsWith(TokenTypeNumber) && !types.Contains(TokenTypeDash):
		return b.Factory.Number(parts[0].Value)
	}
	return nil, errors.New("numberer does not match any valid patterns")
}

// NewStepNumbererFactory will create a new instance of StepNumbererFactory wrapping a base
// factory
func NewStepNumbererFactory(base NumbererProvider) StepNumbererFactory {
	return StepNumbererFactory{
		Base: base,
	}
}

// StepNumbererFactory is a decorator of a Provider
type StepNumbererFactory struct {
	Base NumbererProvider
}

// Numberer implements Provider and will return a StepNumberer if the
// part contains a step, otherwise it will return the base
func (s StepNumbererFactory) Numberer(part Part) (Numberer, error) {
	types := part.Types()
	if !types.Contains(TokenTypeSlash) {
		return s.Base.Numberer(part)
	}
	base, err := s.Base.Numberer(part)
	if err != nil {
		return nil, err
	}
	return NewStepNumberer(base, part[len(part)-1].Value)
}

// NewStepNumberer will return a Numberer that will only output numbers in steps
// from the Numberer that it decorates
func NewStepNumberer(base Numberer, stepStr string) (Numberer, error) {
	step, err := strconv.ParseInt(stepStr, 10, 32)
	if err != nil {
		return nil, err
	}
	return StepNumberer{
		Base: base,
		Step: int(step),
	}, nil
}

// StepNumberer is a type that decorates Numberer, and only returns values from
// the base if they are in the step of the Step member
type StepNumberer struct {
	Base Numberer
	Step int
}

// Numbers implements Numberer and will return the number only if
// step mod index == 0
func (s StepNumberer) Numbers() []int {
	numbers := s.Base.Numbers()
	output := make([]int, 0, len(numbers)/s.Step)
	for i, number := range numbers {
		if i%s.Step != 0 {
			continue
		}
		output = append(output, number)
	}
	return output
}
