package gnmatcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/gogna/gnparser"
)

var (
	parser gnparser.GNparser
)

func TestGnmatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gnmatcher Suite")
}

var _ = BeforeSuite(func() {
	parser = gnparser.NewGNparser()
})
