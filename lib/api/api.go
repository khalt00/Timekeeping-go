package api

import "github.com/gin-gonic/gin"

func ResponseWithErrorAndMessage(status int, err error, c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if err != nil {
		c.AbortWithStatusJSON(status, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(status, nil)
}

func ResponseWithStatusAndData(status int, data interface{}, c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, data)
}
