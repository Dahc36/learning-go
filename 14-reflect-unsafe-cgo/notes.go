package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"math/bits"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

var isLittleEndian bool

func init() {
	var x uint16 = 0xFF00
	xb := *(*[2]byte)(unsafe.Pointer(&x))
	isLittleEndian = (xb[0] == 0x00)
}

// Convert external binary data
// Say we are reading from a network, where the protocol has the following structure:
// - Value: 4 bytes, representing an unsigned, big-endian 32-bit int
// - Label: 10 bytes, ASCII name for the value
// - Active: 1 byte, boolean flag to indicate if the field is active
// - Padding: 1 byte, because we want to use 16 bytes
type Data struct {
	Value  uint32   // 4 bytes
	Label  [10]byte // 10 bytes
	Active bool     // 1 byte
	// Go padded this with 1 byte to make it align
}

// If we're reading:
// [0 132 95 237 80 104 111 110 101 0 0 0 0 0 1 0]
// With safe code, we would map it like:
func DataFromBytes(b [16]byte) Data {
	d := Data{}
	d.Value = binary.BigEndian.Uint32(b[:4])
	copy(d.Label[:], b[4:14])
	d.Active = b[14] != 0
	return d
}

// Or, we could use unsafe.Pointer instead:
func DataFromBytesUnsafe(b [16]byte) Data {
	data := *(*Data)(unsafe.Pointer(&b))
	if isLittleEndian {
		data.Value = bits.ReverseBytes32(data.Value)
	}
	return data
}

func BytesFromData(d Data) [16]byte {
	out := [16]byte{}
	binary.BigEndian.PutUint32(out[:4], d.Value)
	copy(out[4:14], d.Label[:])
	if d.Active {
		out[14] = 1
	}
	return out
}

func BytesFromDataUnsafe(d Data) [16]byte {
	if isLittleEndian {
		d.Value = bits.ReverseBytes32(d.Value)
	}
	b := *(*[16]byte)(unsafe.Pointer(&d))
	return b
}

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

	data := `name,age,has_pet
Jon,"100",true
"Fred ""The Hammer"" Smith",42,false
Martha,37,"true"
`

	r := csv.NewReader(strings.NewReader(data))
	allData, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	var entries []MyData
	Unmarshal(allData, &entries)
	fmt.Printf("%+v", entries)

	//now to turn entries into output
	out, err := Marshal(entries)
	if err != nil {
		panic(err)
	}
	sb := &strings.Builder{}
	// fmt.Println(sb)
	w := csv.NewWriter(sb)
	w.WriteAll(out)
	fmt.Println(sb)

	// You can build functions
	timed := makeTimedFunction(timeMe).(func(int) int)
	fmt.Println(timed(2))

	// Reflection is not free, it uses significantly more memory and performs more allocations.
	// It also makes your code less type safe

	fmt.Println("-- Unsafe --")
	// The unsafe library allows for memory management, it provides a fairly small API:
	// Methods:
	// - SizeOf takes in any variable and returns how many bytes it uses
	// - Offsetof takes in a field of a struct and returns the number of bytes from the start of the struct to the start of the field
	// - Alignof takes in a field or a variable and returns the byte alignment it requires
	// Types:
	// - unsafe.Pointer is a special type for a single purpose, a pointer of any type can be converted to or from it.
	//   It can also be converted to and from a special integer type, uintptr to do math with it,
	//   this allows for individual byte extraction for any type, pointer arithmetic and byte manipulation

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

func timeMe(a int) int {
	time.Sleep(time.Duration(a) * time.Second)
	result := a * 2
	return result
}

func makeTimedFunction(f interface{}) interface{} {
	ft := reflect.TypeOf(f)
	fv := reflect.ValueOf(f)
	wrapperF := reflect.MakeFunc(ft, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := fv.Call(in)
		end := time.Now()
		fmt.Println(end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}
