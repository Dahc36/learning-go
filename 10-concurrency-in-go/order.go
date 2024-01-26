package main

import "fmt"

func foo1() {
	ch1 := make(chan int)
	go func() {
		fmt.Println("foo1:sub - Reading from ch1") // 2
		x := <-ch1
		fmt.Println("foo1:sub - Writing to ch1") // 3
		ch1 <- 2
		fmt.Println("foo1:sub - x:", x) // Reached eventually, sometimes
	}()

	fmt.Println("foo1:main - Writing to ch1") // 1
	ch1 <- 1
	fmt.Println("foo1:main - Reading from ch1") // 4
	x := <-ch1
	fmt.Println("foo1:main - x:", x) // 5
}

func foo2() {
	ch1 := make(chan int)
	go func() {
		fmt.Println("foo2:sub - Writing to ch1") // 3
		ch1 <- 2
		fmt.Println("foo2:sub - Reading from ch1") // 2
		x := <-ch1
		fmt.Println("foo2:sub - x:", x) // Reached eventually, sometimes
	}()

	fmt.Println("foo2:main - Reading from ch1") // 4
	x := <-ch1
	fmt.Println("foo2:main - Writing to ch1") // 1
	ch1 <- 1
	fmt.Println("foo2:main - x:", x) // 5
}

func foo3() {
	ch1 := make(chan int)
	go func() {
		fmt.Println("foo3:sub - Writing to ch1") // 3
		ch1 <- 2
		fmt.Println("foo3:sub - Writing to ch1 again") // 2
		ch1 <- 3
		fmt.Println("foo3:sub - Writing to ch1 again, with 4") // 2
		ch1 <- 4
		fmt.Println("foo3:sub - Writing to ch1 again, with 5") // 2
		ch1 <- 5
	}()

	fmt.Println("foo3:main - Reading from ch1") // 4
	x := <-ch1
	fmt.Println("foo3:main - x:", x) // 5
	x = <-ch1
	fmt.Println("foo3:main - x:", x) // 5
}

// Looks like when a goroutine already wrote to a channel and another one reads the channel's buffer wasn't empty
// If you read first and then someone else writes
func main() {
	fmt.Println("-- Foo --")
	foo1()
	fmt.Println("-- Foo 2 --")
	foo2()
	fmt.Println("-- Foo 3 --")
	foo3()
}
