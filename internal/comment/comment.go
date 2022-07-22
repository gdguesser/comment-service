package comment

import (
	"context"
	"fmt"
)

// Comment - a representation of the comment structure for our service
type Comment struct {
	ID     int
	Slug   string
	Body   string
	Author string
}

type Store interface {
	GetComment(context.Context, string) (Comment, error)
}

// Service - is the struct which all our logic will be build on top of
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service instance
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {
	fmt.Println("Retrieving a comment")
	return Comment{}, nil
}
