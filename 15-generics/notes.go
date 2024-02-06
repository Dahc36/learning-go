package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Stack[T comparable] struct {
	vals []T
}

func (s *Stack[T]) Push(val T) {
	s.vals = append(s.vals, val)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}
	top := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return top, true
}

func (s Stack[T]) Contains(val T) bool {
	for _, v := range s.vals {
		if v == val {
			return true
		}
	}
	return false
}

func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
	r := make([]T2, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func Reduce[T1, T2 any](s []T1, initializer T2, f func(T2, T1) T2) T2 {
	r := initializer
	for _, v := range s {
		r = f(r, v)
	}
	return r
}

func Filter[T1 any](s []T1, f func(T1) bool) []T1 {
	var r []T1
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

type BuiltInOrdered interface {
	// ~ before the type means that it will be valid for any type that has that underlying type
	string | ~int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint64 | uintptr
}

func Min[T BuiltInOrdered](v1, v2 T) T {
	if v1 < v2 {
		return v1
	}
	return v2
}

func BuiltInOrderable[T BuiltInOrdered](v1, v2 T) int {
	if v1 < v2 {
		return -1
	}
	if v1 > v2 {
		return 1
	}
	return 0
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Order(other Person) int {
	out := p.Age - other.Age
	if out == 0 {
		out = strings.Compare(p.Name, other.Name)
	}
	return out
}

func main() {
	// Go 1.18 introduced generics (type parameters), along with some utility types:
	// - any (same as interface{})
	// - comparable (all types that can be == or !=)
	var st Stack[int]
	st.Push(10)
	st.Push(20)
	st.Push(30)
	fmt.Println(st.Contains(10))
	fmt.Println(st.Contains(5))

	fmt.Println("-- Generic Functions --")
	// You can now write generic functions
	s := []int{1, 2, 3, 4, 5}
	fmt.Println(Map(s, func(i int) string {
		return "Value: " + strconv.FormatInt(int64(i*2), 10)
	}))
	fmt.Println(Reduce(s, 0, func(total int, v int) int {
		return total + v
	}))
	fmt.Println(Filter(s, func(v int) bool {
		return v%2 != 0
	}))

	fmt.Println("-- Generic Interfaces --")
	pair2Da := Pair[Point2D]{Point2D{1, 1}, Point2D{5, 5}}
	pair2Db := Pair[Point2D]{Point2D{1, 1}, Point2D{5, 5}}
	closer := FindCloserPair(pair2Da, pair2Db)
	fmt.Println(closer)

	pair3Da := Pair[Point3D]{Point3D{1, 1, 10}, Point3D{5, 5, 0}}
	pair3Db := Pair[Point3D]{Point3D{10, 10, 10}, Point3D{11, 5, 0}}
	closer2 := FindCloserPair(pair3Da, pair3Db)
	fmt.Println(closer2)

	fmt.Println("-- Type terms --")
	a := 10
	b := 20
	fmt.Println(Min(a, b))

	type myInt int
	var myA myInt = 0
	var myB myInt = 3
	// This wouldn't work without ~int in the BuiltInOrdered declaration
	fmt.Println(Min(myA, myB))

	// Go will infer the type whenever possible (when calling generic functions or types)
	// But there are situations where you have to specify it.
	// You can use the generic types within a function (like for type conversion T2(v1), see page 338).
	// Type terms will limit the operations and values that can be done with variables of the type.

	fmt.Println("-- Generic Functions with Generic Data Structures")
	t1 := NewTree(BuiltInOrderable[int])
	t1.Add(10)
	t1.Add(30)
	t1.Add(15)
	fmt.Println(t1.Contains(15))
	fmt.Println(t1.Contains(40))

	t2 := NewTree(Person.Order)
	t2.Add(Person{"David", 34})
	t2.Add(Person{"Marla", 34})
	t2.Add(Person{"Bob", 47})
	fmt.Println(t2.Contains(Person{"Marla", 34}))
	fmt.Println(t2.Contains(Person{"Fred", 20}))
}
