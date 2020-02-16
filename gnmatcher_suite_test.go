package gnmatcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGnmatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gnmatcher Suite")
}
