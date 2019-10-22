package stemmer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gnames/gnmatcher/stemmer"
)

var _ = Describe("Stemmer", func() {
	Describe("Stem", func() {
		It("treats que suffix with exceptions", func() {
			Expect(Stem("detorque").Stem).To(Equal("detorque"))
			Expect(Stem("somethingque").Stem).To(Equal("something"))
		})
		It("removes suffixes correctly", func() {
			for k, v := range stemsDict {
				Expect(Stem(k).Stem).To(Equal(v))
			}
		})
	})

})
