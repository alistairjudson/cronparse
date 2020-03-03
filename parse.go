package cronparse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alistairjudson/cronparse/internal/parse"
)

// CronParser is a type that can parse the components of a cron expression and
// expand them into the values that they run on
var CronParser = Parser{
	newComponentParser("minute", parse.MinuteParser.Parse),
	newComponentParser("hour", parse.HourParser.Parse),
	newComponentParser("day of month", parse.DayOfMonthParser.Parse),
	newComponentParser("month", parse.MonthParser.Parse),
	newComponentParser("day of week", parse.DayOfWeekParser.Parse),
}

// Parser is a type that can hold multiple ComponentParsers
type Parser []ComponentParser

// Parse will parse all of the components with their respective component parser
func (p Parser) Parse(components []string) ([]ParsedComponent, error) {
	if len(components) != len(p) {
		return nil, fmt.Errorf("expected (%d) components, got (%d) components", len(p), len(components))
	}
	parsed := make([]ParsedComponent, 0, len(p))
	for i, componentParser := range p {
		num, err := componentParser.Parser.Parse(components[i])
		if err != nil {
			return nil, fmt.Errorf("(%s): %w", componentParser.Name, err)
		}
		parsed = append(parsed, ParsedComponent{
			Name:    componentParser.Name,
			Numbers: num.Numbers(),
		})
	}
	return parsed, nil
}

// PartParser is a type that can parse part of a cron expression
type PartParser interface {
	Parse(component string) (Numberer, error)
}

// Numberer is a type that can give you the expanded cron values
type Numberer interface {
	Numbers() []int
}

// ComponentParser is a parser for a specific portion of a cron expression
type ComponentParser struct {
	Name   string
	Parser PartParser
}

// ParsedComponent is the result of parsing a cron component
type ParsedComponent struct {
	Name    string
	Numbers []int
}

// String implements fmt.Stringer and pretty prints a parsed cron component
func (p ParsedComponent) String() string {
	stringNums := make([]string, 0, len(p.Numbers))
	for _, num := range p.Numbers {
		stringNums = append(stringNums, strconv.Itoa(num))
	}
	return fmt.Sprintf(
		"%-14s %s",
		p.Name,
		strings.Join(stringNums, " "),
	)
}

func newComponentParser(name string, parserFunc parserFunc) ComponentParser {
	return ComponentParser{
		Name:   name,
		Parser: parserFunc,
	}
}

type parserFunc func(component string) (parse.Numberer, error)

func (p parserFunc) Parse(component string) (Numberer, error) {
	return p(component)
}
