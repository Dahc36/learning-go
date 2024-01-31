package main

import "fmt"

func main() {
	fmt.Println("-- Arrays --")
	// Arrays can be created in multiple ways
	var az [3]int
	fmt.Println(az)
	var al = [3]int{1, 2, 3}
	fmt.Println(al)
	var as = [7]int{1, 3: 2, -1, 6: 3}
	fmt.Println(as)

	// When using a literal you can ignore the length
	var al2 = [...]int{1, 2, 3}
	// You can compare for equality
	fmt.Println(al == al2)

	// The built-in function len returns the array's length
	fmt.Println(len([...]int{42}))

	// Only use arrays when you know the exact length ahead of time
	// And it makes things clearer
	// This is an exception, the rule is to use slices

	fmt.Println("-- Slices --")
	// With slices we don't specify the length
	var s = []int{0, 1, 2, 3}
	fmt.Println(s)
	// Otherwise, declaring slices works very similarly to arrays

	// Differences arise when declaring without a literal
	var sz []int
	fmt.Println(sz, sz == nil) // true
	// The zero value for a slice is nil
	// If you use an empty literal the slice won't equal nil
	var sz2 = []int{}
	fmt.Println(sz2, sz2 == nil) // false

	// Comparing slices is a compile-time error, you can only compare slices with nil
	// fmt.Println(sz == s) // errors

	// len (and other built-in functions) works with slices too
	fmt.Println(len(s))

	// You can add items with append and check the capacity with cap
	fmt.Println(s, len(s), cap(s))
	s = append(s, 4)
	fmt.Println(s, len(s), cap(s))
	s = append(s, 5, 6)
	fmt.Println(s, len(s), cap(s))
	s2 := []int{23, 42, 101}
	s = append(s, s2...)
	fmt.Println(s, len(s), cap(s))

	// The Go runtime is compiled into every Go binary

	// If you now the size of the slice ahead of time, declare it with make
	sm := make([]int, 5)
	fmt.Println(sm, len(sm), cap(sm))
	// You can specify both the length and the capacity
	sm = make([]int, 5, 10)
	fmt.Println(sm, len(sm), cap(sm))

	// How to declare slices
	// 1. If you have initial values or the values won't change, use a slice literal
	// 2. If you know the size, use make and index into the corresponding values
	// 3. Other situations can create a 0 length with specified capacity,
	//    using append to add values

	// You can create slices of slices with a slice expression
	fmt.Println(s[:])
	fmt.Println(s[:2])
	fmt.Println(s[8:])
	fmt.Println(s[1:4])
	// Slices created this way share memory
	ss := s[:2]
	ss[1] = -1
	fmt.Println(s)
	fmt.Println(ss)
	// Things get even weirder with append
	ss = append(ss, -2)
	fmt.Println(s)
	fmt.Println(ss)
	// You can avoid the append with three-part slices but not the memory sharing
	// I don't fully understand why, read page 45 for more info
	ss = s[:2:2]
	ss = append(ss, -3)
	fmt.Println(s)
	fmt.Println(ss)
	// Arrays can be sliced (which produces a new slice),
	// but they have the same memory sharing issues

	// You can use copy to prevent all of the problems
	fmt.Println("-- Slices.Copy --")
	s = []int{1, 2, 3, 4}
	ss = make([]int, 2)
	sss := make([]int, 3)
	copy(ss, s)
	copy(sss, s[1:])
	ss[1] = -1
	sss[1] = -2
	fmt.Println(s)
	fmt.Println(ss)
	fmt.Println(sss)

	// You can slice strings and index into them but beware of them being sequences of bytes
	// So if a character is represented by more than one byte weird things can happen
	st := "Smiling sun: 🌞"
	fmt.Println(st)
	fmt.Println(st[len(st)-1:])
	fmt.Println(st[len(st)-3:])
	fmt.Println(st[len(st)-4:])
	// It takes 4 bytes to encode the smiling sun emoji in UTF-8

	// Beware of converting int to string, you can get chars
	// fmt.Println(string(65)) // Shows warning in recent Go versions

	// You can convert strings to slices of bytes or runes
	fmt.Println([]byte(st))
	fmt.Println([]rune(st))

	// Instead of all this slicing strings:
	// You should use the functions provided in the strings and unicode/utf8 packages
	// Or for-range loops

	fmt.Println("-- Maps --")
	// Maps can store any one type of data associated with keys of any comparable value
	var nilMap map[string]int
	fmt.Println(nilMap, nilMap == nil, len(nilMap))
	nilMap2 := map[string]int{}
	fmt.Println(nilMap2, nilMap2 == nil, len(nilMap2))

	teams := map[string][]string{
		"Orcas":   {"Fred", "Ralph", "Bijou"},
		"Lions":   {"Sarah", "Peter", "Billie"},
		"Kittens": {"Waldo", "Raul", "Ze"},
	}
	fmt.Println(teams, teams == nil, len(teams))

	// You can use make to specify an initial size
	ages := make(map[int][]string, 10)
	fmt.Println(ages, ages == nil, len(ages))

	// The key for a map has to be a comparable type
	// So it cannot be a slice or another map

	// You can't use := to assign to a map key
	totalWins := map[string]int{}
	totalWins["Orcas"] = 1
	totalWins["Lions"] = 2
	fmt.Println(totalWins["Orcas"])
	fmt.Println(totalWins["Lions"])
	fmt.Println(totalWins["Kittens"])
	totalWins["Kittens"]++
	fmt.Println(totalWins["Kittens"])
	// Maps return the zero value for their the type when asking for a non-existing key

	// The comma ok idiom allows differentiation between existing and non-existing keys
	m := map[string]int{
		"hello": 5,
		"world": 0,
	}
	v, ok := m["hello"]
	fmt.Println(v, ok)
	v, ok = m["world"]
	fmt.Println(v, ok)
	v, ok = m["goodbye"]
	fmt.Println(v, ok)

	// You can delete with delete
	delete(m, "hello")
	v, ok = m["hello"]
	fmt.Println(v, ok)

	// There are no Sets in Go, but you can use a map for a similar effect
	intSet := map[int]bool{}
	vals := []int{5, 10, 2, 5, 8, 7, 3, 9, 1, 2, 10}
	for _, v := range vals {
		intSet[v] = true
	}
	fmt.Println(len(vals), len(intSet))
	fmt.Println(intSet[5])
	fmt.Println(intSet[500])
	// You can use empty structs instead of booleans, because they use zero bytes
	// However, you rely on the comma ok idiom to check for the existence of values
	// This is probably not necessary unless you have a lot of values

	fmt.Println("-- Structs --")
	type person struct {
		name string
		age  int
		pet  string
	}

	// The zero value sets the zero value for each field
	var fred person
	fmt.Printf("%+v\n", fred)

	// You can also use a struct literal
	bob := person{}
	fmt.Printf("%+v\n", bob)
	// Unlike with maps, there's no difference when declaring with literal and not assigning

	// For non-empty struct literals, you can use a list of values
	julia := person{
		"Julia", 36, "Cat",
	}
	fmt.Printf("%+v\n", julia)
	// For this style you have to define every value
	// and they have to follow the definition order
	// Another option is
	beth := person{
		age:  30,
		name: "Beth",
	}
	fmt.Printf("%+v\n", beth)
	// Both styles cannot be mixed together

	// Struct fields are accessed with dotted notation
	fmt.Println(beth.name)
	fmt.Println(julia.pet)

	// Maps use brackets and structs use dots

	// Anonymous structs can be used to declare variables without defining a struct first
	var specialPerson struct {
		name string
		age  int
		pet  string
	}
	specialPerson.name = "Bob"
	specialPerson.age = 3
	specialPerson.pet = "Dog"
	fmt.Printf("%+v\n", specialPerson)

	pet := struct {
		name string
		kind string
	}{name: "Fido", kind: "German Shepherd"}
	fmt.Printf("%+v\n", pet)

	// Structs are comparable if every one of its fields is comparable.
	// You can't compare different struct types, unless:
	// - they have the same fields in the same order and you convert one into the other
	// - they have the same fields in the same order and one of them is anonymous
}
