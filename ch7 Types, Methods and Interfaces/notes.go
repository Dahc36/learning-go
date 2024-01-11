package main

import (
	"fmt"
	"time"
)

// You can use any primitive or compound type to define a concrete type
type Score int
type Converter func(string) Score
type TeamScores map[string]Score
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// Abstract type: specifies what a type should do, but not how
// Concrete type: specifies what and how (provides ways to store data and implementation for any method declared)

// Methods for a type are defined at the package block level
func (p Person) String() string {
	return fmt.Sprintf("%s %s, age %d", p.FirstName, p.LastName, p.Age)
}

// The receiver specification differentiates between a regular function and a method
// Don't use 'this' or 'self', just abbreviate the type

// Methods can only be declared in the same package as their associated type
// You can still create methods in different files within the same package, but don't

// Just like with parameters, receivers can be values or pointers
// 1. If your method modifies the receiver, you must use a pointer receiver
// 2. If your method needs to handle nil instances, you must use a pointer receiver
// 3. If your method doesn't modify the receiver, you can use a value receiver
//    - It's common practice that if the type has other pointer receiver methods,
//      have all of its methods use pointer receivers, for consistency

type Counter struct {
	total       int
	lastUpdated time.Time
}

func (c *Counter) Increment() {
	c.total++
	c.lastUpdated = time.Now()
}

func (c Counter) String() string {
	return fmt.Sprintf("total: %d, last updated: %v", c.total, c.lastUpdated)
}

// Go considers both pointer and value receiver methods to be in the method set of a pointer instance.
// For a value instance, only the value receiver methods are in the method set.

// Don't write getters and setters for Go structs, unless:
// - You need to update multiple fields as a single operation
// - The update isn't a straight forward assignment of a new value

func (p *Person) StringNilInstance() string {
	if p == nil {
		return "No person found"
	}

	return p.String()
}

// You can use a method as a replacement for a function
type Adder struct {
	start int
}

func (a Adder) AddTo(val int) int {
	return a.start + val
}

// Type declarations are not inheritance, you need conversion to assign from one type to another
// and methods defined for one type won't be shared with types that are defined by it

// Types are documentation, they provide names for concepts and describe the kind of data expected

// iota is for Enumerations - Sometimes
type MailCategory int

const (
	Uncategorized MailCategory = iota
	Personal
	Spam
	Social
	Advertisements
)

// Beware that this should only be used when the actual values don't matter,
// otherwise this is fragile to someone adding/removing values in between declarations.
// It should only be used for "internal" purposes, where the constants are only referred to by name

// When you add a type as the field of a struct, you're making it an embedded field
// Any fields or methods declared on an embedded field are promoted to the containing struct
type Employee struct {
	Name string
	Id   string
}

func (e Employee) Description() string {
	return fmt.Sprintf("%s (%s)", e.Name, e.Id)
}

type Manager struct {
	Employee
	Reports []Employee
}

func (m Manager) FindNewEmployees() {
	// ...
}

// You can embed any type within a struct, not just another struct

// If the containing struct has fields or methods with the same name
// you have to use the embedded field's type to refer to the obscured fields
type Inner struct {
	X int
}
type Outer struct {
	Inner
	X int
}

// Embedding is not inheritance
func (i Inner) IntPrinter(val int) string {
	return fmt.Sprintf("Inner: %d", val)
}

func (i Inner) Double() string {
	return i.IntPrinter(i.X * 2)
}

func (o Outer) IntPrinter(val int) string {
	return fmt.Sprintf("Outer: %d", val)
}

// The methods on an embedded field do count toward the method set fo the containing struct
// So they can make the containing struct implement an interface

// Interfaces are the only abstract type in Go
// Methods declared by an interface are called the method set of the interface
// Interfaces are usually named with "er" endings
type Stringer interface {
	String() string
}

// If the method set for a concrete type contains all of the methods in the method set
// for an interface, the concrete type implements the interface.
// This means that the concrete type can be assigned to a variable or field of the interface type
// Read 142 to 145 to learn why this is neat
// TLDR: it combines the flexibility from (non-existent) interfaces in dynamic languages with the
// self-documenting/explicit nature of interfaces in other typed languages

// Like embedding types in structs, you can embed interfaces into other interfaces

// Accept Interfaces, return Structs. Read 146 to 147
// TLDR: Returning interfaces ruins the benefit, coupling your "consumer" code with the interface

func checkType(i interface{}) string {
	switch i.(type) {
	case nil:
		return "nil"
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		return "I don't know"
	}
}

// Function types are a bridge to interfaces
// If your single function is likely to depend on many other functions or state that's
// not specified in its input parameters, use an interface parameter and define a
// function type to bridge a function to the interface
type myHandler func(string)

func (mh myHandler) handle(s string) {
	mh(s)
}

func main() {
	p := Person{
		FirstName: "David",
		LastName:  "Hans",
		Age:       34,
	}
	fmt.Println(p.String())

	fmt.Println("-- Pointer receivers --")
	c := Counter{}
	c.Increment()
	fmt.Println(c.String())
	c.Increment()
	fmt.Println(c.String())

	fmt.Println("-- nil receivers --")
	var p2 *Person
	// fmt.Println(p2.String()) // panics
	fmt.Println(p2.StringNilInstance())
	p2 = &Person{FirstName: "Marla", LastName: "Michell", Age: 34}
	fmt.Println(p2.StringNilInstance())
	// fmt.Println(p2.String()) // would't panic anymore

	fmt.Println("-- Methods are functions --")
	myAdder := Adder{start: 10}
	fmt.Println(myAdder.AddTo(5))
	// This is called a method value:
	f1 := myAdder.AddTo
	fmt.Println(f1(6))
	// This is called a method expression (it's created from the type itself):
	f2 := Adder.AddTo
	fmt.Println(f2(myAdder, 7))
	// The first parameter is the receiver for the method

	fmt.Println("-- Embedding for composition --")
	m := Manager{
		Employee: Employee{
			Name: "Bob Bobson",
			Id:   "42",
		},
	}
	fmt.Println(m.Description())

	o := Outer{
		Inner: Inner{X: 10},
		X:     20,
	}
	fmt.Println(o.X)
	fmt.Println(o.Inner.X)
	fmt.Println(o.Double()) // Doesn't use Outer.IntPrinter

	fmt.Println("-- Interfaces --")
	// Interfaces are implemented as 2 pointers, one to the underlying type and one to the underlying value
	var s *string
	fmt.Println(s == nil)
	var i interface{}
	fmt.Println(i == nil)
	i = s
	fmt.Println(i == nil) // false
	// for an interface to be nil, both the type and the value have to be nil

	// An empty interface means that the variable can hold any value
	i = 20
	i = "hello"
	i = struct {
		FirstName string
		LastName  string
	}{FirstName: "Fred", LastName: "Fredson"}

	// You can check for the type with type switches
	// And type assertions
	i = true
	_, ok := i.(string)
	if !ok {
		fmt.Println("i is not a string")
	}
	_, ok = i.(int)
	if !ok {
		fmt.Println("i is not an int")
	}
	fmt.Printf("i is: %s\n", checkType(i))

	// Type assertion can be useful when an interface may implement another interface (thus saving work)
	// See page 152 for an example

	var handler myHandler = func(s string) {
		fmt.Println(s)
	}
	handler.handle("Bridging")

}
