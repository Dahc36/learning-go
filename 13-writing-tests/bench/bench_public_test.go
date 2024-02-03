package bench_test

import (
	"fmt"
	"testing"

	"github.com/dahc36/learning-go/13-writing-tests/bench"
)

func TestFileLen(t *testing.T) {
	result, err := bench.FileLen("testdata/data.txt", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result != 6726 {
		t.Error("Expected 6726, got", result)
	}
}

var blackHole int

func BenchmarkFileLen1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result, err := bench.FileLen("testdata/data.txt", 1)
		if err != nil {
			b.Fatal(err)
		}
		// We add an assignment of the result to a package-level variable to prevent
		// the compiler from optimizing away the call to FileLen
		blackHole = result
	}
}

func BenchmarkFileLen(b *testing.B) {
	for _, v := range []int{1, 10, 100, 1000, 10000, 100000, 1000000} {
		b.Run(fmt.Sprintf("FileLen-%d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result, err := bench.FileLen("testdata/data.txt", v)
				if err != nil {
					b.Fatal(err)
				}
				// We add an assignment of the result to a package-level variable to prevent
				// the compiler from optimizing away the call to FileLen
				blackHole = result
			}
		})
	}

	for _, v := range []int{1, 10, 100, 1000, 10000, 100000, 1000000} {
		b.Run(fmt.Sprintf("FileLen2-%d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result, err := bench.FileLen2("testdata/data.txt", v)
				if err != nil {
					b.Fatal(err)
				}
				// We add an assignment of the result to a package-level variable to prevent
				// the compiler from optimizing away the call to FileLen
				blackHole = result
			}
		})
	}
}
