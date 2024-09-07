package main

import (
	"flag"
	"fmt"
	"github.com/NubeDev/nats-demo/config"
	"github.com/NubeDev/nats-demo/controllers"
	"github.com/NubeDev/nats-demo/models"
	"github.com/NubeDev/nats-demo/natsrouter"

	"github.com/NubeDev/nats-demo/handlers"
	"github.com/NubeDev/nats-demo/routes"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	// Define a command line flag for the UUID (edge device mode)
	uuid := flag.String("uuid", "", "UUID for the edge device")
	port := flag.String("port", "", "port")
	flag.Parse()

	// Initialize NATS connection
	nc, err := config.SetupNATS()
	if err != nil {
		log.Fatalf("Error setting up NATS: %v", err)
	}
	defer nc.Close()
	fmt.Println("UUID:", *uuid)
	natsRouter := natsrouter.New(nc)
	err = natsRouter.CreateObjectStore("apps", nil)
	if err != nil {
		return
	}
	controller := controllers.NewController(nc, natsRouter)
	// Check if UUID is provided (edge device mode)
	go startEdgeDevice(*uuid, nc, natsRouter, controller)
	startCloudServer(nc, *port, controller)

}

func startEdgeDevice(uuid string, nc *nats.Conn, natsRouter *natsrouter.NatsRouter, controller *controllers.Controller) {
	log.Printf("Starting edge device with UUID: %s", uuid)
	// Initialize NATS router

	controller.AddHostCore(models.Host{
		UUID: "abc123",
		Name: "abc123",
		IP:   "0.0.0.0",
	})

	// Register NATS routes
	natsRouter.Handle("host."+uuid+".server", handlers.ServerHandler(controller))
	natsRouter.Handle("host."+uuid+".ping", handlers.PingHandler(uuid))
	natsRouter.Handle("host."+uuid+".time", handlers.TimeHandler(uuid))
	nc.Subscribe("host.ping", handlers.GeneralPingHandler(nc, uuid))

	// Keep the edge device running indefinitely
	select {}
}
func startCloudServer(nc *nats.Conn, port string, controller *controllers.Controller) {
	// Start the Gin server
	router := gin.Default()

	// Register routes
	routes.RegisterRoutes(router, nc, controller)

	// Start server
	log.Println("Starting cloud server on port:", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
