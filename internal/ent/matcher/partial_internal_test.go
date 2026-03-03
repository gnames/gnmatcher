package matcher

import (
	"testing"

	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/stretchr/testify/assert"
)

type exactMatcherMock struct{}

func (exactMatcherMock) Init() error              { return nil }
func (exactMatcherMock) SetConfig(cfg config.Config) {}
func (exactMatcherMock) MatchCanonicalID(uuid string) bool { return false }

// TestProcessPartialGenusNoMatchReturnsEmptyResult verifies the fix for a
// nil pointer dereference: when WithUninomialFuzzyMatch is true but fuzzy
// matching finds nothing, processPartialGenus must return a non-nil
// emptyResult instead of nil.
func TestProcessPartialGenusNoMatchReturnsEmptyResult(t *testing.T) {
	ns := nameString{
		ID:   "abc-123",
		Name: "Pardosa maesta",
		Partial: &partial{
			Genus: "Pardosa",
		},
	}
	m := matcher{
		exactMatcher: exactMatcherMock{},
		fuzzyMatcher: fuzzyMatcherMock{},
		cfg:          config.Config{WithUninomialFuzzyMatch: true},
	}
	res, err := m.processPartialGenus(ns)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, vlib.NoMatch, res.MatchType)
	assert.Equal(t, ns.ID, res.ID)
	assert.Equal(t, ns.Name, res.Name)
}
