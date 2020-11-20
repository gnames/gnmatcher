package exact

type ExactMatcher interface {
	Init()
	MatchCanonicalID(uuid string) bool
	MatchNameStringID(uuid string) bool
}
