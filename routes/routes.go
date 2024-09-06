package routes

import (
	"github.com/NubeDev/nats-demo/controllers"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func RegisterRoutes(router *gin.Engine, nc *nats.Conn) {
	// Host routes
	router.GET("/hosts", controllers.GetHosts)
	router.GET("/hosts/:uuid", controllers.GetHost)
	router.POST("/hosts", controllers.AddHost)
	router.PATCH("/hosts/:uuid", controllers.UpdateHost)
	router.DELETE("/hosts/:uuid", controllers.DeleteHost)
	router.POST("/modules/:moduleID", controllers.PostModuleRequest(nc))

	// NATS based routes
	router.GET("/hosts/remote/ping/all", controllers.PingHostAll(nc))
	router.GET("/hosts/remote/ping/:uuid", controllers.PingHost(nc))
	router.GET("/hosts/remote/time/:uuid", controllers.GetHostTime(nc))
	router.GET("/service/nats", controllers.GetNATSClients(nc))
}
