package handlers

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func PingHandler(uuid string) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		log.Printf("Ping received for UUID %s, replying with pong", uuid)
		m.Respond([]byte("pong"))
	}
}

func TimeHandler(uuid string) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		currentTime := time.Now().Format(time.RFC3339)
		log.Printf("Time request received for UUID %s, replying with current time", uuid)
		m.Respond([]byte(currentTime))
	}
}

func GeneralPingHandler(nc *nats.Conn, uuid string) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		response := fmt.Sprintf("uuid: %s", uuid)
		log.Printf("General ping received for UUID %s, replying with: %s", uuid, response)

		// Publish the response back to the response subject
		err := nc.Publish("host.ping.response", []byte(response))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}
