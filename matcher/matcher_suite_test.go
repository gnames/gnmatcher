package matcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/gogna/gnparser"
)

var (
	parser gnparser.GNparser
)

func TestMatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Matcher Suite")
}

var _ = BeforeSuite(func() {
	parser = gnparser.NewGNparser()
})
