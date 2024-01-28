package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

func main() {
	s := "the quick brown fox jumped over the lazy dog"
	sr := strings.NewReader(s)
	counts, err := countLetters(sr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("s has %+v letters\n", counts)

	r, closer, err := buildGZipReader("./11-the-standard-library/example.txt.gz")
	defer closer()
	if err != nil {
		log.Fatal(err)
	}
	counts, err = countLetters(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("s has %+v letters\n", counts)
	// todo: script to change folder names
}
