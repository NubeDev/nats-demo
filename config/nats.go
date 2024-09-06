package config

import (
	"github.com/nats-io/nats.go"
	"log"
)

func SetupNATS() (*nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
		return nil, err
	}
	log.Println("Connected to NATS")
	return nc, nil
}
