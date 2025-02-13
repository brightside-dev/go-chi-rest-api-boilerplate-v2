/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/brightside-dev/ronin-fitness-be/cmd/command/cmd"
	"github.com/brightside-dev/ronin-fitness-be/database/client"
	"github.com/brightside-dev/ronin-fitness-be/http"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Initialize the database service
	dbService := client.NewMySQL()

	// Initialize the container
	container := http.NewContainer(dbService)

	// Pass the container to the root command
	cmd.SetContainer(container)

	// Execute the root command
	cmd.Execute()
}
