package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gdguesser/comment-service/internal/comment"
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
	cmtService := comment.NewService(db)
	fmt.Println(cmtService.GetComment(
		context.Background(), 
		"1dde6f57-3e47-4b68-80c9-fed100947511",
		))

	fmt.Println("Successfully connected and pinged the database")
	return nil
}

func main() {
	fmt.Println("Go REST API Course.")
	if err := Run(); err != nil {
		log.Printf("Run failed: %s\n", err)
	}
}
