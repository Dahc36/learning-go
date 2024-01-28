package main

import (
	"fmt"
	"time"
)

func fibonacci(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

type result struct {
	value  int
	result int
}

func timeLimit(n int, f func() result) result {
	r := result{value: n, result: -1}
	done := make(chan struct{})
	go func() {
		r = f()
		close(done)
	}()
	select {
	case <-done:
		return r
	case <-time.After(500 * time.Millisecond):
		return r
	}
}

func processAsync(values []int, process func(v int) int, callback func(r result)) {
	ch := make(chan result)
	defer close(ch)
	f := func(v int) {
		ch <- timeLimit(v, func() result {
			return result{value: v, result: process(v)}
		})
	}

	for _, v := range values {
		go f(v)
	}
	for range values {
		callback(<-ch)
	}
}

func main() {
	defer fmt.Println("main - deferred")
	fmt.Println("main - start")

	values := []int{43, 41, 40, 39, 1, 8, 11, 42, 33, 22}
	results := make([]result, 0, len(values))
	processAsync(values, fibonacci, func(r result) {
		fmt.Printf("Fibonacci of %v is %v\n", r.value, r.result)
		results = append(results, r)
	})
	fmt.Printf("%+v\n", results)
}
