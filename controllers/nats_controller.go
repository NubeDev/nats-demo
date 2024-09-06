package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"sync"
	"time"
)

func PingHostAll(nc *nats.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Channel to collect responses
		responseChan := make(chan string, 10) // buffer of 10, can be adjusted
		var wg sync.WaitGroup

		// Subscribe to a subject to collect all responses
		sub, err := nc.SubscribeSync("host.ping.response")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer sub.Unsubscribe()

		// Send the ping request to all hosts
		log.Println("Sending ping request to all hosts")
		err = nc.Publish("host.ping", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Set a 5-second timeout
		timeout := time.NewTimer(5 * time.Second)

		// Start a goroutine to listen for responses until the timeout is reached
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Use select to either wait for a message or time out
				select {
				case <-timeout.C: // Timeout after 5 seconds
					log.Println("Timeout reached, stopping collection of responses")
					return
				default:
					// Listen for responses with a short timeout of 1 second to avoid blocking
					msg, err := sub.NextMsg(1 * time.Second)
					if err == nats.ErrTimeout {
						continue // Continue waiting for messages
					}
					if err != nil {
						log.Printf("Error receiving message: %v", err)
						return
					}
					if msg != nil && len(msg.Data) > 0 { // Ensure non-empty message data
						responseChan <- string(msg.Data)
						log.Printf("Message received: %s", string(msg.Data))
					}
				}
			}
		}()

		// Close the response channel after the goroutine finishes
		go func() {
			wg.Wait()
			close(responseChan)
		}()

		// Collect responses from the channel
		var responses []string
		for res := range responseChan {
			if res != "" { // Avoid adding empty strings to the response
				log.Printf("Adding response: %s", res)
				responses = append(responses, res)
			}
		}

		// Return the collected responses
		if len(responses) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No responses received", "responses": responses})
		} else {
			c.JSON(http.StatusOK, gin.H{"hosts": responses})
		}
	}
}

func PingHost(nc *nats.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		msg, err := nc.Request("host."+uuid+".ping", nil, time.Second)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": string(msg.Data)})
	}
}

func GetHostTime(nc *nats.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		msg, err := nc.Request("host."+uuid+".time", nil, time.Second)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"time": string(msg.Data)})
	}
}

func GetNATSClients(nc *nats.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch clients from NATS here
		// For example, use `nc.Request(...)` to ask NATS for a list of clients
		c.JSON(http.StatusOK, gin.H{"message": "Fetching NATS connected clients is not implemented yet."})
	}
}
