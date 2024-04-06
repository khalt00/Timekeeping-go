package main

import (
	"net/http"
	"timekeeping/db/dbsvc"
	"timekeeping/pooling"

	"github.com/gin-gonic/gin"
)

func main() {

	go pooling.Pooling()

	dbsvc.StartPostgres()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()

}
