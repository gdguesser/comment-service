//go:build integration

package db

import (
	"context"
	"testing"

	"github.com/gdguesser/comment-service/internal/comment"
	"github.com/stretchr/testify/assert"
)

func TestCommentDatabase(t *testing.T) {
	t.Run("test get all comments", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		cmts, err := db.GetAllComments(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, cmts)
	})

	t.Run("test get comment", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		cmt, err := db.GetComment(context.Background(), "01GF70E4M1F5XN8X3GCQPKWKAF")
		assert.NoError(t, err)
		assert.NotEmpty(t, cmt)
	})

	t.Run("test create comment", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug:   "slug",
			Author: "author",
			Body:   "body",
		})
		assert.NoError(t, err)

		newCmt, err := db.GetComment(context.Background(), cmt.ID)
		assert.NoError(t, err)
		assert.Equal(t, "slug", newCmt.Slug)
	})

	t.Run("test delete comment", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug:   "new-slug",
			Author: "gabriel",
			Body:   "body",
		})
		assert.NoError(t, err)

		err = db.DeleteComment(context.Background(), cmt.ID)
		assert.NoError(t, err)

		_, err = db.GetComment(context.Background(), cmt.ID)
		assert.Error(t, err)
	})

	t.Run("test update comment", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug:   "new-slug",
			Author: "gabriel",
			Body:   "body",
		})
		assert.NoError(t, err)

		updatedCmt, err := db.UpdateComment(context.Background(), cmt.ID, comment.Comment{
			Slug:   "updated-slug",
			Author: "updated-gabriel",
			Body:   "updated-body",
		})
		assert.NoError(t, err)
		assert.Equal(t, "updated-slug", updatedCmt.Slug)
	})
}
