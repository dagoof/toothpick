package toothpick

import "testing"

func TestInstruction(t *testing.T) {
	valid_input := "12345"
	invalid_input := "hello world"
	digit := I(`\d+`)
	if digit.Match(valid_input) != Simple(valid_input) {
		t.Fatalf("Digit instruction did not match valid input")
	}
	if digit.Match(invalid_input) != Failure {
		t.Fatalf("Digit instruction matched invalid input")
	}
}
