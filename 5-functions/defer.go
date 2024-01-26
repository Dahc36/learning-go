package main

import "fmt"

func deferredValueInnerScope() int {
	x := 1
	defer func() {
		x = 2
	}()

	return x
}
func deferredValueReturnScope() (x int) {
	x = 1
	defer func() {
		x = 2
	}()

	return x
}
func deferredValueParam(x int) int {
	x = 1
	defer func() {
		x = 2
	}()

	return x
}
func deferredValueParamReference(x *int) int {
	*x = 1
	defer func() {
		*x = 2
	}()

	return *x
}

var y = 0

func deferredValueFile() int {
	y = 1
	defer func() {
		y = 2
	}()

	return y
}

func main() {
	x := 0
	fmt.Println(deferredValueInnerScope())       // defer doesn't change the output
	fmt.Println(deferredValueReturnScope())      // returns 2
	fmt.Println(deferredValueParam(x))           // defer doesn't change the output
	fmt.Println("x:", x)                         // defer doesn't change x (because it gets a copy)
	fmt.Println(deferredValueParamReference(&x)) // defer doesn't change the output
	fmt.Println("x:", x)                         // defer does change the value to 2
	fmt.Println(deferredValueFile())             // defer doesn't change the output
	fmt.Println("y:", y)                         // defer does change the value to 2
}
