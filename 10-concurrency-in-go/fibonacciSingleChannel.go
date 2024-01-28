package main

import (
	"fmt"
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

func processAsync(values []int, process func(v int) int, callback func(r result)) {
	ch := make(chan result)
	defer close(ch)
	f := func(v int) {
		ch <- result{value: v, result: process(v)}
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
