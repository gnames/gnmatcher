package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/encode"
	"github.com/stretchr/testify/assert"
)

const url = "http://:8080/"

var bugs = []struct {
	name           string
	matchType      vlib.MatchType
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
}

func TestBugs(t *testing.T) {
	enc := encode.GNjson{}
	req, err := enc.Encode(params())
	assert.Nil(t, err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"match", "application/json", r)
	assert.Nil(t, err)
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var mtch []mlib.Match
	err = enc.Decode(respBytes, &mtch)
	assert.Nil(t, err)

	for i, v := range bugs {
		assert.Greater(t, len(mtch[i].MatchItems), 0)
		assert.Equal(t, mtch[i].MatchType.String(), v.matchType.String())
	}
}

func params() []string {
	ns := make([]string, len(bugs))
	for i, v := range bugs {
		ns[i] = v.name
	}
	return ns
}
