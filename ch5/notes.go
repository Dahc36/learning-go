package main

import (
	"errors"
	"fmt"
	"sort"
)

func main() {
	fmt.Println(div(5, 2))

	fmt.Println("-- Named and Optional parameters --")
	// Go doesn't have named or optional input parameters

	// You must supply all parameters for a function
	MyFunc(MyFuncOpts{
		LastName: "Patel", Age: 33,
	})
	MyFunc(MyFuncOpts{
		FirstName: "Freida", LastName: "Pinto",
	})
	// Functions should have few parameters
	// Prioritize simple functions over complex inputs

	fmt.Println("-- Variadic parameters --")
	// The variadic parameter must be the last (or only) function parameter
	fmt.Println(addTo(1, 2, 3, 4, 5))
	fmt.Println(addTo(11, 31))
	fmt.Println(addTo(-2, []int{4, 8, 15, 16, 23, 42}...))
	fmt.Println(addTo(3))

	fmt.Println("-- Multiple return values --")
	fmt.Println(divAndRemainder(5, 2))
	fmt.Println(divAndRemainder(4, 0))

	fmt.Println("-- Named return values --")
	// You can pre-declare variables to hold the return values
	fmt.Println(divAndRemainderNamed(5, 2))
	// This assignment are local to the function
	// You can only name some of the return values, using _
	// There's no requirement to actually return the named return values, which can be confusing

	// Functions are values!

	// You can declare function types
	type myFuncType func(int, int) int

	fmt.Println("-- Anonymous Functions --")
	for i := 0; i < 3; i++ {
		func(j int) {
			fmt.Println(j)
		}(i)
	}
	// This is not useful in this context (we could just remove the function and call the code directly)
	// But anonymous functions are useful in defer statements and Goroutines

	fmt.Println("-- Functions as parameters --")
	type Person struct {
		FirstName string
		LastName  string
		Age       int
	}
	people := []Person{
		{FirstName: "Tom", LastName: "Thomson", Age: 37},
		{FirstName: "Pat", LastName: "Patterson", Age: 23},
		{FirstName: "Rob", LastName: "Roberson", Age: 18},
	}
	fmt.Println(people)
	sort.Slice(people, func(i int, j int) bool {
		return people[i].LastName < people[j].LastName
	})
	fmt.Println(people)
	sort.Slice(people, func(i int, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Println(people)

	// The defer keyword ensures that some code runs, no matter how we exit a function
	// You can defer multiple calls, they run in LIFO order
	// See example in page 102 to understand using named return values with defer
	// TLDR: you can do clean-up with logic based on the final return values (like the error)
	// defer for clean-up leads to easier to read code by avoiding the nesting in try/catch/finally blocks

	fmt.Println("-- Go is call by value --")
	i := 1
	str := "Hello"
	p := person{name: "Robert"}
	modifyFails(i, str, p)
	fmt.Println(i, str, p)
	// Primitives and structs are not modified by passing them as parameters

	fmt.Println("-- Go is call by value - With pointers --")
	m := map[int]string{
		1: "first",
		2: "second",
	}
	modMap(m)
	fmt.Println(m)

	s := []int{1, 2, 3}
	modSlice(s)
	fmt.Println(s)
	// Maps and Slices can be modified by passing them as parameters
	// because they are implemented with pointers
}

func div(numerator int, denominator int) int {
	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

type MyFuncOpts struct {
	FirstName string
	LastName  string
	Age       int
}

func MyFunc(opts MyFuncOpts) {
	fmt.Printf("%+v\n", opts)
}

func addTo(base int, vals ...int) []int {
	out := make([]int, 0, len(vals))
	for _, val := range vals {
		out = append(out, base+val)
	}

	return out
}

func divAndRemainder(numerator int, denominator int) (int, int, error) {
	if denominator == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}

	return numerator / denominator, numerator % denominator, nil
}

func divAndRemainderNamed(numerator int, denominator int) (result int, remainder int, _ error) {
	if denominator == 0 {
		return result, remainder, errors.New("cannot divide by zero")
	}

	result = numerator / denominator
	remainder = numerator % denominator
	return result, remainder, nil
	// You could just do:
	// return
	// But don't, ever. It requires scanning back to understand what's actually returned
}

type person struct {
	age  int
	name string
}

func modifyFails(i int, s string, p person) {
	i = i * 2
	s = "Goodbye"
	p.name = "Bob"
}

func modMap(m map[int]string) {
	m[2] = "Hello"
	m[3] = "Goodbye"
	delete(m, 1)
}

func modSlice(s []int) {
	for k, v := range s {
		s[k] = v * 2
	}
	s = append(s, 10)
}
