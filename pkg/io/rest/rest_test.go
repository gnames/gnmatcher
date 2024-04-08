package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"slices"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/gnvers"
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const url = "http://:8080/api/v1/"

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
	var response mlib.Output
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
	matches := response.Matches
	assert.Equal(t, 9, len(matches))

	bad := matches[0]
	assert.Equal(t, "Not name", bad.Name)
	assert.Equal(t, vlib.NoMatch, bad.MatchType)
	assert.Nil(t, bad.MatchItems)

	good := matches[1]
	assert.Equal(t, "Bubo bubo", good.Name)
	assert.Equal(t, vlib.Exact, good.MatchType)
	assert.Equal(t, "Bubo bubo", good.MatchItems[0].MatchStr)
	assert.Equal(t, 0, good.MatchItems[0].EditDistance)
	assert.Equal(t, 0, good.MatchItems[0].EditDistanceStem)

	full := matches[4]
	assert.Equal(t, "Plantago major var major", full.Name)
	assert.Equal(t, vlib.Exact, full.MatchType)
	assert.Equal(t, "Plantago major major", full.MatchItems[0].MatchStr)

	virus := matches[5]
	assert.Equal(t, "Cytospora ribis mitovirus 2", virus.Name)
	assert.Equal(t, vlib.Virus, virus.MatchType)
	assert.Equal(t, "Cytospora ribis mitovirus 2", virus.MatchItems[0].MatchStr)

	noParse := matches[6]
	assert.Equal(t, "A-shaped rods", noParse.Name)
	assert.Equal(t, vlib.NoMatch, noParse.MatchType)
	assert.Nil(t, noParse.MatchItems)

	abbr := matches[7]
	assert.Equal(t, "Alb. alba", abbr.Name)
	assert.Equal(t, vlib.NoMatch, abbr.MatchType)
	assert.Nil(t, abbr.MatchItems)

	cand := matches[8]
	assert.Equal(t, "Candidatus Aenigmarchaeum subterraneum", cand.Name)
	assert.Equal(t, vlib.Exact, cand.MatchType)
	assert.Equal(t,
		"Aenigmarchaeum subterraneum",
		cand.MatchItems[0].MatchStr,
	)
}

func TestFuzzy(t *testing.T) {
	var response mlib.Output
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

	matches := response.Matches
	bad := matches[0]
	assert.Equal(t, "Not name", bad.Name)
	assert.Equal(t, vlib.NoMatch, bad.MatchType)
	assert.Nil(t, bad.MatchItems)

	uni := matches[1]
	assert.Equal(t, "Pomatomusi", uni.Name)
	assert.Equal(t, vlib.NoMatch, uni.MatchType)
	assert.Nil(t, uni.MatchItems)

	suffix := matches[2]
	assert.Equal(t, "Pardosa moeste", suffix.Name)
	assert.Equal(t, vlib.Fuzzy, suffix.MatchType)
	assert.Equal(t, 1, len(suffix.MatchItems))
	assert.Equal(t, 1, suffix.MatchItems[0].EditDistance)
	assert.Equal(t, 0, suffix.MatchItems[0].EditDistanceStem)

	// We do not have yet support for lost spaces.
	space := matches[3]
	assert.Equal(t, "Pardosamoeste", space.Name)
	assert.Equal(t, vlib.NoMatch, space.MatchType)
	assert.Nil(t, space.MatchItems)

	fuzzy := matches[4]
	assert.Equal(t, "Accanthurus glaucopareus", fuzzy.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy.MatchType)
	assert.Equal(t, 2, len(fuzzy.MatchItems))
	assert.Equal(t, 1, fuzzy.MatchItems[0].EditDistanceStem)

	// Added because stem was missing in canonicals table.
	// Still not sure why it was possible, the solution is in
	// creating canonical_stems if they are empty.
	fuzzy2 := matches[5]
	assert.Equal(t, "Tillaudsia utriculata", fuzzy2.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy2.MatchType)
	assert.Equal(t, 1, len(fuzzy2.MatchItems))
	assert.Equal(t, 1, fuzzy2.MatchItems[0].EditDistanceStem)

	// Added because stem for Drosophila melanogaster was missing.
	// It was missing because canonical and stem are the same.
	fuzzy3 := matches[6]
	assert.Equal(t, "Drosohila melanogaster", fuzzy3.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy3.MatchType)
	assert.Equal(t, 2, len(fuzzy3.MatchItems))
	assert.Equal(t, 1, fuzzy3.MatchItems[0].EditDistanceStem)

	fuzzy4 := matches[7]
	assert.Equal(t, "Acanthobolhrium crassicolle", fuzzy4.Name)
	assert.Equal(t, vlib.Fuzzy, fuzzy4.MatchType)
	assert.Equal(t, 2, len(fuzzy4.MatchItems))
	assert.Equal(t, 3, fuzzy4.MatchItems[0].EditDistance)
	assert.Equal(t, 1, fuzzy4.MatchItems[1].EditDistance)
}

