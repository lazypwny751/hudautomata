package main

import (
	"time"
	"flag"
	"strconv"
	"net/http"
	"log/slog"

	"github.com/lazypwny751/hudautomata/pkg/routes"
	"github.com/gin-gonic/gin"
)

var (
	host       = flag.String("host", "127.0.0.1", "Gorgi Chat host")
	port       = flag.Int("port", 8080, "Gorgi Chat running port")
	// language   = flag.String("language", "en", "Language for the chat")

	dbHost     = flag.String("db-host", "127.0.0.1", "Database host")
	dbPort     = flag.Int("db-port", 3306, "Database port")
	dbUser     = flag.String("db-user", "gorgi", "Database user")
	dbPassword = flag.String("db-password", "", "Database password")
	dbName     = flag.String("db-name", "gorgi_chat", "Database name")
)

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)

	s := &http.Server{
		Addr:         *host + ":" + strconv.Itoa(*port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Starting server at", "url", "http://"+ *host + ":" + strconv.Itoa(*port))
	if err := s.ListenAndServe(); err != nil {
		slog.Error("Error starting server: ", "err", err)
		return
	}
}
