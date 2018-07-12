package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"code.cloudfoundry.org/go-loggregator"
)

func main() {
	tlsConfig, err := loggregator.NewIngressTLSConfig(
		os.Getenv("CA_CERT_PATH"),
		os.Getenv("CERT_PATH"),
		os.Getenv("KEY_PATH"),
	)
	if err != nil {
		log.Fatal("Could not create TLS config", err)
	}

	client, err := loggregator.NewIngressClient(
		tlsConfig,
		loggregator.WithAddr(os.Getenv("LOGGREGATOR_URL")),
	)

	if err != nil {
		log.Fatal("Could not create client", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.EmitEvent(
		ctx,
		"Starting sample V2 Client",
		"This sample V2 client is about to emit 50 log envelopes",
	)
	if err != nil {
		log.Fatalf("Failed to emit event: %s", err)
	}

	for i := 0; i < 50; i++ {
		client.EmitLog("some log goes here",
			loggregator.WithSourceInfo("something", "platform", "v2-example-source-instance"),
		)

		time.Sleep(10 * time.Millisecond)
	}

	for i := 0; i < 5; i++ {
		client.EmitGauge(
			loggregator.WithGaugeValue("foo", float64(i), "ms"),
			loggregator.WithGaugeSourceInfo("something", "0"),
		)
	}

	client.CloseSend()
	sleepInSeconds, _ := strconv.Atoi(os.Getenv("SLEEP_IN_SECONDS"))
	time.Sleep(time.Duration(sleepInSeconds) * time.Second)
}
