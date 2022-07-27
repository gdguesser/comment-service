package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gdguesser/go-rest-api/internal/db"
)

// Run - is going to be responsible for the instantiation and startup of our go application
func Run() error {
	fmt.Println("Starting up our application")
	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return err
	}

	if err := db.Ping(context.Background()); err != nil {
		return err
	}

	fmt.Println("Successfully connected and pinged the database")
	return nil
}

func main() {
	fmt.Println("Go REST API Course.")
	if err := Run(); err != nil {
		log.Printf("Run failed: %s\n", err)
	}
}
