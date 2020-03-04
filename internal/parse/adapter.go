package parse

import (
	"sort"

	"github.com/alistairjudson/cronparse/internal/numberer"
)

var (
	// MinuteParser is a parser to parse the minute component of a cronparse expression
	MinuteParser = NewParser(numberer.MinuteFactory)

	// HourParser is a parser to parse the hour component of a cronparse expression
	HourParser = NewParser(numberer.HourFactory)

	// DayOfMonthParser is a parser to parse the day of month component of a cronparse expression
	DayOfMonthParser = NewParser(numberer.DayOfMonthFactory)

	// MonthParser is a parser to parse the month component of a cronparse expression
	MonthParser = NewParser(numberer.MonthFactory)

	// DayOfWeekParser is a parser to parse the day of week component of a cronparse expression
	DayOfWeekParser = NewParser(numberer.DayOfWeekFactory)
)

// NewParser will return you a new parser for a given numberer.Provider
func NewParser(provider numberer.Provider) Parser {
	return NewAdapter(
		NewStepNumbererFactory(
			NewBaseNumbererFactory(provider),
		),
	)
}

// Parser is the interface that represents a type that can parse numberer
// of a cronparse statement into a Numberer
type Parser interface {
	Parse(input string) (Numberer, error)
}

// NewAdapter will return you a new instance of an Adapter, with a
// Partitioner as its PartsFactory
func NewAdapter(provider NumbererProvider) Adapter {
	return Adapter{
		PartsFactory:    NewPartitioner(),
		NumbererFactory: provider,
	}
}

// Adapter is a type that can adapt the PartsProvider and the Provider
// into a Parser
type Adapter struct {
	PartsFactory    PartsProvider
	NumbererFactory NumbererProvider
}

// Parse implements Parser and will parse the statement into a Numberer, by
// parsing it into parts, and then into an AggregateNumberer
func (a Adapter) Parse(statement string) (Numberer, error) {
	parts, err := a.PartsFactory.Parts(statement)
	if err != nil {
		return nil, err
	}
	aggregate := make(AggregateNumberer, 0, len(parts))
	for _, part := range parts {
		partNumberer, err := a.NumbererFactory.Numberer(part)
		if err != nil {
			return nil, err
		}
		aggregate = append(aggregate, partNumberer)
	}
	return aggregate, nil
}

// AggregateNumberer is a Numberer that will grab values from all of the
// contained numberers combining them
type AggregateNumberer []Numberer

// Numbers implements the Numberer interface and will produce an ordered,
// de-duplicated list of all of the numbers contained within the AggregateNumberer
func (a AggregateNumberer) Numbers() []int {
	numberSet := map[int]struct{}{}
	for _, unaggregatedNumberer := range a {
		for _, number := range unaggregatedNumberer.Numbers() {
			numberSet[number] = struct{}{}
		}
	}
	numbers := make([]int, 0, len(numberSet))
	for number := range numberSet {
		numbers = append(numbers, number)
	}
	sort.Ints(numbers)
	return numbers
}
