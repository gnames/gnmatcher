package rpc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var (
	conn *grpc.ClientConn
)

func TestRpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rpc Suite")
}

var _ = BeforeSuite(func() {
	var err error
	conn, err = grpc.Dial(":8778", grpc.WithInsecure(), grpc.WithBlock())
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	conn.Close()
})
