package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

var (
	examplePath string
)

func TestV2Ingress(t *testing.T) {
	BeforeSuite(func() {
		var err error

		examplePath, err = gexec.Build("github.com/cloudfoundry/go-loggregator/examples/v2_ingress")
		Expect(err).ShouldNot(HaveOccurred())
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "V2Ingress Suite")
}
