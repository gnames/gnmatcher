package rest_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gnames/gnfmt"
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/stretchr/testify/assert"
)

var bugs = []struct {
	msg, name, matchCanonical string
	matchType                 vlib.MatchTypeValue
	editDistance              int
}{
	{
		msg:            "#7 gnidump, Misspelling of Tillandsia",
		name:           "Tillaudsia utriculata",
		matchCanonical: "Tillandsia utriculata",
		matchType:      vlib.Fuzzy,
		editDistance:   1,
	},
	{
		msg:            "#123 502 for Phegoptera",
		name:           "Phegoptera",
		matchCanonical: "Phegoptera",
		matchType:      vlib.Exact,
		editDistance:   0,
	},
	{
		msg:            "#23 gnmatcher, Misspelling of Drosophila",
		name:           "Drosohila melanogaster",
		matchCanonical: "Drosophila melanogaster",
		matchType:      vlib.Fuzzy,
		editDistance:   1,
	},
	{
		msg:            "#24 checking Exact match",
		name:           "Acacia horrida",
		matchCanonical: "Acacia horrida",
		matchType:      vlib.Exact,
		editDistance:   0,
	},
	{
		msg:            "#24 PartialExact match does not work",
		name:           "Acacia horrida nur",
		matchCanonical: "Acacia horrida",
		matchType:      vlib.PartialExact,
		editDistance:   0,
	},
	{
		msg:            "PartialExact 'Acacia nur'",
		name:           "Acacia nur",
		matchCanonical: "Acacia",
		matchType:      vlib.PartialExact,
		editDistance:   0,
	},
	{
		msg:            "#31 'Bubo bubo' matches partial instead of exact",
		name:           "Bubo bubo",
		matchCanonical: "Bubo bubo",
		matchType:      vlib.Exact,
		editDistance:   0,
	},
	{
		msg:            "#45 'Isoetis longisima' fuzzy match is not found",
		name:           "Isoetes longisima",
		matchCanonical: "Isoetes longissima",
		matchType:      vlib.Fuzzy,
		editDistance:   1,
	},
	{
		msg:            "#31 'Bubo' uninomials do not match",
		name:           "Bubo",
		matchCanonical: "Bubo",
		matchType:      vlib.Exact,
		editDistance:   0,
	},
	{
		msg:            "Should not be parsed",
		name:           "Acetothermia bacterium enrichment culture clone B13-B-61",
		matchCanonical: "",
		matchType:      vlib.NoMatch,
		editDistance:   0,
	},
	{
		msg:            "#48 should not have -1 edit distance",
		name:           "Vesicaria deltoideum creticum",
		matchCanonical: "Vesicaria cretica",
		matchType:      vlib.PartialFuzzy,
		editDistance:   2,
	},
	{
		msg:            "#59 should provide correct match type",
		name:           "Pleurotoma godiea avita F. Edwards, 1860",
		matchCanonical: "Pleurotoma anita",
		matchType:      vlib.PartialFuzzy,
		editDistance:   1,
	},
}

func TestBugs(t *testing.T) {
	enc := gnfmt.GNjson{}
	req, err := enc.Encode(params())
	assert.Nil(t, err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"matches", "application/json", r)
	assert.Nil(t, err)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	var res mlib.Output
	err = enc.Decode(respBytes, &res)
	assert.Nil(t, err)
	matches := res.Matches

	for i, v := range bugs {
		msg := fmt.Sprintf("%s -> %s", v.name, v.matchCanonical)
		assert.Equal(t, v.matchType.String(), matches[i].MatchType.String(), msg)
		if matches[i].MatchType == vlib.NoMatch {
			assert.Empty(t, matches[i].MatchItems)
			continue
		}
		assert.Greater(t, len(matches[i].MatchItems), 0, msg)
		var hasItem, hasEDist bool
		ed := 100
		for _, mi := range matches[i].MatchItems {
			if mi.MatchStr == v.matchCanonical {
				hasItem = true
			}
			if mi.EditDistance < ed {
				ed = mi.EditDistance
			}
			if mi.EditDistance == v.editDistance {
				hasEDist = true
			}
		}
		msg = fmt.Sprintf("%s -> %s", matches[i].Name, v.matchCanonical)
		msgEDist := fmt.Sprintf("%s, ed: %d instead of %d", msg, ed, v.editDistance)
		assert.True(t, hasItem, msg)
		assert.True(t, ed >= 0)
		assert.True(t, hasEDist, msgEDist)

	}
}

// Test #47: Make sure that infraspecies do match as fuzzy even if their
// stems are the same as matched name.
func TestFuzzyInfrasp(t *testing.T) {
	assert := assert.New(t)
	params := mlib.Input{
		Names:       []string{"Teucrium pyrenaicum subsp. guarense"},
		DataSources: []int{196},
	}
	enc := gnfmt.GNjson{}
	req, err := enc.Encode(params)
	assert.Nil(err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"matches", "application/json", r)
	assert.Nil(err)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	var res mlib.Output
	err = enc.Decode(respBytes, &res)
	assert.Nil(err)
	matches := res.Matches
	assert.Equal(
		"Teucrium pyrenaicum guarensis",
		matches[0].MatchItems[0].MatchStr,
	)
	assert.Equal(2, matches[0].MatchItems[0].EditDistance)
	assert.Equal(vlib.Fuzzy, matches[0].MatchItems[0].MatchType)
}

func params() mlib.Input {
	ns := make([]string, len(bugs))
	for i, v := range bugs {
		ns[i] = v.name
	}
	return mlib.Input{Names: ns}
}
