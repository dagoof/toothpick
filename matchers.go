package toothpick

import (
	"fmt"
	"regexp"
)

// A Matcher consumes a string and returns a result
type Matcher interface {
	Match(string) Result
}

// Stores each successful match, acting as a memoized Matcher
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

// A Matcher that returns Annotated Results using the given string
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

// Regex matcher that only considers matches at the beginning of the provided
// string
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

// A list of matchers to be applied to an input
type Or Many

// Returns the first Matcher to successfuly return a Result rather than Failure
func (ms Or) Match(s string) Result {
	for _, m := range ms {
		matched := m.Match(s)
		if matched != Failure {
			return matched
		}
	}
	return Failure
}

type Sequence Many

// Every Matcher in the sequence must successfuly match the input, returns the
// Multi Result
func (ms Sequence) Match(s string) Result {
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

type Repeated struct{ M Matcher }

// Repeatedly attempts to match against a string, returning a Multi Result
func (m Repeated) Match(s string) Result {
	matched := Multi{}
	match := m.M.Match(s)
	for match != Failure {
		matched = append(matched, match)
		match = m.M.Match(s[matched.Count():])
	}
	return matched
}

type Maybe struct{ Always, Sometimes Matcher }

// Assuming the Always matcher succeeds, Maybe returns either the Result of
// Always if Sometimes fails, or the Static combination of Always' and
// Sometimes' Counts, along with Always' Result
func (m Maybe) Match(s string) Result {
	a := m.Always.Match(s)
	if a == Failure {
		return Failure
	}
	b := m.Sometimes.Match(s[a.Count():])
	if b == Failure {
		return a
	}
	return Static{a, a.Count() + b.Count()}
}

type Pair struct{ A, B Matcher }

type Precond Pair

// Assuming both Matchers succeed, the prior Matcher's Result is returned
func (p Precond) Match(s string) Result {
	a := p.A.Match(s)
	if a == Failure {
		return Failure
	}
	return p.B.Match(s[a.Count():])
}

type Postcond Pair

// Assuming both Matchers succeed, the latter Matcher's Result is returned
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

// Rules is a container that allows for recursive access of terms in a grammar.
type Rules map[string]Matcher

// Requesting a stored term through Use allows you to create mutually recursive
// definitions. When you request a term, a promise is returned, which will
// return the appropriate matcher only when the entire grammar is Matched
// against
func (r Rules) Use(s string) Matcher { return Promise{r, s} }

// An alternate method of setting terms in a ruleset, where the added Matcher is
// first converted to an Annotate object with the given key as the Annotate
// string.
func (r Rules) Set(s string, m Matcher) Matcher {
	r[s] = Annotate{m, s}
	return r[s]
}

// A complete grammar is composed of a set of Rules and a root term to start
// matching based on
type Grammar struct {
	R    Rules
	Root string
}

func (g Grammar) Match(s string) Result {
	return g.R[g.Root].Match(s)
}

// A Promise holds a set of rules and a key to be evaluated when matched against
type Promise struct {
	R Rules
	K string
}

func (p Promise) Match(s string) Result {
	return p.R[p.K].Match(s)
}
