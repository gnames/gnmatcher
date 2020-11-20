package matcher

import mlib "github.com/gnames/gnlib/domain/entity/matcher"

type Matcher interface {
	Init()
	MatchNames(names []string) []*mlib.Match
}
