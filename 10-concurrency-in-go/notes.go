package main

import "fmt"

func main() {
	// Concurrency is the term for breaking up a single process into independent
	// component and specifying how these components safely share data.

	// Concurrency is not free!
	// Only use it when slow operations (like I/O) can run independently.
	// If you're not sure, write benchmarks

	fmt.Println("-- Goroutines --")
	// Definitions:
	// - Process
	//   - Instance of a program being run by the computers operating system
	//   - The OS assigns resources (like memory) to each process
	//   - Each process is made up of one or more threads
	// - Threads
	//   - Unit of execution that is given time to run by the OS
	//   - Threads within a process share access to resources
	//   - The CPU can execute instructions from one or more threads at the same time
	//     (depending on the number of cores)
	//   - One of the jobs of the OS is to schedule threads on the CPU so every process
	//     and every thread has a chance to run
	// - Goroutines
	//   - Lightweight processes managed by the Go runtime
	//   - When a program starts, the Go runtime creates a number of threads and launches
	//     a single goroutine to run the program
	//   - All goroutines created are assigned to these treads by the Fo runtime scheduler,
	//     just like the OS schedules threads across CPU cores
	//   - Goroutine creation is faster than thread creation, because no OS-level resource
	//     is being created
	//   - Goroutines' initial stack sizes are smaller than threads and can grow as needed
	//   - Switching between goroutines is faster than switching threads because it happens
	//     withing the same process. Avoiding OS calls
	//   - The scheduler can optimize decisions because it's part of the Go process
	//   - Go programs can spawn thousands of simultaneous goroutines

	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		v := 1
		ch1 <- v
		v2 := <-ch2
		// This is never reached for some reason?
		fmt.Println("From the sub-goroutine", v, v2)
	}()
	v := 2
	var v2 int
	// This select prevents a race condition where both goroutines are waiting for a read
	select {
	case ch2 <- v:
	case v2 = <-ch1:
	}
	fmt.Println("From the main goroutine", v, v2)

	// Keep your APIs concurrency-free, exposing channels or mutexes puts the responsibility
	// of managing them to the users of you API.
	// That doesn't mean that you can't use them as function parameters or struct fields,
	// just don't export those values.
	// The exception is when your API is a library with concurrency helper funcs (like time.After)

	a := []int{2, 4, 6, 8, 10}
	ch := make(chan int, len(a))
	for _, v :=d// Todo: continue with page 213
}
