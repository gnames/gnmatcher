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
			Expect(bi.Parts).To(BeNil())

			tri := nstrings[2].Partial
			Expect(tri.Genus).To(Equal("Pabstia"))
			Expect(len(tri.Parts)).To(Equal(1))
			Expect(tri.Parts[0].Tail).To(Equal("Pabstia parviflora"))
			Expect(tri.Parts[0].Head).To(Equal("Pabstia viridis"))

			quat := nstrings[3].Partial
			Expect(quat.Genus).To(Equal("Abedus"))
			Expect(len(quat.Parts)).To(Equal(2))
			Expect(quat.Parts[0].Tail).To(Equal("Abedus utahensis"))
			Expect(quat.Parts[0].Head).To(Equal("Abedus deinostoma herberti"))
		})
	})
})
