package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"
)

func errorOnOne(i int) error {
	if i == 1 {
		return errors.New("value is one")
	}
	return nil
}

func errorWithFancyMessage(i int) error {
	if i == 1 {
		return fmt.Errorf("%d, isn't different to one", i)
	}
	return nil
}

// Errors are values that implement the error interface
type Status int

const (
	InvalidLogin Status = iota + 1
	NotFound
)

type StatusErr struct {
	Status  Status
	Message string
	Err     error
}

func (se StatusErr) Error() string {
	return se.Message
}

// We implement Unwrap to be able to wrap errors with our custom error
func (se StatusErr) Unwrap() error {
	return se.Err
}

// We can implement our own Is method for error.Is
func (se StatusErr) Is(target error) bool {
	if statusErr, ok := target.(StatusErr); ok {
		return se.Status == statusErr.Status
	}
	return false
}

// Don't define the returned error as the custom StatusErr, because that won't ever == nil
// That's because the zero value of a struct is not nil
func specialError(i int) error {
	if i == 1 {
		return StatusErr{
			Status:  InvalidLogin,
			Message: "Malo el login",
			Err:     errors.New("login sub-error"),
		}
	}
	if i == 2 {
		return StatusErr{
			Status:  NotFound,
			Message: "Malo el found",
			Err:     errors.New("not found sub-error"),
		}
	}
	return nil
}

// fmt.Errorf has a special verb %w, the new error will include the formatted
// string of another error and wrap that error as well.
// All errors wrapped this way make up the error chain.
func fileChecker(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("in fileChecker: %w", err)
	}
	f.Close()
	return nil
}

func main() {
	// Go handles errors by returning a value of type error as the last return
	// value for a function. This is a convention that should never be breached.
	// The calling code can check for errors by comparing the error with nil
	err := errorOnOne(1)
	if err != nil {
		// fmt.Println calls the Error method when it receives an error
		fmt.Println(err)
		fmt.Println("Error on one")
	}

	// error is a built-in interface that defines a single method
	// type error interface{
	// 	Error() string
	// }
	// We use nil to indicate no error, because that's the zero value for all interfaces

	// You can format the error messages with fmt.Errorf
	err = errorWithFancyMessage(1)
	if err != nil {
		// fmt.Println calls the Error method when it receives an error
		fmt.Println(err)
		fmt.Println("Fancy error on one")
	}

	fmt.Println("-- Sentinel Errors --")
	// Sentinel errors are used to indicate a problem with the current state
	// They are declared at the package level and their names should start with Err
	data := []byte("This is not a zip file")
	notAZipFile := bytes.NewReader(data)
	_, err = zip.NewReader(notAZipFile, int64(len(data)))
	if err == zip.ErrFormat {
		fmt.Println("Sentinel error thrown")
	}
	// Sentinel errors are the right choice when the application has reached a state where
	// no further processing is possible and no contextual information needs to be used
	// to explain the error state

	fmt.Println("-- Custom and Wrapped Errors --")
	// Errors are values that implement the error interface
	err = specialError(1)
	if err != nil {
		// fmt.Println calls the Error method when it receives an error
		fmt.Println(err)
	}

	// You can unwrap errors (if there's no wrapped error nil is returned)
	err = fileChecker("not_here.txt")
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Unwrap(err))
	}
	// Usually you don't call errors.Unwrap directly, instead you use errors.Is and errors.As

	fmt.Println("-- error.Is and error.As --")
	// error.Is can be used to check for a specific instance in the error chain
	err = fileChecker("not_here")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("That file doesn't exist")
		}
	}
	// It can also be implemented with custom logic
	err = specialError(2)
	if err != nil {
		if errors.Is(err, StatusErr{Status: NotFound}) {
			fmt.Println("Error is NotFound!")
		}
		if errors.Is(err, StatusErr{Status: InvalidLogin}) {
			fmt.Println("Error is InvalidLogin")
		} else {
			fmt.Println("Error is not InvalidLogin")
		}
	}

	// error.As can be used to check an error against a specific type
	err = specialError(1)
	var myErr StatusErr
	if errors.As(err, &myErr) {
		fmt.Println("Is is a StatusErr")
	}

	// Go provides a recover() function to gracefully handle panics
	// But it's not recommended as idiomatic go should prefer explicit errors
	// See pages 174 to 176.
}
