package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("v2_ingress example", func() {

	var (
		fakeLoggregator *FakeLoggregatorIngressServer
	)

	BeforeEach(func() {
		var err error
		fakeLoggregator, err = NewFakeLoggregatorIngressServer(
			"./metron_agent_cert.crt",
			"./metron_agent_cert.key",
			"./loggregator_ca.crt")
		Expect(err).ShouldNot(HaveOccurred())
		err = fakeLoggregator.Start()
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		if fakeLoggregator != nil {
			fakeLoggregator.Stop()
		}
	})

	It("emits all the metrics without an sleep", func() {
		command := exec.Command(examplePath)
		command.Env = append(command.Env, os.Environ()...)
		command.Env = append(command.Env, "CA_CERT_PATH=./loggregator_ca.crt")
		command.Env = append(command.Env, "CERT_PATH=./metron_agent_cert.crt")
		command.Env = append(command.Env, "KEY_PATH=./metron_agent_cert.key")
		command.Env = append(command.Env, fmt.Sprintf("LOGGREGATOR_URL=%s", fakeLoggregator.Addr))
		command.Env = append(command.Env, "SLEEP_IN_SECONDS=0")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(session, 5*time.Second).Should(gexec.Exit(0))

		var envelope *loggregator_v2.Envelope

		Eventually(fakeLoggregator.ReceivedEnvelopes, 2*time.Second).Should(Receive(&envelope))

		Expect(envelope.GetEvent()).NotTo(BeNil())
		Expect(envelope.GetEvent().GetTitle()).To(Equal("Starting sample V2 Client"))

		for i := 0; i < 50; i++ {
			Eventually(fakeLoggregator.ReceivedEnvelopes, 2*time.Second).Should(Receive(&envelope))
			Expect(envelope.GetLog()).NotTo(BeNil())
		}
		for i := 0; i < 5; i++ {
			Eventually(fakeLoggregator.ReceivedEnvelopes, 2*time.Second).Should(Receive(&envelope))
			Expect(envelope.GetGauge()).NotTo(BeNil())
		}
	})

	It("emits all the metrics with sleep", func() {
		command := exec.Command(examplePath)
		command.Env = append(command.Env, os.Environ()...)
		command.Env = append(command.Env, "CA_CERT_PATH=./loggregator_ca.crt")
		command.Env = append(command.Env, "CERT_PATH=./metron_agent_cert.crt")
		command.Env = append(command.Env, "KEY_PATH=./metron_agent_cert.key")
		command.Env = append(command.Env, fmt.Sprintf("LOGGREGATOR_URL=%s", fakeLoggregator.Addr))
		command.Env = append(command.Env, "SLEEP_IN_SECONDS=2")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(session, 5*time.Second).Should(gexec.Exit(0))

		var envelope *loggregator_v2.Envelope

		Eventually(fakeLoggregator.ReceivedEnvelopes, 2*time.Second).Should(Receive(&envelope))

		Expect(envelope.GetEvent()).NotTo(BeNil())
		Expect(envelope.GetEvent().GetTitle()).To(Equal("Starting sample V2 Client"))

		for i := 0; i < 50; i++ {
			Eventually(fakeLoggregator.ReceivedEnvelopes, 2*time.Second).Should(Receive(&envelope))
			Expect(envelope.GetLog()).NotTo(BeNil())
		}
		for i := 0; i < 5; i++ {
			Eventually(fakeLoggregator.ReceivedEnvelopes, 2*time.Second).Should(Receive(&envelope))
			Expect(envelope.GetGauge()).NotTo(BeNil())
		}
	})

})
