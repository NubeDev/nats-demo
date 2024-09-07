package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/nats-demo/controllers"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type RequestBody struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	Body     string `json:"body"`
}

// ServerHandler creates a handler for NATS messages.
func ServerHandler(controller *controllers.Controller) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		var reqBody RequestBody
		err := json.Unmarshal(m.Data, &reqBody)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			m.Respond([]byte("Error unmarshalling message"))
			return
		}

		handlerFunc := hostHandler(reqBody.Endpoint, reqBody.Method, reqBody.Body, controller)
		handlerFunc = pingHandler(reqBody.Endpoint, reqBody.Method, reqBody.Body, controller)

		if handlerFunc == nil {
			m.Respond([]byte("Unknown endpoint or method"))
			return
		}

		response, err := handlerFunc()
		if err != nil {
			log.Printf("Error processing request: %v", err)
			m.Respond([]byte("Error processing request"))
			return
		}

		m.Respond(response)
	}
}

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
