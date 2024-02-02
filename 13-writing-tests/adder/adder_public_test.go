// Using packagename_test as the package for the test allows for testing of the
// public API of the package, because it doesn't have package level access
package adder_test

import (
	"testing"

	"github.com/dahc36/learning-go/13-writing-tests/adder"
)

func TestAdd(t *testing.T) {
	r := adder.Add(3, 5)
	if r != 8 {
		// We use t.Error() to report that something went wrong
		t.Error("incorrect result: expected 8, got", r)
		// You can also use t.Errorf for a Printf-style formatting
		// Both t.Error and t.Errorf mark the test as failed but the function keeps running
		// If you want it to stop, use t.Fatal or t.Fatalf
	}
}
