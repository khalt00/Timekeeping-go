package model

import (
	"context"
	"errors"
	"fmt"
	"time"
	"timekeeping/db/dbsvc"
)

// AttendanceRecord represents an attendance record
type AttendanceRecord struct {
	ID             int    `json:"id"`
	EmployeeID     string `json:"employeeID"`
	Details        string `json:"details"`
	EventTimestamp string `json:"eventTimestamp"`
	EventDate      string `json:"eventDate"`
	EventHours     string `json:"eventHours"`
	RecordType     string `json:"recordType"`
}

// Function to insert attendance record
func InsertAttendanceRecord(record AttendanceRecord) error {
	conn := dbsvc.GetPostgresConn()
	query := `
		INSERT INTO attendance_records (employee_id, details, event_timestamp, event_date, event_hours, record_type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := conn.Exec(context.Background(), query, record.EmployeeID, record.Details, record.EventTimestamp, record.EventDate, record.EventHours, record.RecordType)
	return err
}

// UpdateAttendanceRecordDetails updates the details field of an attendance record
func UpdateAttendanceRecordDetails(recordID int, newDetails string) error {
	conn := dbsvc.GetPostgresConn()
	query := `
		UPDATE attendance_records
		SET details = $1
		WHERE id = $2
	`
	_, err := conn.Exec(context.Background(), query, newDetails, recordID)
	return err
}

// Function to get attendance records by employee ID
func GetAttendanceRecordsByEmployeeID(employeeID string) ([]AttendanceRecord, error) {
	conn := dbsvc.GetPostgresConn()
	query := `
		SELECT id, employee_id, details, event_timestamp, event_date, event_hours, record_type
		FROM attendance_records
		WHERE employee_id = $1
	`
	rows, err := conn.Query(context.Background(), query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []AttendanceRecord
	for rows.Next() {
		var record AttendanceRecord
		err := rows.Scan(&record.ID, &record.EmployeeID, &record.Details, &record.EventTimestamp, &record.EventDate, &record.EventHours, &record.RecordType)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// Function to get attendance records by date range
func GetAttendanceRecordsByDateRange(startDateStr, endDateStr string) ([]AttendanceRecord, error) {
	// Parse start date string
	startDate, err := time.Parse("02/01/2006", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	// Parse end date string
	endDate, err := time.Parse("02/01/2006", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Ensure end date is after start date
	if endDate.Before(startDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	// Get PostgreSQL connection
	conn := dbsvc.GetPostgresConn()

	// Prepare SQL query
	query := `
		SELECT id, employee_id, details, event_timestamp, event_date, event_hours, record_type
		FROM attendance_records
		WHERE event_date >= $1 AND event_date <= $2
	`

	// Execute query with start and end dates
	rows, err := conn.Query(context.Background(), query, startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse query results
	var records []AttendanceRecord
	for rows.Next() {
		var record AttendanceRecord
		err := rows.Scan(&record.ID, &record.EmployeeID, &record.Details, &record.EventTimestamp, &record.EventDate, &record.EventHours, &record.RecordType)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
