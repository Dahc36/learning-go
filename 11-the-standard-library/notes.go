package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
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
	// If a field should be ignored when marshalling or unmarshaling, use a - for the name
	// If it should be left out when empty, add omitempty after the name

	// todo: script to change folder names
}
