package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gdguesser/comment-service/internal/comment"
	"github.com/gdguesser/comment-service/util"
)

// CommentRow - Represents a CommentRow thats returned by ..Row.. database methods.
type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

// convertCommentRowToComment - Converts a CommentRow to a Comment so it can be saved to the database.
func convertCommentRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Body:   c.Body.String,
		Author: c.Author.String,
	}
}

// GetAllComments - get retrieves all comments (persistence layer).
func (d *Database) GetAllComments(ctx context.Context) ([]comment.Comment, error) {
	rows, err := d.Client.QueryContext(ctx, `select * from comments`)
	if err != nil {
		return []comment.Comment{}, fmt.Errorf("error fetching comments: %w", err)
	}
	defer rows.Close()

	var comments []comment.Comment

	for rows.Next() {
		var cmt comment.Comment
		if err := rows.Scan(&cmt.ID, &cmt.Slug, &cmt.Author, &cmt.Body); err != nil {
			return comments, fmt.Errorf("error fetching comments: %w", err)
		}
		comments = append(comments, cmt)
	}

	if err = rows.Err(); err != nil {
		return comments, fmt.Errorf("error fetching comments: %w", err)
	}

	return comments, nil
}

// GetComment - get comment from the database (persistence layer).
func (d *Database) GetComment(
	ctx context.Context,
	uuid string,
) (comment.Comment, error) {
	var cmtRow CommentRow
	row := d.Client.QueryRowContext(ctx,
		`select id, slug, body, author 
		 from comments where id = $1`,
		uuid)

	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching comment: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}

// PostComment - post a comment in the database (persistence layer).
func (d *Database) PostComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error) {
	cmt.ID = util.GenerateULID()
	postRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(ctx,
		`insert into comments
	(id, slug, body, author)
	values
	(:id, :slug, :body, :author)`,
		postRow,
	)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to insert a comment in the database: %w", err)
	}
	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return cmt, nil
}

// UpdateComment - updates a comment in the database (persistence layer).
func (d *Database) UpdateComment(ctx context.Context, id string, cmt comment.Comment) (comment.Comment, error) {
	cmtRow := CommentRow{
		ID: id,
		Slug: sql.NullString{
			String: cmt.Slug,
			Valid:  true,
		},
		Body: sql.NullString{
			String: cmt.Body,
			Valid:  true,
		},
		Author: sql.NullString{
			String: cmt.Author,
			Valid:  true,
		},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`update comments set 
		slug = :slug,
		author = :author,
		body = :body
		where id = :id`,
		cmtRow,
	)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}

// DeleteComment - deletes a comment from the database (persistence layer).
func (d *Database) DeleteComment(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM comments where id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment from the database: %w", err)
	}
	return nil
}
