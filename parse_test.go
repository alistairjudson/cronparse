package cronparse_test

import (
	"testing"

	"github.com/alistairjudson/cronparse"
)

func TestParser_ParseSucceeds(t *testing.T) {
	expression := []string{"*/15", "0", "1,15", "*", "1-5"}
	res, err := cronparse.CronParser.Parse(expression)
	if err != nil {
		t.Fatal(err)
	}
	expectedStrings := []string{
		"minute         0 15 30 45",
		"hour           0",
		"day of month   1 15",
		"month          1 2 3 4 5 6 7 8 9 10 11 12",
		"day of week    1 2 3 4 5",
	}
	for i, res := range res {
		gotString := res.String()
		if expectedStrings[i] != gotString {
			t.Errorf("expected string (%s), got (%s)", expectedStrings[i], gotString)
		}
	}
}

func TestParser_ParseFailsNotEnoughComponents(t *testing.T) {
	expression := []string{"*/15", "0", "1,15", "*"}
	_, err := cronparse.CronParser.Parse(expression)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}

func TestParser_ParseFailsBadComponent(t *testing.T) {
	expression := []string{"*/15", "0", "1,15", "**", "1-5"}
	_, err := cronparse.CronParser.Parse(expression)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}
