package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"net/http"
)

type StoreObject struct {
	StoreName           string `json:"storeName"`
	ObjectName          string `json:"objectName"`
	FilePath            string `json:"filePath"`
	OverwriteIfExisting bool   `json:"overwriteIfExisting"`
}

func (c *Controller) GetStore(storeName string) ([]*nats.ObjectInfo, error) {
	return c.natsRouter.GetStoreObjects(storeName)
}

func (c *Controller) GinGetObject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var b *StoreObject
		if err := ctx.ShouldBindJSON(&b); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		storeName := b.StoreName
		objectName := b.ObjectName

		data, err := c.natsRouter.GetObject(storeName, objectName)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Assuming you want to serve the object as raw data
		ctx.Header("Content-Disposition", "attachment; filename="+objectName)
		ctx.Data(http.StatusOK, "application/octet-stream", data)
	}
}

func (c *Controller) GinGetStore() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("name")
		resp, err := c.GetStore(param)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func (c *Controller) StoreAddFile(body *StoreObject) error {
	return c.natsRouter.NewObject(body.StoreName, body.ObjectName, body.FilePath, body.OverwriteIfExisting)
}

func (c *Controller) GinStoreAddFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var b *StoreObject
		if err := ctx.ShouldBindJSON(&b); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := c.StoreAddFile(b)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, b)
	}
}
