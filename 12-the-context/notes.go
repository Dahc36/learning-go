package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"time"
)

func main() {
	// In Go, the way to handle metadata on individual requests is the context.
	// This metadata can be required to correctly process the request or when to stop processing the request.
	// Context is an instance that meets the context.Context interface.
	// By convention the context should be passed as the first parameter of a function (as ctx)
	// ctx := context.Background()

	// The context was added to Go after the net/http package, so for backwards compatibility
	// it was added as part of new methods on http.Request:
	// - Context returns the context.Context associated with the request
	// - WithContext takes in a context.Context and returns a new request with the old request's
	//   state combined with the supplied context.Context
	// Here's the general pattern:
	// func middleware(handler http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		ctx := r.Context()
	// 		// wrap the context
	// 		r = r.WithContext(ctx)
	// 		handler.ServeHTTP(w, r)
	// 	})
	// }
	// You can also use the context when making an HTTP call, to pass a context through middleware
	// req, err := http.NewRequest(htt.MethodGet, "http://example.com", nil)
	// req = req.WithContext(ctx)

	fmt.Println("-- Cancellation --")
	// The Context helps to stop ongoing processes from completing when they are not necessary,
	// this is called cancellation. Use context.WithCancel to create a new context and a context.CancelFunc
	ss := slowServer()
	defer ss.Close()
	fs := fastServer()
	defer fs.Close()

	ctx := context.Background()
	callBoth(ctx, os.Args[1], ss.URL, fs.URL)

	// You can use time-limited contexts to control how long requests can take
	ctx = context.Background()
	parent, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	child, cancel2 := context.WithTimeout(parent, 3*time.Second)
	defer cancel2()
	start := time.Now()
	<-child.Done()
	end := time.Now()
	fmt.Println(end.Sub(start)) // The child is done in 2 seconds because the parent is done

	// Most of the time you don't need to worry about timeouts or cancellation
	// See pages 264 to 265 for the pattern for supporting context cancellation

	fmt.Println("-- Values --")
	// By default, you should prefer to pass data through explicit function parameters.
	// However, there are some cases where you cannot pass data explicitly. The most common is
	// an HTTP request handler and its associated middleware. Some examples are:
	// - Extracting a user from a JWT
	// - Creating  a per-request GUID that is passed through multiple layers of middleware
	//   and into your handler and business logic
	// See pages 266 to 270 for examples of both, there's a nice implementation of a decorator
	// that uses context to log out extra data
}

func slowServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This implementation would make the main program wait until the server closes
		// time.Sleep(6 * time.Second)
		// w.Write([]byte("Slow response"))
		// But we can use the context to close early
		ctx := r.Context()
		select {
		case <-ctx.Done():
			fmt.Println("server shut down")
			return
		case <-time.After(6 * time.Second):
			w.Write([]byte("Slow response"))
		}
	}))
	return s
}

func fastServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("error") == "true" {
			w.Write([]byte("error"))
			return
		}
		w.Write([]byte("ok"))
	}))
	return s
}

func callServer(ctx context.Context, label string, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(label, "request err:", err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(label, "response err:", err)
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(label, "read err:", err)
		return err
	}
	result := string(data)
	if result != "" {
		fmt.Println(label, "result:", result)
	}
	if result == "error" {
		fmt.Println("cancelling from", label)
		return errors.New("error happened")
	}
	return nil
}

var client = http.Client{}

func callBoth(ctx context.Context, errVal string, slowURL string, fastURL string) {
	ctx, cancel := context.WithCancel(ctx)
	// You must always cancel a cancelable context
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := callServer(ctx, "slow", slowURL)
		if err != nil {
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		err := callServer(ctx, "fast", fastURL+"?error="+errVal)
		if err != nil {
			cancel()
		}
	}()
	wg.Wait()
	fmt.Println("done with both")
}
