package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"net/http"
	"time"
)

// ModuleRequest represents the request format for sending data to a module
type ModuleRequest struct {
	RequestUUID string `json:"requestUUID"`
	Method      string `json:"method"`
	Payload     string `json:"payload"`
}

func PostModuleRequest(nc *nats.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleID := c.Param("moduleID")

		var request struct {
			RequestUUID string `json:"requestUUID"`
			Method      string `json:"method"`
			Payload     string `json:"payload"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Construct NATS request subject
		natsSubject := "module." + moduleID + "." + request.Method

		// Send request to the module via NATS
		msg, err := nc.Request(natsSubject, []byte(request.Payload), 2*time.Second)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the module's response
		c.JSON(http.StatusOK, gin.H{
			"requestUUID": request.RequestUUID,
			"response":    string(msg.Data),
		})
	}
}
