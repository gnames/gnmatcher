package gnmatcher_test

import (
	. "github.com/gnames/gnmatcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gnmatcher", func() {
	Describe("NewConfig", func() {
		It("Creates a default GNparser", func() {
			Expect(1).To(Equal(1))
			cnf := NewConfig()
			Expect(cnf.JobsNum).To(Equal(8))
		})
	})
})
