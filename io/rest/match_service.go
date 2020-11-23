package rest

import (
	"github.com/gnames/gnlib/domain/entity/gn"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnmatcher"
)

// matchMatcherREST implements MatcherService interface.
type matcherService struct {
	gnm  gnmatcher.GNMatcher
	port int
	enc  encode.Encoder
}

// NewMNewMatcherREST is a constructor for MatchREST.
func NewMatcherService(gnm gnmatcher.GNMatcher,
	port int, enc encode.Encoder) MatcherService {
	return matcherService{
		gnm:  gnm,
		port: port,
		enc:  enc,
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

// GetVersion returns version number and build timestamp of gnmatcher.
func (mr matcherService) GetVersion() gn.Version {
	return gn.Version{
		Version: gnmatcher.Version,
		Build:   gnmatcher.Build,
	}
}

// MatchAry takes a list of strings and matches them to known scientific names.
func (mr matcherService) MatchNames(names []string) []*mlib.Match {
	return mr.gnm.MatchNames(names)
}

// Encode encodes an object into a byte slice.
func (mr matcherService) Encode(obj interface{}) ([]byte, error) {
	return mr.enc.Encode(obj)
}

// Decode decodes an object from a bytes slice.
func (mr matcherService) Decode(input []byte, output interface{}) error {
	return mr.enc.Decode(input, output)
}
