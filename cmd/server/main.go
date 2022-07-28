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

	cmtService.Store.PostComment(context.Background(), comment.Comment{
		ID:     "027c1b0d-f295-473a-a767-f4c0847d9406",
		Slug:   "test-test",
		Body:   "inserting comment body context",
		Author: "Gabriel",
	})

	fmt.Println(cmtService.GetComment(
		context.Background(),
		"3225bc62-9f44-4af4-970a-8494e90e1147",
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
