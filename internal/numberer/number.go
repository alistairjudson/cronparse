package numberer

import (
	"fmt"
	"strconv"
)

var (
	// MinuteFactory is a factory to produce valid minute values
	MinuteFactory = Must(NewFactory("minute", 0, 59))

	// HourFactory is a factory to produce valid hour values
	HourFactory = Must(NewFactory("hour", 0, 23))

	// DayOfMonthFactory is a factory to produce valid day of month values
	DayOfMonthFactory = Must(NewFactory("dayOfMonth", 1, 31))

	// MonthFactory is a factory to produce valid month values
	MonthFactory = Must(NewFactory("month", 1, 12))

	// DayOfWeekFactory is a factory that can produce valid day of week values
	DayOfWeekFactory = Must(NewFactory("dayOfWeek", 0, 6))
)

// Must will panic if an error is passed to it, used for factory variables
func Must(factory Factory, err error) Factory {
	if err != nil {
		panic(err)
	}
	return factory
}

// NewFactory will create a new factory for a given range, with a
// name, the name is used for errors.
func NewFactory(name string, start, end int) (Factory, error) {
	rnge, err := NewRange(start, end)
	if err != nil {
		return Factory{}, fmt.Errorf("(%s) %w", name, err)
	}
	return Factory{
		rnge: rnge,
		name: name,
	}, nil
}

// Factory is a type that can create Numberers for different strings
type Factory struct {
	rnge Range
	name string
}

// Number will create a Number Numberer from a string, and validate that
// it is within the range of the factory (along with validating that the
// input is an int)
func (f Factory) Number(numstr string) (Number, error) {
	num, err := strconv.ParseInt(numstr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("(%s) failed to parse (%s) as int got (%w)", f.name, numstr, err)
	}
	if int(num) < f.rnge.Start || int(num) > f.rnge.End {
		return 0, fmt.Errorf("(%s) number (%d) must be in range (%d-%d)", f.name, num, f.rnge.Start, f.rnge.End)
	}
	return Number(num), nil
}

// Range will create a range from the given strings, validating that the
// range that they within the range of the factory (along with validating
// that the inputs is an int)
func (f Factory) Range(startString, endString string) (Range, error) {
	start, err := strconv.ParseInt(startString, 10, 32)
	if err != nil {
		return Range{}, fmt.Errorf("(%s) failed to parse (%s) as int got (%w)", f.name, startString, err)
	}
	end, err := strconv.ParseInt(endString, 10, 32)
	if err != nil {
		return Range{}, fmt.Errorf("(%s) failed to parse (%s) as int got (%w)", f.name, endString, err)
	}
	if int(start) < f.rnge.Start || int(end) > f.rnge.End {
		return Range{}, fmt.Errorf("(%s) range (%d - %d) must be in range (%d - %d)", f.name, start, end, f.rnge.Start, f.rnge.End)
	}
	return NewRange(int(start), int(end))
}

// Any will return a Numberer that will return the entire range of numbers
func (f Factory) Any() Any {
	return f.rnge.Numbers()
}

// Number is a Numberer that will return a single number
type Number int

// Numbers implements Numberer and will return the number as an array of
// numbers
func (n Number) Numbers() []int {
	return []int{int(n)}
}

// NewRange will create a new Range, validating that the end is greater
// than the end
func NewRange(start, end int) (Range, error) {
	if start > end {
		return Range{}, fmt.Errorf("start (%d), cannot be greater than end (%d)", start, end)
	}
	return Range{
		Start: start,
		End:   end,
	}, nil
}

// Range is a Numberer that will return all of the numbers in a given range
type Range struct {
	Start, End int
}

// Numbers implements Numberer that will return all of the numbers between
// the start and end values
func (r Range) Numbers() []int {
	numbers := make([]int, 0, r.End-r.Start)
	for i := r.Start; i <= r.End; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}

// Any is a type that can be filled with all the numbers from a range
type Any []int

// Numbers implements Numberer and will return all of the values it
// contains
func (a Any) Numbers() []int {
	return a
}

// Provider is an interface to abstract the Factory
type Provider interface {
	Number(numstr string) (Number, error)
	Range(startString, endString string) (Range, error)
	Any() Any
}
