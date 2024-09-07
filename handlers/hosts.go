package handlers

import (
	"encoding/json"
	"github.com/NubeDev/nats-demo/controllers"
	"github.com/NubeDev/nats-demo/models"
)

func hostHandler(endpoint, method, body string, controller *controllers.Controller) func() ([]byte, error) {
	switch endpoint {
	case "hosts":
		if method == "GET" {
			return func() ([]byte, error) {
				allHosts := controller.GetHostsCore()
				return json.Marshal(allHosts)
			}
		}
	case "host":
		if method == "GET" {
			return func() ([]byte, error) {
				host, exists := controller.GetHostCore(body)
				if !exists {
					return json.Marshal(map[string]string{"message": "Host not found"})
				}
				return json.Marshal(host)
			}
		}
	case "addHost":
		if method == "POST" {
			return func() ([]byte, error) {
				var newHost models.Host
				err := json.Unmarshal([]byte(body), &newHost)
				if err != nil {
					return nil, err
				}
				controller.AddHostCore(newHost)
				return json.Marshal(newHost)
			}
		}
	case "updateHost":
		if method == "PUT" {
			return func() ([]byte, error) {
				var updatedHost models.Host
				err := json.Unmarshal([]byte(body), &updatedHost)
				if err != nil {
					return nil, err
				}
				success := controller.UpdateHostCore(body, updatedHost)
				if !success {
					return json.Marshal(map[string]string{"message": "Host not found"})
				}
				return json.Marshal(updatedHost)
			}
		}
	case "deleteHost":
		if method == "DELETE" {
			return func() ([]byte, error) {
				success := controller.DeleteHostCore(body)
				if !success {
					return json.Marshal(map[string]string{"message": "Host not found"})
				}
				return []byte("Host deleted"), nil
			}
		}
	}
	return nil
}
