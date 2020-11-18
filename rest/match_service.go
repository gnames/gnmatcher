package rest

import (
	"github.com/gnames/gnlib/domain/entity/gn"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnmatcher"
)

// matchMatcherREST implements MatcherService interface.
type matcherREST struct {
	gnm  gnmatcher.GNMatcher
	port int
	enc  encode.Encoder
}

// NewMNewMatcherREST is a constructor for MatchREST.
func NewMatcherService(gnm gnmatcher.GNMatcher,
	port int, enc encode.Encoder) matcherREST {
	return matcherREST{
		gnm:  gnm,
		port: port,
		enc:  enc,
	}
}

// GetPort returns port number to the service.
func (mr matcherREST) Port() int {
	return mr.port
}

// Ping returns "pong" message if connection to the service did succed.
func (mr matcherREST) Ping() string {
	return "pong"
}

// GetVersion returns version number and build timestamp of gnmatcher.
func (mr matcherREST) GetVersion() gn.Version {
	return gn.Version{
		Version: gnmatcher.Version,
		Build:   gnmatcher.Build,
	}
}

// MatchAry takes a list of strings and matches them to known scientific names.
func (mr matcherREST) MatchAry(names []string) []*mlib.Match {
	return mr.gnm.MatchNames(names)
}

// Encode encodes an object into a byte slice.
func (mr matcherREST) Encode(obj interface{}) ([]byte, error) {
	return mr.enc.Encode(obj)
}

// Decode decodes an object from a bytes slice.
func (mr matcherREST) Decode(input []byte, output interface{}) error {
	return mr.enc.Decode(input, output)
}
