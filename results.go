package toothpick

// An interface that describes the expected result from a successful match. It
// should contain both a count of how many characters were consumed by the match
// as well as a representation of any kind
type Result interface {
	Count() int
	Repr() interface{}
}

// The most primitive Result possible is a string
type Simple string

func (s Simple) Count() int        { return len(s) }
func (s Simple) Repr() interface{} { return s }

// Failure is the singleton Simple empty string
var Failure = Simple("")

// Static allows a result to cheat its length, by storing an integer. This is
// useful for pre and post conditional matchers who may want to leave off
// undesired matched elements such as whitespace
type Static struct {
	O Result
	C int
}

func (s Static) Count() int        { return s.C }
func (s Static) Repr() interface{} { return s.O.Repr() }

// An annotated struct describes its contents when Repr is output. Useful for
// labeling portions of output from a matched grammar. Annotated results are
// automatically created for Rules assigned by the Set function, as opposed to
// manually assigned matchers.
type Annotated struct {
	S string
	O Result
}

func (a Annotated) Count() int { return a.O.Count() }
func (a Annotated) Repr() interface{} {
	return map[string]interface{}{
		a.S: a.O.Repr(),
	}
}

type Multi []Result

func (os Multi) Count() int {
	sum := 0
	for _, o := range os {
		sum += o.Count()
	}
	return sum
}
func (os Multi) Repr() interface{} {
	r := []interface{}{}
	for _, o := range os {
		r = append(r, o.Repr())
	}
	return r
}
