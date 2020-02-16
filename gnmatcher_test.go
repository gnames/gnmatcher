package gnmatcher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gnames/gnmatcher"
)

var _ = Describe("Gnmatcher", func() {
	Describe("NewGNmatcher", func() {
		It("Creates a default GNparser", func() {
			gnm := NewGNmatcher()
			Expect(gnm.JobsNum).To(Equal(4))
		})
	})

})
