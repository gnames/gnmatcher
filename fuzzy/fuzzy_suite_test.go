package fuzzy_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFuzzy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fuzzy Suite")
}
