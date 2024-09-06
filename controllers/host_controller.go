package controllers

import (
	"github.com/NubeDev/nats-demo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var hosts = map[string]models.Host{} // In-memory store for hosts

func GetHosts(c *gin.Context) {
	var allHosts []models.Host
	for _, host := range hosts {
		allHosts = append(allHosts, host)
	}
	c.JSON(http.StatusOK, allHosts)
}

func GetHost(c *gin.Context) {
	uuid := c.Param("uuid")
	host, exists := hosts[uuid]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Host not found"})
		return
	}
	c.JSON(http.StatusOK, host)
}

func AddHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hosts[host.UUID] = host
	c.JSON(http.StatusCreated, host)
}

func UpdateHost(c *gin.Context) {
	uuid := c.Param("uuid")
	var updatedHost models.Host
	if err := c.ShouldBindJSON(&updatedHost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := hosts[uuid]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Host not found"})
		return
	}
	hosts[uuid] = updatedHost
	c.JSON(http.StatusOK, updatedHost)
}

func DeleteHost(c *gin.Context) {
	uuid := c.Param("uuid")
	if _, exists := hosts[uuid]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Host not found"})
		return
	}
	delete(hosts, uuid)
	c.JSON(http.StatusNoContent, nil)
}
