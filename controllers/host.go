package controllers

import (
	"github.com/NubeDev/nats-demo/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetHostsCore retrieves all hosts
func (c *Controller) GetHostsCore() []models.Host {
	var allHosts []models.Host
	for _, host := range c.hosts {
		allHosts = append(allHosts, host)
	}
	return allHosts
}

// GetHostCore retrieves a specific host by UUID
func (c *Controller) GetHostCore(uuid string) (models.Host, bool) {
	host, exists := c.hosts[uuid]
	return host, exists
}

// AddHostCore adds a new host
func (c *Controller) AddHostCore(host models.Host) {
	c.hosts[host.UUID] = host
}

// UpdateHostCore updates an existing host
func (c *Controller) UpdateHostCore(uuid string, updatedHost models.Host) bool {
	if _, exists := c.hosts[uuid]; !exists {
		return false
	}
	c.hosts[uuid] = updatedHost
	return true
}

// DeleteHostCore deletes a host by UUID
func (c *Controller) DeleteHostCore(uuid string) bool {
	if _, exists := c.hosts[uuid]; !exists {
		return false
	}
	delete(c.hosts, uuid)
	return true
}

// GinHandlerGetHosts wraps GetHostsCore for Gin
func (c *Controller) GinHandlerGetHosts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		allHosts := c.GetHostsCore()
		ctx.JSON(http.StatusOK, allHosts)
	}
}

// GinHandlerGetHost wraps GetHostCore for Gin
func (c *Controller) GinHandlerGetHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		host, exists := c.GetHostCore(uuid)
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Host not found"})
			return
		}
		ctx.JSON(http.StatusOK, host)
	}
}

// GinHandlerAddHost wraps AddHostCore for Gin
func (c *Controller) GinHandlerAddHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var host models.Host
		if err := ctx.ShouldBindJSON(&host); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.AddHostCore(host)
		ctx.JSON(http.StatusCreated, host)
	}
}

// GinHandlerUpdateHost wraps UpdateHostCore for Gin
func (c *Controller) GinHandlerUpdateHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		var updatedHost models.Host
		if err := ctx.ShouldBindJSON(&updatedHost); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !c.UpdateHostCore(uuid, updatedHost) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Host not found"})
			return
		}
		ctx.JSON(http.StatusOK, updatedHost)
	}
}

// GinHandlerDeleteHost wraps DeleteHostCore for Gin
func (c *Controller) GinHandlerDeleteHost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")
		if !c.DeleteHostCore(uuid) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Host not found"})
			return
		}
		ctx.JSON(http.StatusNoContent, nil)
	}
}
