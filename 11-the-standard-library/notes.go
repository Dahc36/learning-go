package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func countLetters(r io.Reader) (map[string]int, error) {
	buf := make([]byte, 2048)
	out := map[string]int{}
	for {
		n, err := r.Read(buf)
		for _, b := range buf[:n] {
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func buildGZipReader(fileName string) (*gzip.Reader, func(), error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, nil, err
	}
	return gr, func() {
		gr.Close()
		r.Close()
	}, nil
}

type myStringReader string

// This is not how Read should be implemented, this is just a naive way of returning a string
// The correct way is to use p as a stream to return the value in chunks of len(p)
// See the actual strings.Reader.Read implementation.
func (s myStringReader) Read(p []byte) (n int, err error) {
	if len(p) < len(s) {
		return 0, errors.New("p is too small, make it bigger")
	}
	n = copy(p, s)
	return n, io.EOF
}

type myStringWriter struct{}

func (s myStringWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), io.EOF
}

func main() {
	fmt.Println("-- IO --")
	cl := func(r io.Reader) {
		counts, err := countLetters(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("s has %+v letters\n", counts)
	}
	s := "the quick brown fox jumped over the lazy dog"
	sr := strings.NewReader(s)
	cl(sr)

	r, closer, err := buildGZipReader("./11-the-standard-library/example.txt.gz")
	defer closer()
	if err != nil {
		log.Fatal(err)
	}
	cl(r)

	mr := myStringReader("Hello my baby, hello my honey, hello my ragtime gal")
	cl(mr)

	// io.Copy makes it trivial to take a Reader and pass it directly to a Writer
	io.Copy(myStringWriter{}, mr+"\n")
	f, err := os.Open("./11-the-standard-library/stan.txt")
	if err != nil {
		log.Fatal(err)
	}
	// If you're opening a resource in a loop, don't use defer because it won't run until the function exits
	// Instead call it before the end of each iteration
	defer f.Close()
	io.Copy(myStringWriter{}, f)

	// In general, try to create interfaces as simple and decoupled as the interfaces defined in io
	// They demonstrate the power of simple abstractions.

	fmt.Println("-- Time --")
	// time.duration
	d := 2*time.Hour + 30*time.Minute
	// time.duration implements the Stringer interface and returns a formatted duration
	fmt.Println(d)
	// Formatted durations can be turned into durations with the time.ParseDuration function
	d, err = time.ParseDuration("24h")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hours in a day:", d.Hours())
	fmt.Println("Minutes in a day:", d.Minutes())
	fmt.Println("Seconds in a day:", d.Seconds())

	// A moment in time is represented with time.Time, complete with a timezone
	// To parse time you have to provide a layout string that represents January 2nd 2006
	// at 03:04:05 MST, because it's 1-2-3-4-5-6-7 when represented like this:
	t, err := time.Parse("01/02 03:04:05PM '06 -0700", "11/16 03:30:00AM '89 -0300")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)
	fmt.Println(t.Format(time.RFC822))
	fmt.Println(t.Format("January 2, 2006 at 03:04PM MST"))
	// Don't use == to check for the same moments in time
	// Prefer the Equal method

	fmt.Println("-- JSON --")
	// Marshaling: converting from a Go data type to an encoding
	// Unmarshaling: means converting to a Go data type

	// We define the rules for processing our JSON with struct tags with the json tag name
	// Struct tags are composed of value pairs, written as tagName:"tagValue" and separated by spaces
	// Every field is exported so that the encoding/json package can access them
	type Item struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	type Order struct {
		Id          string    `json:"id"`
		DateOrdered time.Time `json:"date_ordered"`
		CustomerId  string    `json:"customer_id"`
		Items       []Item    `json:"items"`
	}
	// If a field should be ignored when marshaling or unmarshaling, use a - for the name
	// If it should be left out when empty, add omitempty after the name

	var o Order
	err = json.Unmarshal([]byte(`{"id": "1234", "customer_id": "000001", "items": []}`), &o)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", o)
	out, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	// How do Unmarshal and Marshal evaluate struct tags? How do they work with any struct type?
	// With reflection, check out chapter 14

	// JSON Readers and Writers
	// Working with readers, we could get []byte's of our data to use json.Unmarshal
	// but that's cumbersome, instead we can se json.Decoder and json.Encoder which
	// read and write to anything that implements the io.Reader and io.Writer interfaces
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := person{
		Name: "Fred",
		Age:  28,
	}
	// We create a temporary file in a temporary directory for this example
	tmpFile, err := os.CreateTemp(os.TempDir(), "sample-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	// We write to the temporary file
	err = json.NewEncoder(tmpFile).Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Now we can read from the file
	tmpFile2, err := os.Open(tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	var fromFile person
	err = json.NewDecoder(tmpFile2).Decode(&fromFile)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpFile2.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", fromFile)

	// JSON streams
	data := `{"name": "Old Fred", "age": 40}
{"name": "Mary", "age": 21}
{"name": "Pat", "age": 30}`
	sr = strings.NewReader(data)
	dec := json.NewDecoder(sr)
	p = person{}
	for dec.More() {
		err := dec.Decode(&p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", p)
	}
	// Read 246 for custom JSON parsing

	fmt.Println("-- net/http --")
	// Go's standard library includes a production quality HTTP/2 client and server
	// The net/http package defines a Client type to make HTTP requests and receive HTTP responses
	// There's a default client instance (DefaultClient) but don't use it as it defaults to having no timeout
	// Instead instantiate your own, you only need 1 http.Client for your entire program,
	// as it properly handles multiple simultaneous request across goroutines
	client := &http.Client{Timeout: 30 * time.Second}
	// To make a request, you create a new *http.Request instance
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-My-Client", "Learning Go")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatal(fmt.Errorf("unexpected status: got %v", res.StatusCode))
	}
	fmt.Println(res.Header.Get("Content-Type"))
	var resData struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", resData)

	// http.Server is responsible for listening for HTTP requests
	// ser := http.Server{
	// 	Addr:         ":8080",
	// 	ReadTimeout:  30 * time.Second,
	// 	WriteTimeout: 90 * time.Second,
	// 	IdleTimeout:  120 * time.Second,
	// 	Handler:      helloHandler{},
	// }
	// err = ser.ListenAndServe()
	// if err != nil {
	// 	if err != http.ErrServerClosed {
	// 		log.Fatal(err)
	// 	}
	// }

	// http.ServeMux is a request router that can be assigned to the Handler field in http.Server
	personM := http.NewServeMux()
	personM.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("greetings!\n"))
	})
	dogM := http.NewServeMux()
	dogM.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("woof!\n"))
	})
	mux := http.NewServeMux()
	mux.Handle("/person/", http.StripPrefix("/person", personM))
	mux.Handle("/dog/", http.StripPrefix("/dog", dogM))
	mux.Handle("/", terribleSecurityProvider("pass")(helloHandler{}))
	wrappedMux := requestTimer(mux)
	ser := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      wrappedMux,
	}
	err = ser.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
	// *http.ServeMux has some drawbacks:
	// - Doesn't allow you to specify handlers based on an HTTP verb or header
	// - Doesn't provide support for variables in the URL path
	// - Nesting can be clunky
	// Because of these limitations, there are third-party modules that aim at replacing it.
	// The book recommends gorilla mux and chi
}

type helloHandler struct{}

func (hh helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

// You can create a function that takes in a handler and returns a handler to implement the middleware pattern
// These middleware can be chained to augment any handler or other middleware
func requestTimer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		end := time.Now()
		log.Printf("request time for %s: %v", r.URL.Path, end.Sub(start))
	})
}

// Handlers wrapped with this have to be called like `curl --header "X-Secret-Password: <password>" localhost:8080`
func terribleSecurityProvider(password string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("You didn't give the secret password\n"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
