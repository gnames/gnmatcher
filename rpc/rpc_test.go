package rpc_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gnames/gnmatcher/protob"
)

var _ = Describe("Rpc", func() {
	Describe("Ping()", func() {
		It("Gets pong from gRPC server", func() {
			client := protob.NewGNMatcherClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			response, err := client.Ping(ctx, &protob.Void{})
			Expect(err).ToNot(HaveOccurred())
			Expect(response.Value).To(Equal("pong"))
		})
	})

	Describe("Ver()", func() {
		It("Returns gnMatcher Version from gRPC server", func() {
			client := protob.NewGNMatcherClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			response, err := client.Ver(ctx, &protob.Void{})
			Expect(err).ToNot(HaveOccurred())
			Expect(response.Version).To(MatchRegexp(`^v\d+\.\d+\.\d+`))
		})
	})

	Describe("MatchAry()", func() {
		It("Finds matches for entered names", func() {
			client := protob.NewGNMatcherClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			names := protob.Names{
				Names: []string{"Not name", "Bubo bubo", "Pomatomus",
					"Pardosa moesta", "Plantago major var major",
					"Cytospora ribis mitovirus 2"},
			}
			response, err := client.MatchAry(ctx, &names)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(response.Results)).To(Equal(6))

			bad := response.Results[0]
			Expect(bad.Name).To(Equal("Not name"))
			Expect(bad.MatchType).To(Equal(protob.MatchType_NONE))
			Expect(bad.MatchData).To(BeNil())

			good := response.Results[1]
			Expect(good.Name).To(Equal("Bubo bubo"))
			Expect(good.MatchType).To(Equal(protob.MatchType_CANONICAL))
			Expect(good.MatchData[0].MatchStr).To(Equal("Bubo bubo"))

			full := response.Results[4]
			Expect(full.Name).To(Equal("Plantago major var major"))
			Expect(full.MatchType).To(Equal(protob.MatchType_CANONICAL_FULL))
			Expect(full.MatchData[0].MatchStr).To(Equal("Plantago major var. major"))

			virus := response.Results[5]
			Expect(virus.Name).To(Equal("Cytospora ribis mitovirus 2"))
			Expect(virus.MatchType).To(Equal(protob.MatchType_VIRUS))
			Expect(virus.MatchData[0].MatchStr).To(Equal("Cytospora ribis mitovirus 2"))
		})
	})
})
