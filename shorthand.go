package toothpick

// Create a CachedMatcher from a Matcher
func C(m Matcher) Matcher {
	return CachedMatcher{map[string]Result{}, m}
}

// Create an Instruction from a string
func I(s string) Matcher { return Instruction(s) }

// Create an Or from a variable number of Matchers
func O(ms ...Matcher) Matcher { return Or(ms) }

// Create a Sequence from a variable number of Matchers
func S(ms ...Matcher) Matcher { return Sequence(ms) }

// Create a Repeated Matcher from a Matcher
func R(m Matcher) Matcher { return Repeated{m} }

// Create a Maybe from two Matchers
func M(always, sometimes Matcher) Matcher { return Maybe{always, sometimes} }

// Create a Precond from an optional and required Matcher
func Pre(optional, m Matcher) Matcher { return Precond{optional, m} }

// Create a Postcond from a required and optional Matcher
func Post(m, optional Matcher) Matcher { return Postcond{m, optional} }
