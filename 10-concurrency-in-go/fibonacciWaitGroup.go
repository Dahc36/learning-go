package main

import (
	"fmt"
	"sync"
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

func main() {
	defer fmt.Println("main - deferred")
	fmt.Println("main - start")

	values := []int{43, 41, 40, 39, 1, 8, 11, 42, 33, 22}
	results := make([]result, 0, len(values))

	var wg sync.WaitGroup
	wg.Add(len(values))
	for _, v := range values {
		go func(value int) {
			defer wg.Done()
			r := fibonacci(value)
			fmt.Printf("Fibonacci of %v is %v\n", value, r)
			results = append(results, result{value: value, result: r})
		}(v)
	}
	wg.Wait()

	fmt.Printf("%+v\n", results)
}
