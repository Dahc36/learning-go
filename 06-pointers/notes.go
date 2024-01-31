package main

import (
	"fmt"
)

func main() {
	// Pointers are stored using 4 bytes of memory, regardless of the type they point to
	// & is the address operator, returns the address of the memory location
	// * is the indirection operator
	//   - it returns the pointer-to value, when used with a pointer
	//   - it returns a type that represents a pointer, when used with a type
	//
	var x int32 = 10
	var y bool = true
	pointerX := &x
	pointerY := &y
	var pointerZ *string
	fmt.Println(&pointerX, pointerX, *pointerX)
	fmt.Println(&pointerY, pointerY, *pointerY)
	fmt.Println(&pointerZ, pointerZ) // *pointerZ panics
	// See pages 107 and 108 for diagrams

	// You can use the built-in new function to create a pointer variable
	var np = new(int)
	fmt.Println(&np, np, *np)
	// It returns a zero value instance of the type
	// This is rarely used

	// You can also get a pointer using an & before a struct literal
	type fooS struct{ value int }
	fSp := &fooS{}
	fmt.Println(&fSp, fSp, *fSp)
	// You can't do this with primitive literals (numbers, booleans or strings)
	// Instead you have to declare a variable and point to it:
	var s string
	sp := &s
	fmt.Println(&sp, sp, *sp)
	// Or, you can use a helper function
	sp = stringp("")
	fmt.Println(&sp, sp, *sp)

	fmt.Println("-- Mutable parameters --")
	// Pointers in Go are used to indicate mutability
	// When you pass primitives, structs or arrays to a function those are passed by value
	// So they cannot be modified by the function.
	// However, when you pass a pointer it can be mutated by dereferencing the pointer,
	// to access the address where the value is stored.
	// Reassigning the pointer parameter itself won't work because it is also passed by value
	var ip = new(int)
	fmt.Println(&ip, ip, *ip)
	failsMutation(ip)
	fmt.Println(&ip, ip, *ip)
	doesMutation(ip)
	fmt.Println(&ip, ip, *ip)

	// Pointers should be a last resort, they:
	// - make it harder to understand data flow
	// - create extra work for the garbage collector
	// The only time you should use them is when the function expects an interface,
	// this is common when working with JSON

	// When returning values from a function, favor value types.
	// Pointers should only be returned when the data type needs to be modified,
	// an example for this is working with buffers for reading or writing data.
	// There are also data types to work with concurrency that must always be passed as pointers

	// Passing a pointer as parameter can improve performance when a struct gets too large,
	// Since the pointer has a fixed size
	// See more details in page 118.

	// Avoid passing maps as parameters or return values
	// - They say nothing about the values contained within
	// - You have to trace trough the code to figure out it's vales
	// - Makes the code not self-documented
	// - It's worse for garbage collection

	// Slices behavior is pretty weird, because they are defined as a struct of:
	// 1. A pointer to the underlying array
	// 2. An int field for the length
	// 3. An int field for the capacity
	// This leads to:
	// - Functions can modify the contents of the array but not the values for the length or capacity
	// - When they go over the length, the change will be in memory but the original array
	//   won't know that it's supposed to read that far
	// - Then they go over the capacity, the pointer for the copy now points at a different
	//   and bigger memory block

	// There is one use for slices as parameters
	// Using slices of bytes as buffers, we can prevent some unnecessary garbage collection
	// More on page 123

	// Prioritize mechanical sympathy with writing Go. Read page 123 trough the end of the chapter.
	// TLDR: we want to store as much data as possible in the stack, avoiding the heap
	// pointers are stored in the stack (because their size is known at compile time),
	// but their values are stored in the heap. Prefer non-pointers values/params whenever possible.
}

func stringp(s string) *string {
	return &s
}

func failsMutation(ip *int) {
	x := 10
	ip = &x
}

func doesMutation(ip *int) {
	*ip = 10
}
