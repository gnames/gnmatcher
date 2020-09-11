package rest

import (
	"github.com/gnames/gnames/lib/encode"
	"github.com/gnames/gnmatcher"
	"github.com/gnames/gnmatcher/domain/entity"
)

// MatchMatcherREST implements MatcherService interface.
type MatcherREST struct {
	gnm  *gnmatcher.GNMatcher
	port int
	enc  encode.Encoder
}

// NewMNewMatcherREST is a constructor for MatchREST.
func NewMatcherREST(gnm *gnmatcher.GNMatcher, port int, enc encode.Encoder) MatcherREST {
	return MatcherREST{
		gnm:  gnm,
		port: port,
		enc:  enc,
	}
}

// GetPort returns port number to the service.
func (mr MatcherREST) Port() int {
	return mr.port
}

// Ping returns "pong" message if connection to the service did succed.
func (mr MatcherREST) Ping() string {
	return "pong"
}

// GetVersion returns version number and build timestamp of gnmatcher.
func (mr MatcherREST) Version() entity.Version {
	return entity.Version{
		Version: gnmatcher.Version,
		Build:   gnmatcher.Build,
	}
}

// MatchAry takes a list of strings and matches them to known scientific names.
func (mr MatcherREST) MatchAry(names []string) []*entity.Match {
	return mr.gnm.MatchNames(names)
}

// Encode encodes an object into a byte slice.
func (mr MatcherREST) Encode(obj interface{}) ([]byte, error) {
	return mr.enc.Encode(obj)
}

// Decode decodes an object from a bytes slice.
func (mr MatcherREST) Decode(input []byte, output interface{}) error {
	return mr.enc.Decode(input, output)
}
