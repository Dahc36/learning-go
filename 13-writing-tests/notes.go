package main

import (
	"fmt"

	"github.com/dahc36/learning-go/13-writing-tests/race"
)

func main() {
	// Test file names have to end with _test.go
	fmt.Println("Hello 13")
	for range [10]int{} {
		counter := race.GetCounter()
		fmt.Println(counter)
	}
	// For performance profiling check out the blog https://oreil.ly/HHe9c
}
