package rest_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	gn "github.com/gnames/gnames/domain/entity"
	"github.com/gnames/gnames/lib/encode"
	"github.com/gnames/gnmatcher/domain/entity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const url = "http://:8080/"

var _ = Describe("Rest", func() {
	Describe("Ping()", func() {
		It("Gets pong from REST server", func() {
			var req []byte
			enc := encode.GNgob{}
			resp, err := http.Post(url+"ping", "application/x-binary", bytes.NewReader(req))
			Expect(err).To(BeNil())

			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			var response string
			enc.Decode(respBytes, &response)
			Expect(string(response)).To(Equal("pong"))
		})
	})

	Describe("Version()", func() {
		It("Gets Version from REST server", func() {
			var req []byte
			enc := encode.GNgob{}
			resp, err := http.Post(url+"version", "application/x-binary", bytes.NewReader(req))
			Expect(err).To(BeNil())
			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			var response entity.Version
			enc.Decode(respBytes, &response)
			Expect(response.Version).To(MatchRegexp(`^v\d+\.\d+\.\d+`))
		})
	})

	Describe("MatchAry()", func() {
		It("Finds exact matches for entered names", func() {
			var response []entity.Match
			enc := encode.GNgob{}
			request := []string{
				"Not name", "Bubo bubo", "Pomatomus",
				"Pardosa moesta", "Plantago major var major",
				"Cytospora ribis mitovirus 2",
				"A-shaped rods", "Alb. alba",
			}
			req, err := enc.Encode(request)
			Expect(err).To(BeNil())
			r := bytes.NewReader(req)
			resp, err := http.Post(url+"match", "application/x-binary", r)
			Expect(err).To(BeNil())
			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			enc.Decode(respBytes, &response)
			Expect(len(response)).To(Equal(8))

			bad := response[0]
			Expect(bad.Name).To(Equal("Not name"))
			Expect(bad.MatchType).To(Equal(gn.NoMatch))
			Expect(bad.MatchItems).To(BeNil())

			good := response[1]
			Expect(good.Name).To(Equal("Bubo bubo"))
			Expect(good.MatchType).To(Equal(gn.Exact))
			Expect(good.MatchItems[0].MatchStr).To(Equal("Bubo bubo"))

			full := response[4]
			Expect(full.Name).To(Equal("Plantago major var major"))
			Expect(full.MatchType).To(Equal(gn.Exact))
			Expect(full.VirusMatch).To(BeFalse())
			Expect(full.MatchItems[0].MatchStr).To(Equal("Plantago major major"))

			virus := response[5]
			Expect(virus.Name).To(Equal("Cytospora ribis mitovirus 2"))
			Expect(virus.MatchType).To(Equal(gn.Exact))
			Expect(virus.VirusMatch).To(BeTrue())
			Expect(virus.MatchItems[0].MatchStr).To(Equal("Cytospora ribis mitovirus 2"))

			noParse := response[6]
			Expect(noParse.Name).To(Equal("A-shaped rods"))
			Expect(noParse.MatchType).To(Equal(gn.NoMatch))
			Expect(noParse.MatchItems).To(BeNil())

			abbr := response[7]
			Expect(abbr.Name).To(Equal("Alb. alba"))
			Expect(abbr.MatchType).To(Equal(gn.NoMatch))
			Expect(abbr.MatchItems).To(BeNil())
		})

		It("Finds fuzzy matches for entered names", func() {
			var response []entity.Match
			request := []string{
				"Not name", "Pomatomusi",
				"Pardosa moeste", "Pardosamoeste",
				"Accanthurus glaucopareus",
			}
			enc := encode.GNgob{}
			req, err := enc.Encode(request)
			Expect(err).To(BeNil())
			resp, err := http.Post(url+"match", "application/x-binary", bytes.NewReader(req))
			Expect(err).To(BeNil())
			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			enc.Decode(respBytes, &response)

			Expect(len(response)).To(Equal(5))
			bad := response[0]
			Expect(bad.Name).To(Equal("Not name"))
			Expect(bad.MatchType).To(Equal(gn.NoMatch))
			Expect(bad.MatchItems).To(BeNil())

			uni := response[1]
			Expect(uni.Name).To(Equal("Pomatomusi"))
			Expect(uni.MatchType).To(Equal(gn.NoMatch))
			Expect(uni.MatchItems).To(BeNil())

			suffix := response[2]
			Expect(suffix.Name).To(Equal("Pardosa moeste"))
			Expect(suffix.MatchType).To(Equal(gn.Fuzzy))
			Expect(len(suffix.MatchItems)).To(Equal(1))
			Expect(suffix.MatchItems[0].EditDistance).To(Equal(1))
			Expect(suffix.MatchItems[0].EditDistanceStem).To(Equal(0))

			space := response[3]
			Expect(space.Name).To(Equal("Pardosamoeste"))
			Expect(space.MatchType).To(Equal(gn.Fuzzy))
			Expect(len(space.MatchItems)).To(Equal(1))
			Expect(space.MatchItems[0].EditDistance).To(Equal(2))
			Expect(space.MatchItems[0].EditDistanceStem).To(Equal(1))

			multi := response[4]
			Expect(multi.Name).To(Equal("Accanthurus glaucopareus"))
			Expect(multi.MatchType).To(Equal(gn.Fuzzy))
			Expect(len(multi.MatchItems)).To(Equal(2))
			Expect(multi.MatchItems[0].EditDistanceStem).To(Equal(1))
		})
	})
})
