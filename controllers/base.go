package controllers

import (
	"github.com/NubeDev/nats-demo/models"
	"github.com/NubeDev/nats-demo/natsrouter"
	"github.com/nats-io/nats.go"
)

// Controller holds the NATS connection and hosts map
type Controller struct {
	nc         *nats.Conn
	hosts      map[string]models.Host
	natsRouter *natsrouter.NatsRouter
}

// NewController creates a new Controller with the given NATS connection
func NewController(nc *nats.Conn, natsRouter *natsrouter.NatsRouter) *Controller {
	return &Controller{
		nc:         nc,
		hosts:      make(map[string]models.Host),
		natsRouter: natsRouter,
	}
}
