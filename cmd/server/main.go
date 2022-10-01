package main

import (
	"log"

	"github.com/gdguesser/comment-service/internal/comment"
	"github.com/gdguesser/comment-service/internal/db"

	transportHttp "github.com/gdguesser/comment-service/internal/transport/http"
)

// Run - is going to be responsible for the instantiation and startup of our go application
func Run() error {
	db, err := db.NewDatabase()
	if err != nil {
		log.Println("failed to connect to the database")
		return err
	}
	if err := db.MigrateDB(); err != nil {
		log.Println("failed to migrate database")
		return err
	}
	cmtService := comment.NewService(db)

	httpHandler := transportHttp.NewHandler(cmtService)
	log.Printf("starting up our application on %v", httpHandler.Server.Addr)
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Printf("Run failed: %s\n", err)
	}
}
