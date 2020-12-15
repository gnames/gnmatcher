package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/encode"
	"github.com/stretchr/testify/assert"
)

const url = "http://:8080/api/v1/"

var bugs = []struct {
	name           string
	matchType      vlib.MatchTypeValue
	matchCanonical string
	desc           string
}{
	{
		name:           "Tillaudsia utriculata",
		matchType:      vlib.Fuzzy,
		matchCanonical: "Tillandsia utriculata",
		desc:           "#7 gnidump, Misspelling of Tillandsia",
	},
	{
		name:           "Drosohila melanogaster",
		matchType:      vlib.Fuzzy,
		matchCanonical: "Drosophila melanogaster",
		desc:           "#23 gnmatcher, Misspelling of Drosophila",
	},
	{
		name:           "Acacia horrida",
		matchType:      vlib.Exact,
		matchCanonical: "Acacia horrida",
		desc:           "#24 checking Exact match",
	},
	{
		name:           "Acacia horrida nur",
		matchType:      vlib.PartialExact,
		matchCanonical: "Acacia horrida",
		desc:           "#24 PartialExact match does not work",
	},
	{
		name:           "Acacia nur",
		matchType:      vlib.PartialExact,
		matchCanonical: "Acacia",
		desc:           "PartialExact 'Acacia nur'",
	},
	{
		name:           "Bubo bubo",
		matchType:      vlib.Exact,
		matchCanonical: "Bubo bubo",
		desc:           "#31 'Bubo bubo' matches partial instead of exact",
	},
	{
		name:           "Bubo",
		matchType:      vlib.Exact,
		matchCanonical: "Bubo",
		desc:           "#31 'Bubo' uninomials do not match",
	},
	{
		name:           "Acetothermia bacterium enrichment culture clone B13-B-61",
		matchType:      vlib.NoMatch,
		matchCanonical: "",
		desc:           "Should not be parsed",
	},
}

func TestBugs(t *testing.T) {
	enc := encode.GNjson{}
	req, err := enc.Encode(params())
	assert.Nil(t, err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"matches", "application/json", r)
	assert.Nil(t, err)
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var mtch []mlib.Match
	err = enc.Decode(respBytes, &mtch)
	assert.Nil(t, err)

	for i, v := range bugs {
		msg := fmt.Sprintf("%s -> %s", v.name, v.matchCanonical)
		assert.Equal(t, mtch[i].MatchType.String(), v.matchType.String(), msg)
		if mtch[i].MatchType == vlib.NoMatch {
			assert.Empty(t, mtch[i].MatchItems)
			continue
		}
		assert.Greater(t, len(mtch[i].MatchItems), 0, msg)
		hasItem := false
		for _, mi := range mtch[i].MatchItems {
			if mi.MatchStr == v.matchCanonical {
				hasItem = true
			}
		}
		msg = fmt.Sprintf("%s -> %s", mtch[i].Name, v.matchCanonical)
		assert.True(t, hasItem, msg)

	}
}

func params() []string {
	ns := make([]string, len(bugs))
	for i, v := range bugs {
		ns[i] = v.name
	}
	return ns
}
