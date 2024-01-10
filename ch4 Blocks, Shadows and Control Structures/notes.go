package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("-- Shadowing --")
	// Inner block declarations shadow outer block ones
	x := 10
	if x > 5 {
		fmt.Println(x)
		x := 5
		fmt.Println(x)
	}
	fmt.Println(x)

	// Using := inside a block is an easy way of accidentally shadowing a variable
	x2 := 10
	if x2 > 5 {
		x2, y := 5, 20
		fmt.Println(x2, y)
	}
	fmt.Println(x2)
	// You should also be careful about shadowing imports, like fmt or any other

	// Go's built-in types, constants and functions are not keywords but
	// they exist in the Universe Block, so they can be shadowed too
	// fmt.Println(true)
	// true := 10
	// fmt.Println(true)
	// Don't do that...

	fmt.Println("-- If --")
	n := rand.Intn(10)
	if n == 0 {
		fmt.Println("That's too low")
	} else if n > 5 {
		fmt.Println("That's too big", n)
	} else {
		fmt.Println("That's just right", n)
	}

	// Go's If statements allows for creation of variables that are scoped
	// to the if, else if and else blocks (but not the outer scope)
	if n := rand.Intn(10); n == 0 {
		fmt.Println("That's too low")
	} else if n > 5 {
		fmt.Println("That's too big", n)
	} else {
		fmt.Println("That's just right", n)
	}
	fmt.Println("Old n:", n)
	// Technically, you can use any statement before the comparison but don't
	// Also beware of shadowing outer variables

	fmt.Println("-- For --")
	// Go has 4 different formats of for
	// 1. C-style for (a.k.a. complete for loop)
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
	// There are 3 parts, separated by semicolons
	//   Statement
	//     - := is mandatory, var is not allowed here
	//     - You can shadow a variable
	//   Comparison
	//     - Has to evaluate to a bool
	//   Increment
	//     - Any assignment is valid
	// 2. Condition-only for
	i := 10
	for i < 15 {
		fmt.Println(i)
		i++
	}
	// Just the comparison part, pretty much a while statement
	// 3. Infinite for
	i = 20
	for {
		fmt.Println("Hello", i)
		i++
		if i >= 25 {
			break
		}
	}
	// Runs forever... but can be stopped with break
	// Go also supports continue
	// 4. for-range
	// You can only use for-range to iterate over the built-in compound types
	// and user-defined types that are based on them
	evenVals := []int{2, 4, 6, 8, 10, 12}
	for i, v := range evenVals {
		fmt.Println(i, v)
	}
	// i is used for index in this case, but when iterating maps k would be better for key
	// v is short for value but could be something else depending on the types being iterated
	// For longer or nested loops, longer names may be better
	// If you don't need the first value, you can declare it with _
	// If you don't need the second value, you can leave it off
	uniqueNames := map[string]bool{"Fred": true, "Raul": true, "Wilma": true}
	for k := range uniqueNames {
		fmt.Println(k)
	}
	// Iteration on maps doesn't have a guaranteed order
	// Except when you fmt.Println them (they are sorted by key)
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}
	for k, v := range m {
		fmt.Println(k, v)
	}

	fmt.Println("-- For - string --")
	greet := "Good morning 🌞!"
	for i, c := range greet {
		fmt.Println(i, c, string(c))
	}
	// Indexes 14, 15 and 16 will be skipped over as 🌞 takes up 4 bytes to write.
	// This is because for-range iterates over the string's runes, not the bytes.
	// Multibyte runes are converted to a single 32-bit number (from their UTF-8 representation)
	// and the offset is incremented by the number of bytes in the rune.
	// If the loop encounters a byte that doesn't represent a valid UTF-8 value,
	// the Unicode replacement character (0xfffd) is returned instead

	// for-range values are copies
	evenVals = []int{2, 4, 6, 8, 10, 12}
	for _, v := range evenVals {
		v = v * 2 // This won't affect the values in evenVals
	}
	fmt.Println(evenVals)

	fmt.Println("-- For - labeling --")
	// You can label for statements to target them with continue or break statements
	samples := []string{"hello", "apple_n!"}
outer:
	for _, sample := range samples {
		for i, r := range sample {
			fmt.Println(i, r, string(r))
			if r == 'l' {
				continue outer
			}
		}
		fmt.Println()
	}

	// When to use each for loop
	// - for-range is the best way of iterating through the entire content of a compound type
	// - The c-style loop is useful when you only want to iterate a sub-section
	//   - Be careful that this won't work when iterating strings! (because of the bytes vs runes)
	// - The conditional for replaces the classic "while" statement
	// - The infinite for can replicate the "do-while" pattern and implement the iterator pattern

	fmt.Println("-- Switch --")
	// Go has better switch statements than your average language
	words := []string{"a", "cow", "smile", "gopher", "octopus", "anthropologist"}
	for _, word := range words {
		switch size := len(word); size {
		case 1, 2, 3, 4:
			fmt.Println(word, "is a short word!")
		case 5:
			fmt.Println(word, "is exactly the right length:", size)
		case 6, 7, 8, 9:
		default:
			fmt.Println(word, "is a long word!")
		}
	}
	// - There's no fallthrough by default on every case
	//   - If you want fallthrough you can use the fallthrough keyword but try to avoid it refactoring
	// - Cases with empty bodies just don't do anything
	// - break statements are not necessary but you can use them
	//   - Try to avoid and refactor
	//   - You may want to break out of an outer loop from inside a case, you'll need to label that loop

	fmt.Println("-- Switch - Blank Switches --")
	for _, word := range words {
		switch size := len(word); {
		case size < 5:
			fmt.Println(word, "is a short word!")
		case size == 5:
			fmt.Println(word, "is exactly the right length:", size)
		case size > 9:
			fmt.Println(word, "is a long word!")
		default:
		}
	}
	// Don't overdo it with this power, if you write a blank switch with only equality comparisons,
	// you should replace it with an expression (regular) switch statement

	// Use switch statements to group related cases
	// If you want to do unrelated comparisons use if/else

	fmt.Println("-- goto --")
	// Go's goto has some rules:
	// 1. You can't jump over declarations
	// 2. You can't jump into an inner or parallel block

	a := rand.Intn(10)
	for a < 100 {
		if a%5 == 0 {
			goto done
		}
		a = a*2 + 1
	}
	fmt.Println("do something when the loop completes normally")
done:
	fmt.Println("do complicated stuff no matter why we left the loop")
	fmt.Println(a)
	// Try to avoid using goto, but if it is the best option to make the code readable
	// it's theres
}
