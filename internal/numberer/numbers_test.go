package numberer_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alistairjudson/cronparse/internal/numberer"
)

func TestRange_Numbers(t *testing.T) {
	rng := numberer.Range{
		Start: 1,
		End:   5,
	}
	gotRange := rng.Numbers()
	expectedRange := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(expectedRange, gotRange) {
		t.Fatalf("expected (%+v), got (%+v)", expectedRange, gotRange)
	}
}

func TestNewRangeFails(t *testing.T) {
	_, err := numberer.NewRange(10, 0)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}

func TestNewRangeSucceeds(t *testing.T) {
	_, err := numberer.NewRange(0, 10)
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
}

func TestNewFactoryFails(t *testing.T) {
	_, err := numberer.NewFactory("minute", 10, 0)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}

func TestNewFactorySucceeds(t *testing.T) {
	_, err := numberer.NewFactory("minute", 0, 10)
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
}

func TestFactory_NumberFails(t *testing.T) {
	tests := []struct {
		name      string
		numString string
	}{
		{
			name:      "invalid number",
			numString: "zzzzz",
		},
		{
			name:      "number out of range",
			numString: "70",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := numberer.MinuteFactory.Number(test.numString); err == nil {
				t.Fatal("expected an error, got none")
			}
		})
	}
}

func TestFactory_RangeFails(t *testing.T) {
	tests := []struct {
		name             string
		startStr, endStr string
	}{
		{
			name:     "invalid start",
			startStr: "zzz",
			endStr:   "10",
		},
		{
			name:     "invalid end",
			startStr: "10",
			endStr:   "zzz",
		},
		{
			name:     "out of range start",
			startStr: "-1",
			endStr:   "5",
		},
		{
			name:     "out of range end",
			startStr: "2",
			endStr:   "30",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := numberer.DayOfWeekFactory.Range(test.startStr, test.endStr)
			if err == nil {
				t.Fatal("expected an error got none")
			}
		})
	}
}

func TestNumber_Numbers(t *testing.T) {
	expectedNumbers := []int{1}
	gotNumbers := numberer.Number(1).Numbers()
	if !reflect.DeepEqual(expectedNumbers, gotNumbers) {
		t.Fatalf("expected numbers (%+v), got numbers (%+v)", expectedNumbers, gotNumbers)
	}
}

func TestFactory_Any(t *testing.T) {
	expectedNumbers := []int{0, 1, 2, 3, 4, 5, 6}
	gotNumbers := numberer.DayOfWeekFactory.Any().Numbers()
	if !reflect.DeepEqual(expectedNumbers, gotNumbers) {
		t.Fatalf("expected numbers (%+v), got numbers (%+v)", expectedNumbers, gotNumbers)
	}
}

func TestFactory_NumberSucceeds(t *testing.T) {
	_, err := numberer.DayOfWeekFactory.Number("1")
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
}

func TestFactory_RangeSucceeds(t *testing.T) {
	_, err := numberer.DayOfWeekFactory.Range("1", "3")
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
}

func TestMust(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expected a panic, got nil")
		}
	}()
	numberer.Must(numberer.Factory{}, errors.New("an error"))
}
