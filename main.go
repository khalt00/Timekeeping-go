package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"timekeeping/db/dbsvc"
	"timekeeping/lib/auth"
	"timekeeping/lib/config"
	"timekeeping/pooling"
	"timekeeping/routes"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gin-gonic/gin"
)

func main() {

	config, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	dbsvc.StartPostgresConnection(config)
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode,
	)
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	fmt.Println("Migration successfully")
	go pooling.Pooling(config)

	r := gin.Default()
	r.Use(auth.CORSMiddleware())
	routes.InitialRoutes(r)

	// Start HTTP server in a separate goroutine
	go func() {
		if err := r.Run(":8081"); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()

	// Listen for SIGINT and SIGTERM signals to gracefully shut down the server
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// Shutdown server
	log.Println("Shutting down server...")

	// Migrate down
	// if err := m.Down(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal("error rolling back migrations:", err)
	// }
	fmt.Println("Migrations rolled back successfully!")

}
