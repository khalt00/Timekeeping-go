package controllers

import (
	"errors"
	"net/http"
	"timekeeping/lib/api"
	"timekeeping/model"

	"github.com/gin-gonic/gin"
)

type GetAttendanceRecordsResponse struct {
}

func GetAttendancesByID(c *gin.Context) {
	employeeID := c.Param("id")
	if len(employeeID) < 1 {
		api.ResponseWithErrorAndMessage(http.StatusBadRequest, errors.New("empty employee id"), c)
		return
	}
	// Query attendance records by employee ID
	records, err := model.GetAttendanceRecordsByEmployeeID(employeeID)
	if err != nil {
		api.ResponseWithErrorAndMessage(http.StatusInternalServerError, errors.New("failed to retrieve attendance records"), c)
		return
	}

	api.ResponseWithStatusAndData(http.StatusOK, records, c)
}

func GetAttendancesByDateRange(c *gin.Context) {
	// Parse start and end dates from query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Query attendance records by date range
	records, err := model.GetAttendanceRecordsByDateRange(startDateStr, endDateStr)
	if err != nil {
		api.ResponseWithErrorAndMessage(http.StatusInternalServerError, err, c)
		return
	}

	// Return the attendance records
	api.ResponseWithStatusAndData(http.StatusOK, records, c)
}
