package main

import (
	"flag"
	"fmt"
	"github.com/NubeDev/nats-demo/config"
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

	// Check if UUID is provided (edge device mode)
	go startEdgeDevice(nc, *uuid)
	startCloudServer(nc, *port)

}

func startEdgeDevice(nc *nats.Conn, uuid string) {
	log.Printf("Starting edge device with UUID: %s", uuid)

	// Initialize NATS router
	router := natsrouter.New(nc)

	// Register NATS routes
	router.Handle("host."+uuid+".ping", handlers.PingHandler(uuid))
	router.Handle("host."+uuid+".time", handlers.TimeHandler(uuid))
	nc.Subscribe("host.ping", handlers.GeneralPingHandler(nc, uuid))

	// Keep the edge device running indefinitely
	select {}
}
func startCloudServer(nc *nats.Conn, port string) {
	// Start the Gin server
	router := gin.Default()

	// Register routes
	routes.RegisterRoutes(router, nc)

	// Start server
	log.Println("Starting cloud server on port:", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
