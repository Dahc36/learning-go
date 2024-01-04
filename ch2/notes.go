package main

import "fmt"

func main() {
	// Literals can be written in different bases
	fmt.Println(0b01111111)
	fmt.Println(0b11111111)
	fmt.Println(0b10)
	fmt.Println(0o10)
	fmt.Println(0x10)
	// Using only a leading 0 is the same as 0o, but don't
	// Use decimal unless it makes sense for the use case

	// Untyped integer literals can be assigned to any number type
	const ci = 1_000
	var i int = ci
	var i32 int32 = ci
	var i64 int64 = ci
	var f float32 = ci
	var u uint = ci
	fmt.Println(i, i32, i64, f, u)

	// What integer to use
	// 1. If you're working with a format or protocol that requires a specific size / sign
	//    - Use that type
	// 2. If you're writing a function that should work with multiple types
	//    - Write a different function for each type
	// 3. Just use int

	// Avoid floats when possible because of their precision limitations
	// Avoid strict equality between floats (use epsilon instead)

	// Fun fact: Go has complex types (also avoid when possible)
	fmt.Println(complex(1, 1))

	// Rune type is the same as int32
	fmt.Println(int32('a'))
	// Avoid using int32 for characters though, to clarify intent

	// There's no implicit type conversion in Go
	fmt.Println(i + int(i32))
	fmt.Println(int64(i) + i64)
	fmt.Println(byte(i))

	// Always use :=, except when:
	// 0. (Try to avoid this) Declaring outside of a function (in the package scope)
	// 1. Initializing a variable with its zero value
	// 2. Assigning an untyped constant or a literal value,
	//    where the default type doesn't match what you want
	// 3. To avoid shadowing variables

	// Only declare variables in the same line when assigning multiple
	// values returned from a function or the comma ok idiom

	// Constants are just a way to give a name to literals
	// They can be typed or untyped (with default types, just like literals)

	// Use camelCase for your variables
	// An uppercase first letter is used to expose package-level declarations
	// Short variable names can help keep code short and concise
	// Use longer, more descriptive name for package blocks
}
