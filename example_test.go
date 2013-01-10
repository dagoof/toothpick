package toothpick

import (
	"fmt"
)

func ExampleInstruction() {
	digit := I(`\d+`)
	fmt.Println(digit.Match("22015"))
	// Output: 22015
}

