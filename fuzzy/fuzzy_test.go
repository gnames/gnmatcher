package fuzzy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gnames/gnmatcher/fuzzy"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("Fuzzy", func() {
	DescribeTable("StripTags",
		func(s1, s2 string, expected int) {
			Expect(ComputeDistance(s1, s2)).To(Equal(expected))
		},
		Entry("identical", "Hello", "Hello", 0),
		Entry("distance 1", "Pomatomus", "Pom-tomus", 1),
		Entry("distance 2", "Pomatomus", "Poma  tomus", 2),
		Entry("distance utf8 1", "Pomatomus", "Pomщtomus", 1),
	)
})
