package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := longRunningThingManager(ctx, "What", 2*time.Microsecond)
		fmt.Println("result:", result, "- err:", err)
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("cancel")
	cancel()

	wg.Wait()
}

func longRunningThing(ctx context.Context, data string, duration time.Duration) (string, error) {
	fmt.Println("longRunningThing")
	time.Sleep(duration)
	return data, nil
}

func longRunningThingManager(ctx context.Context, data string, duration time.Duration) (string, error) {
	fmt.Println("longRunningThingManager")
	type resultWrapper struct {
		result string
		err    error
	}

	// The buffer where allows the goroutine to exit, even if no one reads the value
	ch := make(chan resultWrapper, 1)
	go func() {
		result, err := longRunningThing(ctx, data, duration)
		fmt.Println("sending result to channel")
		ch <- resultWrapper{result: result, err: err}
		fmt.Println("Done with the goroutine")
	}()

	select {
	case data := <-ch:
		return data.result, data.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
