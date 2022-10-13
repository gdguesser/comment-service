package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	Router  *chi.Mux
	Service CommentService
	Server  *http.Server
}

func NewHandler(service CommentService) *Handler {
	h := &Handler{
		Service: service,
	}
	h.Router = chi.NewRouter()
	h.Router.Use(middleware.Logger)
	h.Router.Use(middleware.Heartbeat("/alive"))
	h.Router.Use(JSONMiddleware)
	h.Router.Use(TimeoutMiddleware)
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.Get("/api/v1/comment", h.GetAllComments)
	h.Router.Get("/api/v1/comment/{id}", h.GetComment)
	h.Router.Post("/api/v1/comment", JWTAuth(h.PostComment))
	h.Router.Put("/api/v1/comment/{id}", JWTAuth(h.UpdateComment))
	h.Router.Delete("/api/v1/comment/{id}", JWTAuth(h.DeleteComment))
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("Shutting down gracefully")

	return nil
}
