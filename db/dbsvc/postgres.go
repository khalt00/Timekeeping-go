package dbsvc

import (
	"context"
	"fmt"
	"log"
	"os"
	"timekeeping/lib/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func StartPostgresConnection(config config.Config) {
	postgresConn(config)
}

func postgresConn(dbConfig config.Config) {

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.SSLMode,
	)

	newPool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		os.Exit(1)
	}
	if err := newPool.Ping(context.Background()); err != nil {
		fmt.Println("error while pinging database", err)
		os.Exit(1)
	}

	pool = newPool

	log.Println("CONNECT POSTGRES SUCCESSFULLY")
}

func GetPostgresConn() *pgxpool.Pool {
	return pool
}
