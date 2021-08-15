package utils

import (
	"log"
	"os"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

const kLocalhostPulsarUrl = "pulsar://localhost:6650"

func getPulsarAddr() string {
	if uri, exists := os.LookupEnv("PULSAR_URL"); exists {
		return uri
	} else {
		return kLocalhostPulsarUrl
	}
}

func CreatePulsarClientOrDie() pulsar.Client {
	options := pulsar.ClientOptions{
		URL:               getPulsarAddr(),
		OperationTimeout:  10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
	if client, err := pulsar.NewClient(options); err != nil {
		log.Fatalf("[FATAL] Failed to create pulsar client: %v", err)
		return nil
	} else {
		return client
	}
}
