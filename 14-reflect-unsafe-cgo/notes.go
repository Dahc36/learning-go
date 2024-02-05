package main

import (
	"fmt"
	"reflect"
)

func main() {
	fmt.Println("-- Reflection --")
	// Reflection allows us to evaluate types at runtime
	// It also allows for the examination, creation and modification of variables,
	// functions and structs at runtime.
	// This is useful in general when interacting with data outside of the Go code we write

	// Type
	// Represents the type  of the variable passed to reflection.TypeOf
	// It is a value of type reflect.Type
	var x int
	xt := reflect.TypeOf(x)
	fmt.Println(xt.Name()) // int
	type foo struct {
		A int `myTag:"value"`
		B int `myTag:"value2"`
	}
	f := foo{}
	ft := reflect.TypeOf(f)
	fmt.Println(ft.Name()) // foo
	xpt := reflect.TypeOf(&x)
	fmt.Println(xpt.Name()) // empty string

	// Kind
	// Constant that says what the type is made of. E.g. if you defined a struct foo,
	// the kind is reflect.Struct and the type is "foo".
	// Some of the methods on reflect.Type and other types oin the reflect package only
	// make sense for certain kinds and will panic if used with the wrong kind.
	fmt.Println(xpt.Name())        // empty string
	fmt.Println(xpt.Kind())        // reflect.Ptr
	fmt.Println(xpt.Elem().Name()) // "int"
	fmt.Println(xpt.Elem().Kind()) // reflect.Int

	// There are methods on reflectType for reflecting on structs
	for i := 0; i < ft.NumField(); i++ {
		curField := ft.Field(i)
		fmt.Println(curField.Name, curField.Type.Name(), curField.Tag.Get("myTag"))
	}

	// Values
	s := []string{"a", "b", "c"}
	sv := reflect.ValueOf(s)               // sv is of type reflect.Value
	fmt.Println(sv.Interface().([]string)) // Using the method Interface we get back the value but loose the type
	// You can use reflection to set a value, but it's a 3 step process
	i := 10
	iv := reflect.ValueOf(&i)
	ivv := iv.Elem()
	// Just like there are special-case methods for reading, there are methods for setting primitives
	ivv.SetInt(20)
	fmt.Println("i:", i) // For other types you have to use the Set methods, which takes a variable of type reflect.Value
	// To create a value, we use reflect.New which takes in a reflect.Type and returns a reflect.Value
	// that's a pointer to a reflect.Value of the specified type.
	// You can also use reflection to replicate the make function, but you must start from a value
	var stringType = reflect.TypeOf((*string)(nil)).Elem() // reflect.Type that represents a string
	var stringSliceType = reflect.TypeOf([]string(nil))
	ssv := reflect.MakeSlice(stringSliceType, 0, 10)
	sv = reflect.New(stringType).Elem()
	sv.SetString("hello")
	ssv = reflect.Append(ssv, sv)
	ss := ssv.Interface().([]string)
	fmt.Println(ss) // [hello]

	// Check if an Interface's value is nil
	fmt.Println(hasNoValue(struct{}{}))
	// Check out page 309 to 314 to see reflection used to marshal and unmarshal CSV files
}

func hasNoValue(i interface{}) bool {
	iv := reflect.ValueOf(i)
	// Returns true if the reflect.Value holds anything other than a nil interface
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case
		reflect.Ptr,
		reflect.Slice,
		reflect.Map,
		reflect.Func,
		reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
