package handlers

import (
	"encoding/json"
	"github.com/NubeDev/nats-demo/controllers"
)

func pingHandler(endpoint, method, body string, controller *controllers.Controller) func() ([]byte, error) {
	switch endpoint {
	case "ping":
		if method == "GET" {
			return func() ([]byte, error) {
				r, err := controller.PingHostCore(body)
				if err != nil {
					return nil, err
				}
				return json.Marshal(r)
			}
		}
	}
	return nil
}
