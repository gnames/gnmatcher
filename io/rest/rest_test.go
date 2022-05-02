package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"sort"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/gnvers"
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const url = "http://:8080/api/v0/"

func TestPing(t *testing.T) {
	resp, err := http.Get(url + "ping")
	assert.Nil(t, err)

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err)
	}

	assert.Equal(t, "pong", string(res))
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
	var response []mlib.Output
	enc := gnfmt.GNjson{}
	names := []string{
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
	request := mlib.Input{Names: names}
	req, err := enc.Encode(request)
	assert.Nil(t, err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"matches", "application/json", r)
	assert.Nil(t, err)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	_ = enc.Decode(respBytes, &response)
	assert.Equal(t, 9, len(response))

	bad := response[0]
	assert.Equal(t, "Not name", bad.Name)
	assert.Equal(t, vlib.NoMatch, bad.MatchType)
	assert.Nil(t, bad.MatchItems)

	good := response[1]
	assert.Equal(t, "Bubo bubo", good.Name)
	assert.Equal(t, vlib.Exact, good.MatchType)
	assert.Equal(t, "Bubo bubo", good.MatchItems[0].MatchStr)
	assert.Equal(t, 0, good.MatchItems[0].EditDistance)
	assert.Equal(t, 0, good.MatchItems[0].EditDistanceStem)

	full := response[4]
	assert.Equal(t, "Plantago major var major", full.Name)
	assert.Equal(t, vlib.Exact, full.MatchType)
	assert.Equal(t, "Plantago major major", full.MatchItems[0].MatchStr)

	virus := response[5]
	assert.Equal(t, "Cytospora ribis mitovirus 2", virus.Name)
	assert.Equal(t, vlib.Virus, virus.MatchType)
	assert.Equal(t, "Cytospora ribis mitovirus 2", virus.MatchItems[0].MatchStr)

	noParse := response[6]
	assert.Equal(t, "A-shaped rods", noParse.Name)
	assert.Equal(t, vlib.NoMatch, noParse.MatchType)
	assert.Nil(t, noParse.MatchItems)

	abbr := response[7]
	assert.Equal(t, "Alb. alba", abbr.Name)
	assert.Equal(t, vlib.NoMatch, abbr.MatchType)
	assert.Nil(t, abbr.MatchItems)

	cand := response[8]
	assert.Equal(t, "Candidatus Aenigmarchaeum subterraneum", cand.Name)
	assert.Equal(t, vlib.Exact, cand.MatchType)
	assert.Equal(t,
		"Aenigmarchaeum subterraneum",
		cand.MatchItems[0].MatchStr,
	)
}

func TestFuzzy(t *testing.T) {
	var response []mlib.Output
	names := []string{
		"Not name", "Pomatomusi",
		"Pardosa moeste", "Pardosamoeste",
		"Accanthurus glaucopareus",
		"Tillaudsia utriculata",
		"Drosohila melanogaster",
		"Acanthobolhrium crassicolle",
	}
	request := mlib.Input{Names: names}
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
	assert.Equal(t, "Not name", bad.Name)
	assert.Equal(t, vlib.NoMatch, bad.MatchType)
	assert.Nil(t, bad.MatchItems)

	uni := response[1]
	assert.Equal(t, "Pomatomusi", uni.Name)
	assert.Equal(t, vlib.NoMatch, uni.MatchType)
	assert.Nil(t, uni.MatchItems)

	suffix := response[2]
	assert.Equal(t, "Pardosa moeste", suffix.Name)
	assert.Equal(t, vlib.Fuzzy, suffix.MatchType)
	assert.Equal(t, 1, len(suffix.MatchItems))
	assert.Equal(t, 1, suffix.MatchItems[0].EditDistance)
	assert.Equal(t, 0, suffix.MatchItems[0].EditDistanceStem)

	// We do not have yet support for lost spaces.
	space := response[3]
	assert.Equal(t, "Pardosamoeste", space.Name)
	assert.Equal(t, vlib.NoMatch, space.MatchType)
	assert.Nil(t, space.MatchItems)

	fuzzy := response[4]
	assert.Equal(t, "Accanthurus glaucopareus", fuzzy.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy.MatchType)
	assert.Equal(t, 2, len(fuzzy.MatchItems))
	assert.Equal(t, 1, fuzzy.MatchItems[0].EditDistanceStem)

	// Added because stem was missing in canonicals table.
	// Still not sure why it was possible, the solution is in
	// creating canonical_stems if they are empty.
	fuzzy2 := response[5]
	assert.Equal(t, "Tillaudsia utriculata", fuzzy2.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy2.MatchType)
	assert.Equal(t, 1, len(fuzzy2.MatchItems))
	assert.Equal(t, 1, fuzzy2.MatchItems[0].EditDistanceStem)

	// Added because stem for Drosophila melanogaster was missing.
	// It was missing because canonical and stem are the same.
	fuzzy3 := response[6]
	assert.Equal(t, "Drosohila melanogaster", fuzzy3.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy3.MatchType)
	assert.Equal(t, 2, len(fuzzy3.MatchItems))
	assert.Equal(t, 1, fuzzy3.MatchItems[0].EditDistanceStem)

	fuzzy4 := response[7]
	assert.Equal(t, "Acanthobolhrium crassicolle", fuzzy4.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy4.MatchType)
	assert.Equal(t, 2, len(fuzzy4.MatchItems))
	assert.Equal(t, 3, fuzzy4.MatchItems[0].EditDistance)
	assert.Equal(t, 1, fuzzy4.MatchItems[1].EditDistance)
}

// Related to issue #43. Send a name with one suffix and get back not only
// names with the same suffix, but also ones with another suffix if available.
func TestStem(t *testing.T) {
	var response []mlib.Output
	request := mlib.Input{
		Names: []string{"Isoetes longissimum"},
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
	assert.Equal(t, 2, len(response[0].MatchItems))
}

// Related to issue #49. Optional search inside of species group.
func TestSpeciesGroup(t *testing.T) {
	assert := assert.New(t)
	var response []mlib.Output
	tests := []struct {
		msg              string
		withSpeciesGroup bool
		itemsNum         int
		matchTypes       []string
	}{
		{"with SpGroup", true, 2, []string{"Exact", "ExactSpeciesGroup"}},
		{"without SpGroup", false, 1, []string{"Exact"}},
	}
	for _, v := range tests {
		request := mlib.Input{
			Names:            []string{"Narcissus minor"},
			WithSpeciesGroup: v.withSpeciesGroup,
		}
		enc := gnfmt.GNjson{}
		req, err := enc.Encode(request)
		assert.Nil(err)
		resp, err := http.Post(url+"matches", "application/json", bytes.NewReader(req))
		assert.Nil(err)
		respBytes, err := io.ReadAll(resp.Body)
		assert.Nil(err)

		err = enc.Decode(respBytes, &response)
		assert.Nil(err)
		assert.Equal(v.itemsNum, len(response[0].MatchItems))

		mts := make([]string, len(response[0].MatchItems))
		for i, v := range response[0].MatchItems {
			mts[i] = v.MatchType.String()
		}
		sort.Strings(mts)
		assert.Equal(v.matchTypes, mts)
	}
}
