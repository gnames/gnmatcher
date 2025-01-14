package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gnames/gnfmt"
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/stretchr/testify/assert"
)

func TestVirus(t *testing.T) {
	tests := []struct {
		msg, name, matchStr string
		matchType           vlib.MatchTypeValue
		matchlen            int
	}{
		{
			msg:       "not virus",
			name:      "Something not a virus",
			matchType: vlib.NoMatch,
			matchlen:  0,
		},
		{
			msg:       "arct vir",
			name:      "Antarctic virus",
			matchStr:  "antarctic virus 1_i_cpgeorsw001ad",
			matchType: vlib.Virus,
			matchlen:  21,
		},
		{
			msg:       "bird",
			name:      "Bubo bubo",
			matchStr:  "bubo bubo",
			matchType: vlib.Exact,
			matchlen:  1,
		},
		{
			msg:       "vector",
			name:      "Cloning vector pAJM.011",
			matchStr:  "cloning vector pajm.011",
			matchType: vlib.Virus,
			matchlen:  1,
		},
		{
			msg:       "tobacco mosaic",
			name:      "Tobacco mosaic virus",
			matchStr:  "tobacco mosaic virus",
			matchType: vlib.Virus,
			matchlen:  17,
		},
		{
			msg:       "influenza overload",
			name:      "Influenza B virus",
			matchStr:  "influenza b virus",
			matchType: vlib.Virus,
			matchlen:  21,
		},
	}
	assert := assert.New(t)
	var response mlib.Output
	enc := gnfmt.GNjson{}
	names := make([]string, len(tests))
	for i := range tests {
		names[i] = tests[i].name
	}
	request := mlib.Input{Names: names}
	req, err := enc.Encode(request)
	assert.Nil(err)
	r := bytes.NewReader(req)
	resp, err := http.Post(url+"matches", "application/json", r)
	assert.Nil(err)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	_ = enc.Decode(respBytes, &response)
	assert.Equal(6, len(response.Matches))

	for i, v := range tests {
		res := response.Matches[i]
		assert.Equal(v.name, res.Name, v.msg)
		assert.Equal(v.matchType, res.MatchType, v.msg)
		assert.Equal(v.matchlen, len(res.MatchItems), v.msg)
		if len(res.MatchItems) > 0 {
			assert.Equal(v.matchStr, strings.ToLower(res.MatchItems[0].MatchStr))
			assert.Equal(v.matchType, res.MatchItems[0].MatchType, v.msg)
			assert.Equal(v.name, res.MatchItems[0].InputStr, v.msg)
		}
	}
}
