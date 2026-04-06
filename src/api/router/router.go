package router

import (
	"github.com/gin-gonic/gin"
	"api/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/events", handlers.HandleSensorEvent)

	return r
}
