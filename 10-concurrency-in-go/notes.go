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

	// Goroutines communicate using channels
	// Channels are a reference type and their zero value is nil

	// Channels:
	// - Unbuffered (default)
	//   - When written to:
	//     - Open: the writing goroutine pauses until another one reads from the same channel
	//     - Closed: PANIC
	//   - When read from:
	//     - Open: the reading goroutine pauses until another one writes to the same channel
	//     - Closed: returns the zero value for the underlying channel type (use the
	//       comma-ok idiom to differentiate a closed channel read)
	//   - When closed:
	//     - Open: closes the channel
	//     - Closed: PANIC
	// - Buffered
	//   - When written to:
	//     - Open: doesn't pause until the buffer is full
	//     - Closed: PANIC
	//   - When read from:
	//     - Open: only pauses if the buffer is empty
	//     - Closed: returns the remaining values until empty, then the zero value
	//   - When closed:
	//     - Open: closes the channel but the remaining values are still there
	//     - Closed: PANIC
	// - nil
	//   - When written to: Hangs forever
	//   - When read from: Hangs forever
	//   - When closed: PANIC
	// To create a buffered channel do `make(chan int), i` where i is the buffer size

	// The select statement provides a way of solving conflicts when different goroutines
	// want to continue their operations at the same time.
	// Each case in a select is a read or a write to a channel, if a read or a write is possible
	// for any of the cases, it is executed along with the body of the case. If multiple cases
	// have channels that can be read or written, select chooses randomly one of them.

	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		v := 1
		// This writes to a channel
		ch1 <- v
		// This reads from a channel
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
	// You can close a channel with close(ch)

	// Keep your APIs concurrency-free, exposing channels or mutexes puts the responsibility
	// of managing them to the users of you API.
	// That doesn't mean that you can't use them as function parameters or struct fields,
	// just don't export those values.
	// The exception is when your API is a library with concurrency helper funcs (like time.After)

	// Every time your goroutine depends on a value that might change, pass the current value
	a := []int{2, 4, 6, 8, 10}
	ch := make(chan int, len(a))
	for _, v := range a {
		go func(val int) {
			ch <- val * 2
		}(v)
	}
	for i := 0; i < len(a); i++ {
		fmt.Println(<-ch)
	}

	// Always close your goroutines
	for i := range countTo(10) {
		fmt.Println(i)
	}

	// This one won't close the channel
	for i := range countTo(10) {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}

	// The done channel pattern provides a way to stop or cancel processing
	chCancel, cancel := countToWithCancel(10)
	for i := range chCancel {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}
	cancel()

	// You can turn off a case in a select, which can be useful when one of the channels is closed
	// You just have to make the channel = nil. Look into page 219 for example.

	// There's a utility type sync.Once, that ensures we only run some code once, see page 222.

	// ToDo: document for-range with channels
	// ToDo: document for-select
	// ToDo: implement input listener and return fibonacci
	// ToDo: implement parallel fetches
	// ToDo: implement map-reduce with benchmark to compare it to sequential processing
}

func countTo(max int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < max; i++ {
			ch <- i
		}
		fmt.Println("Closing channel")
		close(ch)
	}()
	return ch
}

func countToWithCancel(max int) (<-chan int, func()) {
	ch := make(chan int)
	done := make(chan struct{})
	cancel := func() {
		fmt.Println("Closing done")
		close(done)
	}
	go func() {
		for i := 0; i < max; i++ {
			select {
			case <-done:
				break
			case ch <- i:
			}
		}
		fmt.Println("Closing ch")
		close(ch)
	}()
	return ch, cancel
}
