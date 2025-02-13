package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/brightside-dev/ronin-fitness-be/database"
	"github.com/brightside-dev/ronin-fitness-be/database/client"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   client.DatabaseService
}

func New() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	container := NewContainer(NewServer.db)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(container),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
