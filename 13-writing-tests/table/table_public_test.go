package table_test

import (
	"testing"

	"github.com/dahc36/learning-go/13-writing-tests/table"
)

func TestDoMath(t *testing.T) {
	data := []struct {
		name     string
		num1     int
		num2     int
		op       string
		expected int
		errMsg   string
	}{
		{"addition", 2, 2, "+", 4, ""},
		{"subtraction", 2, 2, "-", 0, ""},
		{"multiplication", 2, 2, "*", 4, ""},
		{"another_multiplication", 2, 3, "*", 6, ""},
		{"division", 2, 2, "/", 1, ""},
		{"bad_division", 2, 0, "/", 0, "division by zero"},
		{"bad_division", 2, 2, "?", 0, "operator ? not supported"},
	}
	for _, d := range data {
		result, err := table.DoMath(d.num1, d.num2, d.op)
		if result != d.expected {
			t.Errorf("Expected %d, got %d", d.expected, result)
		}
		// If the error had a custom type we could use errors.Is or errors.As for more robust checking
		// Instead of just comparing the strings
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		if errMsg != d.errMsg {
			t.Errorf("Expected error message %s, got %s", d.errMsg, errMsg)
		}
	}
}
