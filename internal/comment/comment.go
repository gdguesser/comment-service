package comment

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	ErrFetchingComment = errors.New("failed to fetch comment by id")
	ErrNotImplemented  = errors.New("not implemented")
)

// Comment - a representation of the comment structure for our service.
type Comment struct {
	ID     string `json:"id"`
	Slug   string `json:"slug"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

// Store - this interface defines all the methods that our service needs in order to operate.
type Store interface {
	GetAllComments(context.Context) ([]Comment, error)
	GetComment(context.Context, string) (Comment, error)
	PostComment(context.Context, Comment) (Comment, error)
	UpdateComment(context.Context, string, Comment) (Comment, error)
	DeleteComment(context.Context, string) error
}

// Service - is the struct which all our logic will be build on top of.
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service instance.
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

// GetAllComments - get comment from the database (repository layer).
func (s *Service) GetAllComments(ctx context.Context) ([]Comment, error) {
	cmts, err := s.Store.GetAllComments(ctx)
	if err != nil {
		log.Printf("Error retrieving all comments: %s\n", err)
		return []Comment{}, err
	}

	return cmts, nil
}

// GetComment - get comment from the database (repository layer).
func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {
	fmt.Println("Retrieving a comment")
	cmt, err := s.Store.GetComment(ctx, id)
	if err != nil {
		log.Printf("Error retrieving comment: %s\n", err)
		return Comment{}, ErrFetchingComment
	}

	return cmt, nil
}

// UpdateComment - updates a comment in the database (repository layer).
func (s *Service) UpdateComment(ctx context.Context, id string, updatedCmt Comment) (Comment, error) {
	cmt, err := s.Store.UpdateComment(ctx, id, updatedCmt)
	if err != nil {
		log.Println("error updating comment")
		return Comment{}, err
	}

	return cmt, nil
}

// DeleteComment - deletes a comment from the database (repository layer).
func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return s.Store.DeleteComment(ctx, id)
}

// PostComment - post a comment in the database (repository layer).
func (s *Service) PostComment(ctx context.Context, cmt Comment) (Comment, error) {
	insertedCmt, err := s.Store.PostComment(ctx, cmt)
	if err != nil {
		return Comment{}, err
	}

	return insertedCmt, nil
}
