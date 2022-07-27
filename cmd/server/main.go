package main

import (
	"fmt"
	"log"

	"github.com/gdguesser/comment-service/internal/db"
)

// Run - is going to be responsible for the instantiation and startup of our go application
func Run() error {
	fmt.Println("Starting up our application")
	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return err
	}
	if err := db.MigrateDB(); err != nil {
		log.Println("failed to migrate database")
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
