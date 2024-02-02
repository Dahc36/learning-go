package adder

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

var testTime time.Time

// When there's common state that you want to set up before any tests you can use TestMain
// If there is a TestMain, it is run instead of the individual tests
// It runs only once (not before and after each test)
func TestMain(m *testing.M) {
	fmt.Println("Set up stuff for tests here")
	testTime = time.Now()
	// You have to run the tests yourself calling m.Run()
	exitVal := m.Run()
	fmt.Println("Clean up stuff after tests here")
	// You have to call os.Exit with the result from m.Run()
	os.Exit(exitVal)
}

// Test functions start with the word Test and take a single parameter of type *testing.T
// The name of the test is supposed to document what it's testing
// When testing unexported functions you can use _ before the name of the function
func Test_addNumbers(t *testing.T) {
	result := addNumbers(2, 3)
	now := time.Now()
	fmt.Println("Time taken:", now.Sub(testTime))
	if result != 5 {
		// We use t.Error() to report that something went wrong
		t.Error("incorrect result: expected 5, got", result)
		// You can also use t.Errorf for a Printf-style formatting
		// Both t.Error and t.Errorf mark the test as failed but the function keeps running
		// If you want it to stop, use t.Fatal or t.Fatalf
	}
}

func TestSecond(t *testing.T) {
	fmt.Println("TestSecond also uses stuff set up in TestMain", testTime)
}

// For logic that needs cleanup, you can use t.Cleanup
func createFile(t *testing.T) (string, error) {
	f, err := os.Create("testdata/tempFile")
	if err != nil {
		return "", err
	}
	// write some data to f
	io.Copy(f, strings.NewReader("How you doin'?"))
	t.Cleanup(func() {
		os.Remove(f.Name())
	})
	return f.Name(), nil
}

func TestFileProcessing(t *testing.T) {
	fName, err := createFile(t)
	if err != nil {
		t.Fatal(err)
	}
	// test and don't worry about cleanup
	fmt.Println(fName)
}

func TestTestData(t *testing.T) {
	// You can use the testdata directory with relative path from within tests
	f, err := os.Open("testdata/data.txt")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(os.Stdout, f)

	t.Cleanup(func() {
		f.Close()
	})
}
