package routes

import (
	"fmt"
	"net/http"
	"timekeeping/controllers"
	"timekeeping/lib/api"
	"timekeeping/lib/auth"

	"github.com/gin-gonic/gin"
)

func InitialRoutes(router *gin.Engine) {

	r := router.Group("/api")

	r.GET("/ping", func(c *gin.Context) {
		api.ResponseWithStatusAndData(http.StatusOK, "pong", c)
	})

	r.POST("/login", controllers.Login)

	authorized := r.Group("/")
	authorized.Use(auth.JWTMiddleware())
	{
		authorized.GET("/attendance/:id", controllers.GetAttendancesByID)
		authorized.GET("/attendance", controllers.GetAttendancesByDateRange)
	}

	fmt.Println("ROUTER CONNECTED")

}
