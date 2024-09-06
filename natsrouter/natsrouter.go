package natsrouter

import (
	"github.com/nats-io/nats.go"
	"log"
)

type NatsRouter struct {
	nc *nats.Conn
}

type NatsHandlerFunc func(*nats.Msg)

// New creates a new NatsRouter
func New(nc *nats.Conn) *NatsRouter {
	return &NatsRouter{nc: nc}
}

// Handle registers a handler for a NATS subject
func (r *NatsRouter) Handle(subject string, handler NatsHandlerFunc) {
	r.nc.Subscribe(subject, nats.MsgHandler(handler))
}

// QueueHandle registers a handler for a NATS subject with a queue group
func (r *NatsRouter) QueueHandle(subject string, queue string, handler NatsHandlerFunc) {
	r.nc.QueueSubscribe(subject, queue, nats.MsgHandler(handler))
}

func (r *NatsRouter) Publish(subject string, reply string, data []byte) {
	err := r.nc.Publish(subject, data)
	if err != nil {
		log.Println("Error publishing message:", err)
	}
}
