package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gdguesser/comment-service/internal/comment"
)

type CommentRow struct {
	ID string
	Slug sql.NullString
	Body sql.NullString
	Author sql.NullString
}

func convertCommentRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Body:   c.Body.String,
		Author: c.Author.String,
	}
}

func (d *Database) GetComment(
	ctx context.Context,
	uuid string,
) (comment.Comment, error) {
	var cmtRow CommentRow
	row := d.Client.QueryRowContext(ctx,
		 `select id, slug, body, author 
		 from comments where id = $1`,
		uuid,)

	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching comment: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}