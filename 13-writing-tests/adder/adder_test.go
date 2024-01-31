package adder

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var testTime time.Time

// If there is a TestMain, it is run instead of the individual tests
// It runs only once (not before and after each test)
func TestMain(m *testing.M) {
	fmt.Println("Set up stuff for tests here")
	testTime = time.Now()
	exitVal := m.Run()
	fmt.Println("Clean up stuff after tests here")
	os.Exit(exitVal)
}

func Test_addNumbers(t *testing.T) {
	result := addNumbers(2, 3)
	now := time.Now()
	fmt.Println("Time taken:", now.Sub(testTime))
	if result != 5 {
		t.Error("incorrect result: expected 5, got", result)
	}
}
