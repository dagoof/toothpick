package toothpick

type Result interface {
	Count() int
	Repr() interface{}
}

type Simple string
func (s Simple) Count() int { return len(s) }
func (s Simple) Repr() interface{} { return s }

var Failure = Simple("")

type Static struct {
	O Result
	C int
}
func (s Static) Count() int { return s.C }
func (s Static) Repr() interface{} { return s.O.Repr() }

type Annotated struct {
	S string
	O Result
}
func (a Annotated) Repr() interface{} {
	return map[string]interface{}{
		a.S: a.O.Repr(),
	}
}
func (a Annotated) Count() int { return a.O.Count() }

/*
type Annotated map[string]Result
func (a Annotated) Repr() interface{} {
	r := map[string]interface{}{ }
	for k, v := range a {
		r[k] = v.Repr()
	}
	return r
}
func (a Annotated) Count() int {
	sum := 0
	for _, o := range a {
		sum += o.Count()
	}
	return sum
}
*/

type Multi []Result
func (os Multi) Repr() interface{} {
	r := []interface{}{ }
	for _, o := range os {
		r = append(r, o.Repr())
	}
	return r
}
func (os Multi) Count() int {
	sum := 0
	for _, o := range os {
		sum += o.Count()
	}
	return sum
}

