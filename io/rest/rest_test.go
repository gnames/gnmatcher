package rest_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gnames/gnlib/domain/entity/gn"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/encode"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const url = "http://:8080/api/v1/"

func TestPing(t *testing.T) {
	resp, err := http.Get(url + "ping")
	assert.Nil(t, err)

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, string(res), "pong")
}

func TestVer(t *testing.T) {
	enc := encode.GNjson{}
	resp, err := http.Get(url + "version")
	assert.Nil(t, err)
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var response gn.Version
	_ = enc.Decode(respBytes, &response)
	assert.Regexp(t, `^v\d+\.\d+\.\d+`, response.Version)
}

func TestExact(t *testing.T) {
	var response []mlib.Match
	enc := encode.GNjson{}
	request := []string{
		"Not name", "Bubo bubo", "Pomatomus",
		"Pardosa moesta", "Plantago major var major",
		"Cytospora ribis mitovirus 2",
		"A-shaped rods", "Alb. alba",
	}
	req, err := enc.Encode(request)
	assert.Nil(t, err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"matches", "application/json", r)
	assert.Nil(t, err)
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	_ = enc.Decode(respBytes, &response)
	assert.Equal(t, len(response), 8)

	bad := response[0]
	assert.Equal(t, bad.Name, "Not name")
	assert.Equal(t, bad.MatchType, vlib.NoMatch)
	assert.Nil(t, bad.MatchItems)

	good := response[1]
	assert.Equal(t, good.Name, "Bubo bubo")
	assert.Equal(t, good.MatchType, vlib.Exact)
	assert.Equal(t, good.MatchItems[0].MatchStr, "Bubo bubo")
	assert.Equal(t, good.MatchItems[0].EditDistance, 0)
	assert.Equal(t, good.MatchItems[0].EditDistanceStem, 0)

	full := response[4]
	assert.Equal(t, full.Name, "Plantago major var major")
	assert.Equal(t, full.MatchType, vlib.Exact)
	assert.False(t, full.VirusMatch)
	assert.Equal(t, full.MatchItems[0].MatchStr, "Plantago major major")

	virus := response[5]
	assert.Equal(t, virus.Name, "Cytospora ribis mitovirus 2")
	assert.Equal(t, virus.MatchType, vlib.Exact)
	assert.True(t, virus.VirusMatch)
	assert.Equal(t, virus.MatchItems[0].MatchStr, "Cytospora ribis mitovirus 2")

	noParse := response[6]
	assert.Equal(t, noParse.Name, "A-shaped rods")
	assert.Equal(t, noParse.MatchType, vlib.NoMatch)
	assert.Nil(t, noParse.MatchItems)

	abbr := response[7]
	assert.Equal(t, abbr.Name, "Alb. alba")
	assert.Equal(t, abbr.MatchType, vlib.NoMatch)
	assert.Nil(t, abbr.MatchItems)
}

func TestFuzzy(t *testing.T) {
	var response []mlib.Match
	request := []string{
		"Not name", "Pomatomusi",
		"Pardosa moeste", "Pardosamoeste",
		"Accanthurus glaucopareus",
		"Tillaudsia utriculata",
		"Drosohila melanogaster",
		"Acanthobolhrium crassicolle",
	}
	enc := encode.GNjson{}
	req, err := enc.Encode(request)
	assert.Nil(t, err)
	resp, err := http.Post(url+"matches", "application/json", bytes.NewReader(req))
	assert.Nil(t, err)
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	err = enc.Decode(respBytes, &response)
	assert.Nil(t, err)

	bad := response[0]
	assert.Equal(t, bad.Name, "Not name")
	assert.Equal(t, bad.MatchType, vlib.NoMatch)
	assert.Nil(t, bad.MatchItems)

	uni := response[1]
	assert.Equal(t, uni.Name, "Pomatomusi")
	assert.Equal(t, uni.MatchType, vlib.NoMatch)
	assert.Nil(t, uni.MatchItems)

	suffix := response[2]
	assert.Equal(t, suffix.Name, "Pardosa moeste")
	assert.Equal(t, suffix.MatchType, vlib.Fuzzy)
	assert.Equal(t, len(suffix.MatchItems), 1)
	assert.Equal(t, suffix.MatchItems[0].EditDistance, 1)
	assert.Equal(t, suffix.MatchItems[0].EditDistanceStem, 0)

	// support for missing spaces is limited, because we cannot
	// generate correct stemmed version from them, so many
	// of such names are not matched due to edit distance bigger than
	// the threshold.
	space := response[3]
	assert.Equal(t, space.Name, "Pardosamoeste")
	assert.Equal(t, space.MatchType, vlib.Fuzzy)
	assert.Equal(t, len(space.MatchItems), 1)
	assert.Equal(t, space.MatchItems[0].EditDistance, 2)
	assert.Equal(t, space.MatchItems[0].EditDistanceStem, 1)

	fuzzy := response[4]
	assert.Equal(t, fuzzy.Name, "Accanthurus glaucopareus")
	assert.Equal(t, fuzzy.MatchType, vlib.Fuzzy)
	assert.Equal(t, len(fuzzy.MatchItems), 2)
	assert.Equal(t, fuzzy.MatchItems[0].EditDistanceStem, 1)

	// Added because stem was missing in canonicals table.
	// Still not sure why it was possible, the solution is in
	// creating canonical_stems if they are empty.
	fuzzy2 := response[5]
	assert.Equal(t, fuzzy2.Name, "Tillaudsia utriculata")
	assert.Equal(t, fuzzy2.MatchType, vlib.Fuzzy)
	assert.Equal(t, len(fuzzy2.MatchItems), 1)
	assert.Equal(t, fuzzy2.MatchItems[0].EditDistanceStem, 1)

	// Added because stem for Drosophila melanogaster was missing.
	// It was missing because canonical and stem are the same.
	fuzzy3 := response[6]
	assert.Equal(t, fuzzy3.Name, "Drosohila melanogaster")
	assert.Equal(t, fuzzy3.MatchType, vlib.Fuzzy)
	assert.Equal(t, len(fuzzy3.MatchItems), 2)
	assert.Equal(t, fuzzy3.MatchItems[0].EditDistanceStem, 1)

	fuzzy4 := response[7]
	assert.Equal(t, fuzzy4.Name, "Acanthobolhrium crassicolle")
	assert.Equal(t, fuzzy4.MatchType, vlib.Fuzzy)
	assert.Equal(t, len(fuzzy4.MatchItems), 2)
	assert.Equal(t, fuzzy4.MatchItems[0].EditDistance, 1)
	assert.Equal(t, fuzzy4.MatchItems[1].EditDistance, 3)
}