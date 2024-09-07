package routes

import (
	"github.com/NubeDev/nats-demo/controllers"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func RegisterRoutes(router *gin.Engine, nc *nats.Conn, controller *controllers.Controller) {
	// Host routes
	router.GET("/hosts", controller.GinHandlerGetHosts())
	router.GET("/hosts/:uuid", controller.GinHandlerGetHost())
	router.POST("/hosts", controller.GinHandlerAddHost())
	router.PATCH("/hosts/:uuid", controller.GinHandlerUpdateHost())
	router.DELETE("/hosts/:uuid", controller.GinHandlerDeleteHost())
	router.POST("/modules/:uuid", controller.PostModuleRequest(nc))

	// Store
	router.GET("/store/:name", controller.GinGetStore())
	router.POST("/store/add/object", controller.GinStoreAddFile())
	router.POST("/store/get/object", controller.GinGetObject())

	// NATS based routes
	router.GET("/hosts/remote/ping/all", controller.GinHandlerPingHostAll())
	router.GET("/hosts/remote/ping/:uuid", controller.GinHandlerPingHost())
	router.GET("/hosts/remote/time/:uuid", controller.GinHandlerGetHostTime())
}
