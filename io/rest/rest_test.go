package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/gnvers"
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const url = "http://:8080/api/v1/"

func TestPing(t *testing.T) {
	resp, err := http.Get(url + "ping")
	assert.Nil(t, err)

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, string(res), "pong")
}

func TestVer(t *testing.T) {
	enc := gnfmt.GNjson{}
	resp, err := http.Get(url + "version")
	assert.Nil(t, err)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	var response gnvers.Version
	_ = enc.Decode(respBytes, &response)
	assert.Regexp(t, `^v\d+\.\d+\.\d+`, response.Version)
}

func TestExact(t *testing.T) {
	var response []mlib.Match
	enc := gnfmt.GNjson{}
	request := []string{
		"Not name",
		"Bubo bubo",
		"Pomatomus",
		"Pardosa moesta",
		"Plantago major var major",
		"Cytospora ribis mitovirus 2",
		"A-shaped rods",
		"Alb. alba",
		"Candidatus Aenigmarchaeum subterraneum",
	}
	req, err := enc.Encode(request)
	assert.Nil(t, err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"matches", "application/json", r)
	assert.Nil(t, err)
	respBytes, err := io.ReadAll(resp.Body)
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

	cand := response[8]
	assert.Equal(t, cand.Name, "Candidatus Aenigmarchaeum subterraneum")
	assert.Equal(t, cand.MatchType, vlib.Exact)
	assert.Equal(t,
		cand.MatchItems[0].MatchStr,
		"Aenigmarchaeum subterraneum",
	)
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
	enc := gnfmt.GNjson{}
	req, err := enc.Encode(request)
	assert.Nil(t, err)
	resp, err := http.Post(url+"matches", "application/json", bytes.NewReader(req))
	assert.Nil(t, err)
	respBytes, err := io.ReadAll(resp.Body)
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

	// We do not have yet support for lost spaces.
	space := response[3]
	assert.Equal(t, space.Name, "Pardosamoeste")
	assert.Equal(t, space.MatchType, vlib.NoMatch)
	assert.Nil(t, space.MatchItems)

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
