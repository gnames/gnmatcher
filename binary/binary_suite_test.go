package binary_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBinary(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Binary Suite")
}
