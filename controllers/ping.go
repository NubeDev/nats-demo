package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"net/http"
	"sync"
	"time"
)

// PingHostCore handles the core logic for pinging a host
func (c *Controller) PingHostCore(uuid string) (string, error) {
	msg, err := c.nc.Request("host."+uuid+".ping", nil, time.Second)
	if err != nil {
		return "", err
	}
	return string(msg.Data), nil
}

// GetHostTimeCore handles the core logic for getting the time from a host
func (c *Controller) GetHostTimeCore(uuid string) (string, error) {
	msg, err := c.nc.Request("host."+uuid+".time", nil, time.Second)
	if err != nil {
		return "", err
	}
	return string(msg.Data), nil
}

// PingHostAllCore handles the core logic for sending a ping request to all hosts
func (c *Controller) PingHostAllCore() ([]string, error) {
	responseChan := make(chan string, 10) // buffer of 10, can be adjusted
	var wg sync.WaitGroup

	sub, err := c.nc.SubscribeSync("host.ping.response")
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	err = c.nc.Publish("host.ping", nil)
	if err != nil {
		return nil, err
	}

	timeout := time.NewTimer(5 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-timeout.C:
				return
			default:
				msg, err := sub.NextMsg(1 * time.Second)
				if err == nats.ErrTimeout {
					continue
				}
				if err != nil {
					return
				}
				if msg != nil && len(msg.Data) > 0 {
					responseChan <- string(msg.Data)
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(responseChan)
	}()

	var responses []string
	for res := range responseChan {
		if res != "" {
			responses = append(responses, res)
		}
	}

	return responses, nil
}

// GinHandlerPingHost wraps PingHostCore for Gin
func (c *Controller) GinHandlerPingHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		response, err := c.PingHostCore(uuid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

// GinHandlerGetHostTime wraps GetHostTimeCore for Gin
func (c *Controller) GinHandlerGetHostTime() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		time, err := c.GetHostTimeCore(uuid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"time": time})
	}
}

// GinHandlerPingHostAll wraps PingHostAllCore for Gin
func (c *Controller) GinHandlerPingHostAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		responses, err := c.PingHostAllCore()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(responses) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"message": "No responses received", "responses": responses})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"hosts": responses})
		}
	}
}
