package toothpick

func MM(m Matcher) Matcher {
	return CachedMatcher{map[string]Result{}, m}
}

func I(s string) Matcher { return MM(Instruction(s)) }

func O(ms ...Matcher) Matcher { return MM(Or(ms)) }

func S(ms ...Matcher) Matcher { return MM(Seq(ms)) }

func R(m Matcher) Matcher { return MM(Rep{m}) }

func M(a, b Matcher) Matcher { return MM(Maybe{a, b}) }
