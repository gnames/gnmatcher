package gnmatcher

type MatchResult struct {
	Error  error
	Output string
}

func (gnm GNmatcher) MatchStream(in <-chan string, out chan<- MatchResult) {

}
