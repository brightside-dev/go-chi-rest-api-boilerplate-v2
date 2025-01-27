/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/cmd/command/cmd"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/database"
	"github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Initialize the database service
	dbService := database.New()

	// Initialize the container
	container := server.NewContainer(dbService)

	// Pass the container to the root command
	cmd.SetContainer(container)

	// Execute the root command
	cmd.Execute()
}
