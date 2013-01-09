package toothpick

import (
	"fmt"
	"regexp"
)

type Matcher interface {
	Match(string) Result
}

type CachedMatcher struct {
	Cache map[string]Result
	M     Matcher
}

func (m CachedMatcher) Match(s string) Result {
	if match, ok := m.Cache[s]; ok {
		return match
	}
	m.Cache[s] = m.M.Match(s)
	return m.Cache[s]
}

type Annotate struct {
	Matcher
	S string
}

func (a Annotate) Match(s string) Result {
	if match := a.Matcher.Match(s); match != Failure {
		return Annotated{a.S, match}
	}
	return Failure
}

type Instruction string

func (i Instruction) Match(s string) Result {
	pattern, _ := regexp.Compile(fmt.Sprintf(`^%s`, i))
	if pattern.MatchString(s) {
		matched := pattern.FindString(s)
		return Simple(matched)
	}
	return Failure
}

type Many []Matcher

type Or Many

func (ms Or) Match(s string) Result {
	for _, m := range ms {
		matched := m.Match(s)
		if matched != Failure {
			return matched
		}
	}
	return Failure
}

type Seq Many

func (ms Seq) Match(s string) Result {
	matched := Multi{}
	for _, m := range ms {
		match := m.Match(s[matched.Count():])
		if match == Failure {
			return Failure
		}
		matched = append(matched, match)
	}
	return matched
}

type Rep struct{ M Matcher }

func (m Rep) Match(s string) Result {
	matched := Multi{}
	match := m.M.Match(s)
	for match != Failure {
		matched = append(matched, match)
		match = m.M.Match(s[matched.Count():])
	}
	return matched
}

type Pair struct{ A, B Matcher }
type Maybe Pair

func (m Maybe) Match(s string) Result {
	a := m.A.Match(s)
	if a == Failure {
		return Failure
	}
	b := m.B.Match(s[a.Count():])
	if b == Failure {
		return a
	}
	return Static{a, a.Count() + b.Count()}
}

type Precond Pair

func (p Precond) Match(s string) Result {
	a := p.A.Match(s)
	if a == Failure {
		return Failure
	}
	return p.B.Match(s[a.Count():])
}

type Postcond Pair

func (p Postcond) Match(s string) Result {
	a := p.A.Match(s)
	if a == Failure {
		return Failure
	}
	b := p.B.Match(s[a.Count():])
	if b == Failure {
		return Failure
	}
	return a
}

type Rules map[string]Matcher
type Grammar struct {
	R    Rules
	Root string
}

func (r Rules) Use(s string) Matcher {
	return MM(Promise{r, s})
}

func (r Rules) Set(s string, m Matcher) Matcher {
	r[s] = Annotate{m, s}
	return r[s]
}

func (g Grammar) Match(s string) Result {
	return g.R[g.Root].Match(s)
}

type Promise struct {
	R Rules
	K string
}

func (p Promise) Match(s string) Result {
	return p.R[p.K].Match(s)
}
