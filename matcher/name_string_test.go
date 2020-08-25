package matcher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gnames/gnmatcher/matcher"
)

var _ = Describe("NameString", func() {
	Describe("NewNameString", func() {
		It("Creates partials from a name string.", func() {
			names := []string{
				"Pomatomus",
				"Homo habilis L. Leakey, Tobias & Napier, 1964",
				"Pabstia viridis var. parviflora (Hoehne) Garay",
				"Abedus deinostoma herberti utahensis Menke 1960",
			}
			nstrings := make([]NameString, len(names))
			for i, name := range names {
				ns, _ := NewNameString(parser, name)
				nstrings[i] = ns
			}

			uni := nstrings[0].Partial
			Expect(uni).To(BeNil())

			bi := nstrings[1].Partial
			Expect(bi.Genus).To(Equal("Homo"))
			Expect(bi.Multinomials).To(BeNil())

			tri := nstrings[2].Partial
			Expect(tri.Genus).To(Equal("Pabstia"))
			Expect(len(tri.Multinomials)).To(Equal(1))
			Expect(tri.Multinomials[0].Tail).To(Equal("Pabstia parviflora"))
			Expect(tri.Multinomials[0].Head).To(Equal("Pabstia viridis"))

			quat := nstrings[3].Partial
			Expect(quat.Genus).To(Equal("Abedus"))
			Expect(len(quat.Multinomials)).To(Equal(2))
			Expect(quat.Multinomials[0].Tail).To(Equal("Abedus utahensis"))
			Expect(quat.Multinomials[0].Head).To(Equal("Abedus deinostoma herberti"))
		})
	})

})
