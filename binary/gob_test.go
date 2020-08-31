package binary_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gnames/gnmatcher/binary"
	"github.com/gnames/gnmatcher/model"
)

var _ = Describe("Gob", func() {
	Describe("Encode and Decode", func() {
		It("encodes and decodes a string", func() {
			s := model.Pong("pong")
			res, err := Encode(s)
			Expect(err).To(BeNil())
			var pong model.Pong
			Decode(res, &pong)
			Expect(string(pong)).To(Equal("pong"))
		})
	})
})