func TestRelaxedFuzzy(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, name, res string
		editDist       int
	}{
		{"bubo", "Bbo bbo", "Bubo bubo", 2},
	}

	for _, v := range tests {
		params := mlib.Input{
			Names:                 []string{v.name},
			WithRelaxedFuzzyMatch: true,
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
		assert.GreaterOrEqual(len(matches[0].MatchItems), 1)
	}
}

// Related to issue #43. Send a name with one suffix and get back not only
// names with the same suffix, but also ones with another suffix if available.
func TestStem(t *testing.T) {
	var response mlib.Output
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
	assert.Equal(t, 2, len(response.Matches[0].MatchItems))
}

// Related to issue #49. Optional search inside of species group.
func TestSpeciesGroup(t *testing.T) {
	assert := assert.New(t)
	var response mlib.Output
	tests := []struct {
		msg, name        string
		withSpeciesGroup bool
		itemsNum         int
		matchTypes       []string
	}{
		{"with SpGroup", "Narcissus minor", true, 2, []string{"Exact", "ExactSpeciesGroup"}},
		{"without SpGroup", "Narcissus minor", false, 1, []string{"Exact"}},
		{"with nil SpGroup", "Pardosa moesta", true, 1, []string{"Exact"}},
		{"without nil SpGroup", "Pardosa moesta", false, 1, []string{"Exact"}},
	}
	for _, v := range tests {
		request := mlib.Input{
			Names:            []string{v.name},
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
		matches := response.Matches
		assert.Equal(v.itemsNum, len(matches[0].MatchItems))

		mts := make([]string, len(matches[0].MatchItems))
		for i, v := range matches[0].MatchItems {
			mts[i] = v.MatchType.String()
		}
		slices.Sort(mts)
		assert.Equal(v.matchTypes, mts)
	}
}

// related to issue #58, add an option too fuzzy match uninomials.
func TestFuzzyUninomial(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, name, res string
		ds             int
		matchType      vlib.MatchTypeValue
	}{
		{
			msg:       "fuzzy",
			name:      "Simulidae",
			res:       "Simuliidae",
			ds:        3,
			matchType: vlib.Fuzzy,
		},
		{
			msg:       "partialFuzzy",
			name:      "Pomatmus abcdefg",
			res:       "Pomatomus",
			ds:        1,
			matchType: vlib.PartialFuzzy,
		},
	}

	for _, v := range tests {
		params := mlib.Input{
			Names:                   []string{v.name},
			DataSources:             []int{v.ds},
			WithUninomialFuzzyMatch: true,
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
		assert.GreaterOrEqual(len(matches[0].MatchItems), 1)
		var isFuzzy bool
		for _, vv := range matches[0].MatchItems {
			if vv.MatchStr == v.res {
				assert.Equal(1, vv.EditDistance)
				assert.Equal(v.matchType, vv.MatchType)
				isFuzzy = true
			}
		}
		assert.True(isFuzzy)

	}
}

func TestDataSources(t *testing.T) {
	assert := assert.New(t)
	var response mlib.Output
	enc := gnfmt.GNjson{}
	tests := []struct {
		msg, name  string
		dss        []int
		matchDss   []int
		itemsNum   int
		matchTypes []string
	}{
		{"no ds", "Narcissus minor minor", []int{}, []int{165, 169}, 1, []string{"Exact"}},
		{"grin txmy", "Narcissus minor minor", []int{6}, []int{6}, 1, []string{"PartialExact"}},
	}

	for _, v := range tests {
		request := mlib.Input{
			Names:       []string{v.name},
			DataSources: v.dss,
		}
		req, err := enc.Encode(request)
		assert.Nil(err)
		resp, err := http.Post(url+"matches", "application/json", bytes.NewReader(req))
		assert.Nil(err)
		respBytes, err := io.ReadAll(resp.Body)
		assert.Nil(err)

		err = enc.Decode(respBytes, &response)
		assert.Nil(err)
		matches := response.Matches
		assert.Equal(v.itemsNum, len(matches[0].MatchItems))
		assert.Equal(v.matchDss, matches[0].MatchItems[0].DataSources)

		mts := make([]string, len(matches[0].MatchItems))
		for i, v := range matches[0].MatchItems {
			mts[i] = v.MatchType.String()
		}
		slices.Sort(mts)
		assert.Equal(v.matchTypes, mts)
	}
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	resp, err := http.Get(url + "matches/Narcissus+minor+minor?data_sources=6|7&species_group=true")
	assert.Nil(err)

	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	var response mlib.Output
	err = gnfmt.GNjson{}.Decode(respBytes, &response)
	assert.Nil(err)
	assert.Equal(true, response.Meta.WithSpeciesGroup)
	assert.Equal([]int{6, 7}, response.Meta.DataSources)
	assert.Equal(1, response.Meta.NamesNum)
	matches := response.Matches
	assert.Equal("Narcissus minor minor", matches[0].Name)
}
