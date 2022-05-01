package rest

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

const url = "http://:8080/api/v0/"

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
		msg:            "#47 compares result with canonical simple",
		name:           "Teucrium pyrenaicum subsp. guarense",
		matchCanonical: "Teucrium pyrenaicum guarensis",
		matchType:      vlib.Fuzzy,
		editDistance:   2,
	},
	{
		msg:            "#48 should not have -1 edit distance",
		name:           "Vesicaria deltoideum creticum",
		matchCanonical: "Vesicaria cretica",
		matchType:      vlib.PartialFuzzy,
		editDistance:   2,
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

	var mtch []mlib.Output
	err = enc.Decode(respBytes, &mtch)
	assert.Nil(t, err)

	for i, v := range bugs {
		msg := fmt.Sprintf("%s -> %s", v.name, v.matchCanonical)
		assert.Equal(t, v.matchType.String(), mtch[i].MatchType.String(), msg)
		if mtch[i].MatchType == vlib.NoMatch {
			assert.Empty(t, mtch[i].MatchItems)
			continue
		}
		assert.Greater(t, len(mtch[i].MatchItems), 0, msg)
		var hasItem, hasEDist bool
		ed := 100
		for _, mi := range mtch[i].MatchItems {
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
		msg = fmt.Sprintf("%s -> %s", mtch[i].Name, v.matchCanonical)
		msgEDist := fmt.Sprintf("%s, ed: %d instead of %d", msg, ed, v.editDistance)
		assert.True(t, hasItem, msg)
		assert.True(t, ed >= 0)
		assert.True(t, hasEDist, msgEDist)

	}
}

func params() mlib.Input {
	ns := make([]string, len(bugs))
	for i, v := range bugs {
		ns[i] = v.name
	}
	return mlib.Input{Names: ns}
}
