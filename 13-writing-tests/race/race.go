package race

import "sync"

// The -race flag for the test command (also available when building non-test code)
// detects race conditions and prints a report to the screen.
// It makes the compiled code a lot slower though, so it should not be used by default
// Test by running either:
// - go test ./13-writing-tests/race -race -count=1
// - go run ./13-writing-tests -race
func GetCounter() int {
	var counter int
	var wg sync.WaitGroup
	// Comment out mx to get the race issue
	var mx sync.Mutex
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			for i := 0; i < 1000; i++ {
				mx.Lock()
				counter++
				mx.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return counter
}
