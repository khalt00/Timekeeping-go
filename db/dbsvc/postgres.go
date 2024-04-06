package dbsvc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var _conn *pgx.Conn

func StartPostgres() {
	var conn *pgx.Conn
	var err error

	// Retry up to 5 times
	for i := 0; i < 5; i++ {
		conn, err = pgx.Connect(context.Background(), "postgres://avnadmin:AVNS_g0TeOpoVwDSLwS3Uko_@pg-3a3ef85c-normals3210-5dd3.a.aivencloud.com:21277/defaultdb?sslmode=require")
		if err == nil {
			// Connection successful, break the loop
			break
		}

		// Print error and sleep before retrying
		fmt.Fprintf(os.Stderr, "Attempt %d: Unable to connect to database: %v\n", i+1, err)
		time.Sleep(5 * time.Second) // Adjust the sleep duration as needed
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database after 5 attempts: %v\n", err)
		os.Exit(1)
	}

	_conn = conn
}

func GetPostgresConnection() *pgx.Conn {
	return _conn
}
