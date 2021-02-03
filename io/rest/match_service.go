package rest

import (
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnmatcher"
)

// matchMatcherREST implements MatcherService interface.
type matcherService struct {
	gnmatcher.GNmatcher
	port int
	enc  gnfmt.Encoder
}

// NewMNewMatcherREST is a constructor for MatchREST.
func NewMatcherService(gnm gnmatcher.GNmatcher,
	port int, enc gnfmt.Encoder) MatcherService {
	return matcherService{
		GNmatcher: gnm,
		port:      port,
		enc:       enc,
	}
}

// GetPort returns port number to the service.
func (mr matcherService) Port() int {
	return mr.port
}

// Ping returns "pong" message if connection to the service did succed.
func (mr matcherService) Ping() string {
	return "pong"
}

// Encode encodes an object into a byte slice.
func (mr matcherService) Encode(obj interface{}) ([]byte, error) {
	return mr.enc.Encode(obj)
}

// Decode decodes an object from a bytes slice.
func (mr matcherService) Decode(input []byte, output interface{}) error {
	return mr.enc.Decode(input, output)
}
