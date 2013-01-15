package toothpick

import (
	"fmt"
)

func ExampleInstruction() {
	digit := I(`\d+`)
	fmt.Println(digit.Match("22015"))
	// Output: 22015
}

func ExampleMaybe() {
	space := I(`\s+`)
	digit := I(`\d+`)
	fmt.Println(M(digit, space).Match("194   "))
	// Output: {194 6}
}

func ExamplePrecond() {
	op := I(`[\+\*-/]`)
	digit := I(`\d+`)
	fmt.Println(Pre(op, digit).Match("+194"))
	// Output: {194 4}
}

func ExamplePostcond() {
	op := I(`[\+\*-/]`)
	digit := I(`\d+`)
	fmt.Println(Post(op, digit).Match("+194"))
	// Output: {+ 4}
}

func ExampleSequence() {
	space := I(`\s+`)
	word := M(I(`\w+`), space)
	op := M(I(`[\+\*-/]`), space)
	fmt.Println(S(word, op, word).Match("Hello +\thello"))
	// Output: [{Hello 6} {+ 2} hello]
}

func ExampleRepeated() {
	space := I(`\s+`)
	word := M(I(`\w+`), space)
	fmt.Println(R(word).Match("Hello there\thello"))
	// Output: [{Hello 6} {there 6} hello]
}

func ExampleGrammar() {
	space := I(`\s+`)
	r := Rules{
		"integer":  M(I(`\d+`), space),
		"operator": M(I(`[\+\*-/]`), space),
	}
	r["statement"] = S(
		r.Use("integer"),
		r.Use("operator"),
		r.Use("expression"))
	r["expression"] = O(r.Use("statement"), r.Use("integer"))
	g := Grammar{r, "expression"}
	fmt.Println(g.Match("2+3*5-14"))
	// Output: [2 + [3 * [5 - 14]]]
}

func ExampleAnnotate() {
	r := Rules{}

	space := I(`\s+`)
	lp := M(I(`'?\(`), space)
	rp := M(I(`\)`), space)
	word := M(I(`'?\w+`), space)

	r["atom"] = O(word, r.Use("list"))
	r.Set("list", Pre(lp, Post(
		S(r.Use("atom"), R(r.Use("atom"))), rp)))
	g := Grammar{r, "atom"}
	fmt.Println(g.Match("(cons 'a '(b c))"))
	// Output: {list {{[{cons 5} [{'a 3} {list {{[{b 2} [c]] 4} 6}}]] 15} 16}} 
}
