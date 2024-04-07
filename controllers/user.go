package controllers

import (
	"net/http"
	"timekeeping/lib/api"
	"timekeeping/model"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		api.ResponseWithErrorAndMessage(http.StatusBadRequest, err, c)
		return
	}

	token, err := model.Login(req)
	if err != nil {
		api.ResponseWithErrorAndMessage(http.StatusBadRequest, err, c)
		return
	}

	api.ResponseWithStatusAndData(http.StatusOK, token, c)
}
